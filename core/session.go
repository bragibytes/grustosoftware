package core

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Session struct {
	ID        bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId `bson:"_user"`
	Expires   time.Time     `bson:"expires"`
	*Core     `bson:"-"`
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func NewSession(id bson.ObjectId, exp time.Time, core *Core) *Session {
	x := &Session{
		UserId:  id,
		Expires: exp,
		Core:    core,
	}

	return x
}

func (x *Session) Link(core *Core) {
	x.Core = core
}

func (x *Session) Save() error {

	x.ID = bson.NewObjectId()
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	if err := x.C("sessions").Insert(x); err != nil {
		return err
	}

	return nil
}
