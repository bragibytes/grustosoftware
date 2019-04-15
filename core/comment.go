package core

import (
	. "gopkg.in/mgo.v2/bson"
	"time"
)

type Comment struct {
	ID        ObjectId  `bson:"_id"`
	Body      string    `bson:"body"`
	Parent    ObjectId  `bson:"_parent"`
	Author    ObjectId  `bson:"_author"`
	core      *Core     `bson:"-"`
	createdAt time.Time `bson:"_createdAt"`
	updatedAt time.Time `bson:"_updatedAt"`
}

func NewComment(body string, core *Core) *Comment {
	x := &Comment{

		Body: body,
		core: core,
	}

	return x
}

func (x *Comment) Link(core *Core) {
	x.core = core
}

func (x *Comment) Comments() []*Comment {
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
