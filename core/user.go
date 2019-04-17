package core

import (
	"golang.org/x/crypto/bcrypt"
	. "gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	ID        ObjectId  `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Perm      int       `bson:"perm"` // 0 = normal user, 1 = admin
	Valid     bool      `bson:"valid"`
	*Core     `bson:"-"`
	createdAt time.Time `bson:"_createdAt"`
	updatedAt time.Time `bson:"_updatedAt"`
}

func NewUser(pu PotentialUser, core *Core) *User {
	x := &User{
		Name:     pu.Name,
		Email:    pu.Email,
		Password: pu.Password,
		Core:core,
	}

	return x
}

func (x *User) Link(core *Core) {
	x.Core = core
}

func (x *User) Comments() []*Comment {
	var comments []*Comment
	if err := x.C("comments").Find(M{"_author": x.ID}).All(&comments); err != nil {
		x.AddError(err)
		return nil
	}
	for _, v := range comments {
		v.Link(x.Core)
	}

	return comments
}

func (x *User) Posts() []*Post {
	var posts []*Post
	if err := x.C("posts").Find(M{"_author": x.ID}).All(&posts); err != nil {
		x.AddError(err)
		return nil
	}
	for _, v := range posts {
		v.Link(x.Core)
	}

	return posts
}

func (x *User) ComparePasswordWith(password string) bool {


	err := bcrypt.CompareHashAndPassword([]byte(x.Password), []byte(password))
	if err != nil {
		x.AddError(err)
		return false
	}

	return true
}

func (x *User) Save() {

	// pre save
	x.ID = NewObjectId()
	x.Valid = false
	x.Perm = 0
	x.createdAt = time.Now()
	x.updatedAt = time.Now()

	// encrypt password
	hashed, err := bcrypt.GenerateFromPassword([]byte(x.Password), bcrypt.MinCost)
	if err != nil {
		x.AddError(err)
		return
	}
	x.Password = string(hashed)

	// ready to save`
	err = x.C("users").Insert(x)
	if err != nil {
		x.AddError(err)
		return
	}

}
