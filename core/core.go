package core

import (
	"html/template"
	"net/http"
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Core struct {
	*template.Template
	*mgo.Database
	*errContainer
	LoggedIn *User
	Path     string
}

func NewCore(db *mgo.Database) *Core {

	x := &Core{
		initTemplates(),
		db,
		NewErrorContainer(),
		nil,
		"",
	}

	return x
}

func (x *Core) View(w http.ResponseWriter, tpl string, data interface{}) {
	if err := x.ExecuteTemplate(w, tpl, data); err != nil {
		x.AddError(err)
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

// Posts : gets and returns all posts in the database, if there are none or there was an error it returns nil
func (x *Core) Posts() []*Post {

	var posts []*Post
	err := x.C("posts").Find(bson.M{}).All(&posts)
	if err != nil {
		x.AddError(err)
		return nil
	}
	for _, p := range posts {
		p.Link(x)
		p.CalculateScore()
	}

	x.SortHighestPostTo("top", posts)

	return posts

}

func (x *Core) Projects() []*Project {
	var projects []*Project

	if err := x.C("projects").Find(bson.M{}).All(&projects); err != nil {
		x.AddError(err)
		return nil
	}

	for _, p := range projects {
		p.Link(x)
	}

	return projects

}

func (x *Core) SortHighestPostTo(where string, posts []*Post) {
	if where == "top" {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Score > posts[j].Score
		})
	}
}

func (x *Core) IconState() string {
	if x.LoggedIn != nil {
		return ""
	}
	return "pulse"
}

func (x *Core) IconClick() string {
	if x.LoggedIn != nil {
		return "#"
	}
	return "#modal"
}
