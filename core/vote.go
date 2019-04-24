package core

import (
	. "gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Vote struct {
	*Core `bson:"-"`
	ID ObjectId `bson:"_id"`
	VoterID ObjectId `bson:"_voter"`
	ParentID ObjectId `bson:"_parent" json:"_parent"`
	Value string `bson:"value" json:"value"` // "up" or "down"
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}
func (x *Vote) ShowSelf(){
	log.Printf("\n\n\n ID : %v", x.ID)
	log.Printf("VoterID : %v", x.VoterID)
	log.Printf("ParentID : %v", x.ParentID)
	log.Printf("Value : %v", x.Value)
	log.Printf("CreatedAt : %v", x.CreatedAt)
	log.Printf("UpdatedAt : %v \n\n\n", x.UpdatedAt)
}

func (v *Vote) Link(core *Core){
	v.Core = core
}

func(v *Vote) Save(){
	v.ParentID = ObjectId(v.ParentID)

	v.ID = NewObjectId()

	v.VoterID = v.LoggedIn.ID

	v.CreatedAt = time.Now()

	v.UpdatedAt = time.Now()

	v.ShowSelf()

	if err := v.C("votes").Insert(v); err != nil {
		v.AddError(err)
	}
}

