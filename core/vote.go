package core

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"

	. "gopkg.in/mgo.v2/bson"
)

type Vote struct {
	*Core     `bson:"-"`
	ID        ObjectId  `bson:"_id"`
	VoterID   ObjectId  `bson:"_voter"`
	ParentID  ObjectId  `bson:"_parent" json:"_parent"`
	Value     string    `bson:"value" json:"value"` // "up" or "down"
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func (x *Vote) ShowSelf() {
	log.Printf("\n\n\n ID : %v", x.ID)
	log.Printf("VoterID : %v", x.VoterID)
	log.Printf("ParentID : %v", x.ParentID)
	log.Printf("Value : %v", x.Value)
	log.Printf("CreatedAt : %v", x.CreatedAt)
	log.Printf("UpdatedAt : %v \n\n\n", x.UpdatedAt)
}

// get all votes from parent
func (x *Vote) ExistingVote() (bool, *Vote) {
	var vote *Vote

	if err := x.C("votes").Find(M{"_parent": x.ParentID}).One(&vote); err != nil {
		return false, nil
	}

	vote.Link(x.Core)
	return true, vote
}

func (v *Vote) UpdateValue(val string) {
	if err := v.C("votes").Update(M{"_id": v.ID}, M{"$set": M{"value": val}}); err != nil {
		v.AddError(err)
		return
	}
}

func (x *Vote) Link(core *Core) {
	x.Core = core
}

func (v *Vote) Save(c *mgo.Collection) {
	v.ParentID = ObjectId(v.ParentID)

	v.ID = NewObjectId()

	v.VoterID = v.LoggedIn.ID

	v.CreatedAt = time.Now()

	v.UpdatedAt = time.Now()

	if err := c.Insert(v); err != nil {
		v.AddError(err)
	}
}
