package core

import (
	"log"
	"time"

	"github.com/pkg/errors"
	. "gopkg.in/mgo.v2/bson"
)

type Post struct {
	*Core     `bson:"-"`
	ID        ObjectId  `bson:"_id"`
	Title     string    `bson:"title"`
	Body      string    `bson:"body"`
	ParentID ObjectId	`bson:"_parent"`
	AuthorID  ObjectId  `bson:"_author"`
	Score     int8      `bson:"-"`
	CreatedAt time.Time `bson:"_createdAt"`
	UpdatedAt time.Time `bson:"_updatedAt"`
}

func (p *Post) Link(core *Core) {
	p.Core = core

}

func (p *Post) CalculateScore() {
	var votes []Vote
	if err := p.C("votes").Find(M{"_parent": p.ID}).All(&votes); err != nil {
		p.AddError(errors.New("the error is comming from trying to find votes in the db " + err.Error()))
		return
	}

	for _, v := range votes {
		if v.Value == "up" {
			p.Score++
		} else if v.Value == "down" {
			p.Score--
		}
	}
}

func (x *Post) Comments() []*Comment {

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

func (x *Post) ShowSelf() {
	log.Printf("\n\n\n ID : %v", x.ID)
	log.Printf("Title : %v", x.Title)
	log.Printf("Body : %v", x.Body)
	log.Printf("AuthorID : %v", x.AuthorID)
	log.Printf("CreatedAt : %v", x.CreatedAt)
	log.Printf("UpdatedAt : %v \n\n\n", x.UpdatedAt)
}

func (x *Post) Validate() bool {

	// title validation
	if len(x.Title) < 5 {
		x.AddError(errors.New("Title too short"))
	}
	log.Print("\n=== checked title")
	// body validation
	if len(x.Body) < 10 {
		x.AddError(errors.New("Body too short"))
	}
	log.Print("\n=== checked body")

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
