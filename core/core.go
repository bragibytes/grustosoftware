package core

import (
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Core struct {
	*template.Template
	*mgo.Database
	*errContainer
	Path     string
	LoggedIn *User
}

func NewCore(db *mgo.Database) *Core {

	x := &Core{
		initTemplates(),
		db,
		NewErrorContainer(),
		"",
		nil,
	}

	return x
}

func (x *Core) View(w http.ResponseWriter) {
	if err := x.ExecuteTemplate(w, "index", x);err != nil {
		x.AddError(err)
	}
}

func initTemplates() *template.Template {
	tpl := template.Must(template.ParseGlob("views/*.gohtml"))
	template.Must(tpl.ParseGlob("views/components/*.gohtml"))
	template.Must(tpl.ParseGlob("views/pages/*.gohtml"))

	return tpl
}

func (x *Core) validateLogin(pu PotentialUser) bool {

	var user *User
	err := x.C("users").Find(bson.M{"name": pu.Name}).One(&user)
	if err != nil {
		x.AddError(err)
		return false
	}

	match := user.ComparePasswordWith(pu.Password)
	if !match {
		x.AddError(NewError("incorrect password", http.StatusBadRequest))
		return false
	}
	return true
}

func (x *Core) Login(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {

		x.AddError(NewError(err.Error(), http.StatusInternalServerError))
		http.Redirect(w, r, x.Path, http.StatusSeeOther)
	}

	var pu PotentialUser
	if err := schema.NewDecoder().Decode(&pu, r.PostForm); err != nil {
		x.AddError(NewError(err.Error(), http.StatusInternalServerError))
		http.Redirect(w, r, x.Path, http.StatusSeeOther)
	}
	if ok := x.validateLogin(pu); !ok {
		http.Redirect(w, r, x.Path, http.StatusSeeOther)
	}
	// Everything checks out, making session and cookie
	exp := time.Now().Add(24 * time.Hour)
	maxAge := (60 * 60) * 24

	usr := NewUser(pu, x)
	sess := NewSession(usr.ID, exp, x)

	if _, err := x.C("sessions").RemoveAll(bson.M{"user_id": sess.UserId}); err != nil {
		x.AddError(err)
	}

	sess.Save()

	cookie := &http.Cookie{
		Name:   "session",
		Value:  sess.ID.Hex(),
		MaxAge: maxAge,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)

}

func (x *Core) Logout(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session")
	if err != nil { return }

	if err := x.C("sessions").Remove(bson.M{"_id": bson.ObjectIdHex(c.Value)}); err != nil {
		x.AddError(NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.MaxAge = -1
	http.SetCookie(w, c)
}

func (x *Core) CheckState(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		x.LoggedIn = nil
		return
	}

	if _, err = x.C("sessions").RemoveAll(bson.M{"expires": bson.M{"$lt": time.Now()}}); err != nil {
		x.AddError(NewError("error removing expired sessions : "+err.Error(), http.StatusInternalServerError))
	}

	var session *Session
	if err := x.C("sessions").Find(bson.M{"_id": bson.ObjectIdHex(cookie.Value)}).One(&session);err != nil {
		x.LoggedIn = nil
		cookie.MaxAge = -1
		log.Print("there is a cookie but no session, deleting cookie")
		http.SetCookie(w, cookie)
		return
	}

	var user *User
	if err := x.C("users").Find(bson.M{"_id": session.UserId}).One(&user); err != nil {
		x.LoggedIn = nil
		x.AddError(NewError("there's a cookie and a session, but no user : "+err.Error(), http.StatusConflict))
		return
	}

	x.LoggedIn = user
}

func (x *Core) CountUsers() int {
	var users []*User

	if err := x.C("users").Find(bson.M{}).All(&users); err != nil {
		x.AddError(err)
		return 0
	}
	return len(users)
}

func (x *Core) CountActiveUsers() int {
	var sessions []Session
	err := x.C("sessions").Find(bson.M{}).All(&sessions)
	if err != nil {
		return 0
	}
	return len(sessions)
}

func (x *Core) Posts() []*Post {

	var posts []*Post
	err := x.C("posts").Find(bson.M{}).All(&posts)
	if err != nil {
		x.AddError(err)
		return nil
	}
	for _, p := range posts {
		p.Link(x)
	}

	return posts

}
