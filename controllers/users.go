package controllers

import (
	"grustosoftware/core"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	*core.Core
	_mux *mux.Router
}

func NewUserController(v *core.Core) *UserController {
	x := &UserController{
		v,
		mux.NewRouter(),
	}
	x.InitMux()

	return x
}

func (x *UserController) validate(user core.PotentialUser) bool {

	if user.Name == "" || user.Email == "" || user.Password == "" || user.Cpassword == "" {
		x.AddError(errors.New("forget something??? null data in register request"))
		return false
	}
	//password validation
	if user.Password != user.Cpassword {
		x.AddError(errors.New("Passwords do not match dummy"))
	}
	if len(user.Password) < 8 {
		x.AddError(errors.New("Password must be at least 8 characters"))
	}

	// name validation
	if len(user.Name) < 3 {
		x.AddError(errors.New("name must be at least 3 characters"))
	}

	var u core.User
	if err := x.C("users").Find(bson.M{"name": user.Name}).One(&u); err == nil {
		x.AddError(errors.Errorf("username already exists"))
	}

	// email validation
	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		x.AddError(err)
	}

	if x.ErrorCount() > 0 {
		return false
	}
	return true

}

// Create : creates a new user
func (x *UserController) Create(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)

	if err := r.ParseForm(); err != nil {
		x.AddError(err)
	}
	var pu core.PotentialUser

	if err := schema.NewDecoder().Decode(&pu, r.PostForm); err != nil {
		x.AddError(err)
		return
	}

	if ok := x.validate(pu); !ok {
		return
	}

	u := core.NewUser(pu, x.Core)

	u.Save()

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

}

func (x *UserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	x._mux.ServeHTTP(w, r)
}

func (x *UserController) InitMux() {

	x._mux.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)

}
