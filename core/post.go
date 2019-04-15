package core

import (
	. "gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Post struct {
	ID        ObjectId  `bson:"_id"`
	Title     string    `bson:"title"`
	Body      string    `bson:"body"`
	Author    ObjectId  `bson:"_author"`
	core      *Core     `bson:"-"`
	createdAt time.Time `bson:"_createdAt"`
	updatedAt time.Time `bson:"_updatedAt"`
}

func NewPost(title, body string, author ObjectId, con *Core) *Post {
	x := &Post{
		Title:  title,
		Body:   body,
		Author: author,
		core:   con,
	}

	return x
}

func (x *Post) Link(con *Core) {
	x.core = con
}

func (x *Post) Comments() []*Comment {

	var comments []*Comment
	if err := x.core.C("comments").Find(M{"_parent": x.ID}).All(&comments); err != nil {
		x.core.AddError(err)
		return nil
	}

	for _, v := range comments {
		v.Link(x.core)
	}

	return comments
}

func (x *Post) Validate() bool {

	// title validation
	if len(x.Title) < 5 {
		x.core.AddError(NewError("Title too short", http.StatusBadRequest))
	}
	// body validation
	if len(x.Body) < 10 {
		x.core.AddError(NewError("Body too short", http.StatusBadRequest))
	}

	if x.core.ErrorCount() > 0 {
		return false
	}
	return true
}

func (x *Post) Save() {

	x.ID = NewObjectId()
	x.createdAt = time.Now()
	x.updatedAt = time.Now()

	if err := x.core.C("posts").Insert(x); err != nil {
		x.core.AddError(err)
	}
}
