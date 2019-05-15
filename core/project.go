package core

import (
	"github.com/pkg/errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	*Core       `bson:"-"`
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	ImageLink   string        `bson:"imageLink"`
	Description string        `bson:"description"`
	ProjectLink string        `bson:"link"`
	CodeLink    string        `bson:"codeLink"`
	AuthorID    bson.ObjectId `bson:"_author"`
	CreatedAt   time.Time     `bson:"_createdAt"`
	UpdatedAt   time.Time     `bson:"_updatedAt"`
}

func (x *Project) Link(c *Core) {
	x.Core = c
}

func (x *Project) Validate() bool {

	if len(x.Name) < 1 {
		x.AddError(errors.New("Project name is required"))
	}
	if len(x.Description) < 1 {
		x.AddError(errors.New("Project description is required"))
	}

	if x.ErrorCount() > 0 {
		return false
	}
	return true
}

func (x *Project) Save() {
	x.ID = bson.NewObjectId()
	x.AuthorID = x.LoggedIn.ID
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	if err := x.C("projects").Insert(x); err != nil {
		x.AddError(err)
	}
}
