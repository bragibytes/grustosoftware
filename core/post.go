package core

import (
	"github.com/pkg/errors"
	. "gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Post struct {
	ID        ObjectId `bson:"_id"`
	Title     string   `bson:"title"`
	Body      string   `bson:"body"`
	AuthorID  ObjectId `bson:"_author"`
	*Core     `bson:"-"`
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func NewPost(title, body string, author ObjectId, con *Core) *Post {
	x := &Post{
		Title: title,
		Body:  body,
		Core:  con,
	}

	return x
}

func (x *Post) Link(con *Core) {
	x.Core = con
}

func (x *Post) Comments() []*Comment {

	var comments []*Comment
	if err := x.C("comments").Find(M{"_parent": x.ID}).All(&comments); err != nil {
		x.AddError(err)
		return nil
	}

	for _, v := range comments {
		v.Link(x.Core)
	}

	return comments
}

func (x *Post) CommentCount() int {
	var comments []*Comment
	if err := x.C("comments").Find(M{"_parent": x.ID}).All(&comments); err != nil {
		x.AddError(err)
		return 0
	}
	return len(comments)
}

func (x *Post) IDHex() string {
	return x.ID.Hex()
}

func (x *Post) Author() *User {
	var user *User
	if err := x.C("users").Find(M{"_id": x.AuthorID}).One(&user); err != nil {
		x.AddError(err)
		return nil
	}

	user.Link(x.Core)
	return user
}

func (x *Post) Validate() bool {

	// title validation
	if len(x.Title) < 5 {
		x.AddError(errors.New("Title too short"))
	}
	// body validation
	if len(x.Body) < 10 {
		x.AddError(errors.New("Body too short"))
	}

	if x.ErrorCount() > 0 {
		return false
	}
	return true
}

func (x *Post) Save() {

	x.ID = NewObjectId()
	x.AuthorID = x.LoggedIn.ID
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	if err := x.C("posts").Insert(x); err != nil {
		log.Printf("error saving post to db : %v", err.Error())
		x.AddError(err)
	}
}
