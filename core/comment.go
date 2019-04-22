package core

import (
	"github.com/pkg/errors"
	. "gopkg.in/mgo.v2/bson"
	"time"
)

type Comment struct {
	ID        ObjectId `bson:"_id" json:"_id"`
	Body      string   `bson:"body" json:"body"`
	Parent    ObjectId `bson:"_parent" json:"_parent"`
	AuthorID  ObjectId `bson:"_author" json:"_author"`
	*Core     `bson:"-" json:"-"`
	CreatedAt time.Time `bson:"_createdAt" json:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt" json:"_updatedAt"`
}

func (x *Comment) IDHex() string {
	return x.ID.Hex()
}

//func NewComment(body string, core *Core) *Comment {
//	x := &Comment{
//
//		Body: body,
//		Core: core,
//	}
//
//	return x
//}

func (x *Comment) Link(core *Core) {
	x.Core = core
}

func (x *Comment) Validate() bool {
	if len(x.Body) < 3 {
		x.AddError(errors.New("comment must be at least 3 characters long"))
	}

	if x.ErrorCount() > 0 {
		return false
	}
	return true
}

func (x *Comment) Save() {

	x.ID = NewObjectId()
	x.AuthorID = x.LoggedIn.ID
	x.Parent = ObjectId(x.Parent)
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	if err := x.C("comments").Insert(x); err != nil {
		x.AddError(err)
	}
}

func (x *Comment) CommentCount() int {
	var comments []*Comment
	if err := x.C("comments").Find(M{"_parent": x.ID}).All(&comments); err != nil {
		x.AddError(err)
		return 0
	}
	return len(comments)
}

func (x *Comment) Comments() []*Comment {
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

func (x *Comment) Author() *User {
	var author *User
	if err := x.C("users").Find(M{"_id": x.AuthorID}).One(&author); err != nil {
		x.AddError(err)
		return nil
	}
	return author
}
