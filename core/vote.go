package core

import (
	"github.com/pkg/errors"
	. "gopkg.in/mgo.v2/bson"
	"time"
)

type Vote struct {
	*Core `bson:"-"`
	ID ObjectId `bson:"_id"`
	VoterID ObjectId `bson:"_voter"`
	ParentID ObjectId `bson:"_parent"`
	Value string `bson:"value"` // "up" or "down"
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func (v *Vote) Link(core *Core){
	v.Core = core
}

func (v *Vote) Validate() bool {

	if v.Value != "up" || v.Value != "down" {
		v.AddError(errors.New("Invalid vote value"))
	}

	if v.ErrorCount() > 0 { return false }
	return true
}

func(v *Vote) Save(){
	v.ID = NewObjectId()
	v.VoterID = v.LoggedIn.ID
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	if err := v.C("votes").Insert(v); err != nil {
		v.AddError(err)
	}
}

