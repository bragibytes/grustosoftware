package core

import (
	"time"

	"github.com/pkg/errors"
	. "gopkg.in/mgo.v2/bson"
)

type Comment struct {
	*Core     `bson:"-" json:"-"`
	ID        ObjectId  `bson:"_id" json:"_id"`
	Body      string    `bson:"body" json:"body"`
	Parent    ObjectId  `bson:"_parent" json:"_parent"`
	AuthorID  ObjectId  `bson:"_author" json:"_author"`
	Score     int8      `bson:"-"`
	CreatedAt time.Time `bson:"_createdAt" json:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt" json:"_updatedAt"`
}

func (x *Comment) IDHex() string {

	return x.ID.Hex()
}

func (x *Comment) Link(core *Core) {
	x.Core = core
}

func (c *Comment) CalculateScore() {
	var votes []Vote
	if err := c.C("votes").Find(M{"_parent": c.ID}).All(&votes); err != nil {
		c.AddError(err)
		return
	}

	var score int8 = 0
	for _, v := range votes {
		if v.Value == "up" {
			c.Score++
		} else if v.Value == "down" {
			c.Score--
		}
	}

	c.Score = score
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
		v.CalculateScore()
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
