package core

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type Core struct {
	*template.Template
	*mgo.Database
	*errContainer
	LoggedIn *User
}

func NewCore(db *mgo.Database) *Core {

	x := &Core{
		initTemplates(),
		db,
		NewErrorContainer(),
		nil,
	}

	return x
}

func (x *Core) View(w http.ResponseWriter, tpl string, data interface{}) {
	if err := x.ExecuteTemplate(w, tpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func initTemplates() *template.Template {
	tpl := template.Must(template.ParseGlob("views/components/*.gohtml"))
	template.Must(tpl.ParseGlob("views/pages/*.gohtml"))

	return tpl
}

func (x *Core) UserCount() int {
	var users []*User

	if err := x.C("users").Find(bson.M{}).All(&users); err != nil {
		x.AddError(err)
		return 0
	}
	return len(users)
}

func (x *Core) SessionCount() int {
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

func (x *Core) IconState() string {
	if x.LoggedIn != nil {
		return "sentiment_very_satisfied"
	}
	return "sentiment_dissatisfied"
}


