package core

import (
	"gopkg.in/mgo.v2"
	. "gopkg.in/mgo.v2/bson"
	"time"
)

type Board struct {
	*Core `bson:"-"`
	ID ObjectId `bson:"_id"`
	Topic string
	Description string
	AuthorID ObjectId `bson:"_author"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func(b *Board) Link(c *Core){
	b.Core = c
}

func (b *Board) Posts() []*Post {
	var posts []*Post

	if err := b.C("posts").Find(M{"_parent":b.ID}).All(&posts);err!=nil {
		b.AddError(err)
		return nil
	}


	return posts
}

func (b *Board) Author()*User{
	var author *User

	if err := b.C("users").Find(M{"_id":b.AuthorID}).One(&author);err!=nil{
		b.AddError(err)
		return nil
	}

	return author
}



func(b *Board)Save(c *mgo.Collection){
	b.ID = NewObjectId()
	b.AuthorID = b.LoggedIn.ID
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	if err := c.Insert(b);err!=nil{
		b.AddError(err)
	}
}