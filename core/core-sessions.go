package core

import (
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

///////////
// SESSIONS
///////////

func (x *Core) validateLogin(pu PotentialUser) (*User, bool) {

	var user *User
	err := x.C("users").Find(bson.M{"name": pu.Name}).One(&user)
	if err != nil {
		x.AddError(err)
		return nil, false
	}
	user.Link(x)

	if match := user.ComparePasswordWith(pu.Password); !match {
		x.AddError(errors.New("incorrect password"))
		return nil, false
	}
	return user, true
}

func (x *Core) SetCookie(w http.ResponseWriter, r *http.Request, u *User) {
	exp := time.Now().Add(24 * time.Hour)
	maxAge := (60 * 60) * 24

	sess := NewSession(u.ID, exp, x)

	if _, err := x.C("sessions").RemoveAll(bson.M{"user_id": sess.UserId}); err != nil {
		x.AddError(err)
		return
	}

	if err := sess.Save(); err != nil {
		x.AddError(err)
		http.Redirect(w, r, x.Path, http.StatusSeeOther)
		return
	}

	cookie := &http.Cookie{
		Name:   "session",
		Value:  sess.ID.Hex(),
		MaxAge: maxAge,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/profile/"+u.Name, http.StatusSeeOther)
}

func (x *Core) Login(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		x.AddError(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var pu PotentialUser
	if err := schema.NewDecoder().Decode(&pu, r.PostForm); err != nil {
		x.AddError(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//pu, err := decodeUser(r.Body);
	//if err != nil {
	//	x.AddError(err)
	//	http.Redirect(w, r, x.Path, http.StatusSeeOther)
	//	return
	//}

	usr, ok := x.validateLogin(pu)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Print("\n\n no errors validating \n\n")
	// Everything checks out, making session and cookie
	x.SetCookie(w, r, usr)

}

func (x *Core) Logout(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session")
	if err != nil {
		return
	}

	if err := x.C("sessions").Remove(bson.M{"_id": bson.ObjectIdHex(c.Value)}); err != nil {
		x.AddError(err)
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
		x.AddError(err)
	}

	var session *Session
	if err := x.C("sessions").Find(bson.M{"_id": bson.ObjectIdHex(cookie.Value)}).One(&session); err != nil {
		x.LoggedIn = nil
		cookie.MaxAge = -1
		log.Print("there is a cookie but no session, deleting cookie")
		http.SetCookie(w, cookie)
		return
	}

	var user *User
	if err := x.C("users").Find(bson.M{"_id": session.UserId}).One(&user); err != nil {
		x.LoggedIn = nil
		x.AddError(errors.New("there's a cookie and a session, but no user : " + err.Error()))
		return
	}

	x.LoggedIn = user

}
