package core

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Session struct {
	ID        bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId `bson:"_user"`
	Expires   time.Time     `bson:"expires"`
	core      *Core
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func NewSession(id bson.ObjectId, exp time.Time, con *Core) *Session {
	x := &Session{
		UserId:  id,
		Expires: exp,
		core:    con,
	}

	return x
}

func (x *Session) Link(con *Core) {
	x.core = con
}

func (x *Session) Save() {

	x.ID = bson.NewObjectId()
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	err := x.core.C("sessions").Insert(x)
	if err != nil {
		x.core.AddError(err)
	}
}
