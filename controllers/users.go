package controllers

import (
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"sambragge/go-software-solutions/core"
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
	if err := x.C("users").Find(bson.M{"name":user.Name}).One(&u);err == nil {
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

}

func (x *UserController) Destroy(w http.ResponseWriter, r *http.Request) {

	log.Print("in the user destroy method")

	id := mux.Vars(r)["id"]

	log.Printf("id came in as %v\n", id)
	log.Printf("id is a objectIdHex %v\n", bson.IsObjectIdHex(id))

	objId := bson.ObjectIdHex(id)
	log.Printf("turned the id into an ObjectId and got %v\n", objId)

	err := x.C("users").Remove(bson.D{{"id", objId}})
	if err != nil {
		log.Printf("there was an error deleting the user %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (x *UserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Getting the the user controllers ServeHTTP method with the path as %s and nmethod as %s\n", r.URL.Path, r.Method)
	x._mux.ServeHTTP(w, r)
}

func (x *UserController) InitMux() {

	x._mux.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
	//x.r.Methods(http.MethodGet).Path("/").HandlerFunc(x.GetAll)
	x._mux.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(x.Destroy)

}
