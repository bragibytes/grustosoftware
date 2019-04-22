package controllers

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sambragge/go-software-solutions/core"
)

type ViewController struct {
	*core.Core
	_mux *mux.Router
}

func NewViewController(core *core.Core) *ViewController {
	x := &ViewController{
		core,
		mux.NewRouter(),
	}
	x.InitMux()

	return x
}

func (x *ViewController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	go x.CheckState(w, r)
	x._mux.ServeHTTP(w, r)
}

func (x *ViewController) Blog(w http.ResponseWriter, r *http.Request) {
	x.View(w, "blog", x.Core)
}
func (x *ViewController) Portfolio(w http.ResponseWriter, r *http.Request) {
	x.View(w, "portfolio", x.Core)
}
func (x *ViewController) Register(w http.ResponseWriter, r *http.Request) {
	x.View(w, "register", x.Core)
}
func (x *ViewController) Home(w http.ResponseWriter, r *http.Request) {
	x.View(w, "home", x.Core)
}

func(x *ViewController) Profile (w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]

	var user *core.User
	if err := x.C("users").Find(bson.M{"name":name}).One(&user); err != nil {
		x.AddError(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user.Link(x.Core)
	x.View(w, "profile", user)
}

func (x *ViewController) InitMux() {
	x._mux.Methods(http.MethodGet).Path("/").HandlerFunc(x.Home)
	x._mux.Methods(http.MethodGet).Path("/portfolio").HandlerFunc(x.Portfolio)
	x._mux.Methods(http.MethodGet).Path("/blog").HandlerFunc(x.Blog)
	x._mux.Methods(http.MethodGet).Path("/profile/{name}").HandlerFunc(x.Profile)
	x._mux.Methods(http.MethodGet).Path("/blog/register").HandlerFunc(x.Register)

	x._mux.Methods(http.MethodPost).Path("/login").HandlerFunc(x.Login)
	x._mux.Methods(http.MethodDelete).Path("/logout").HandlerFunc(x.Logout)
}
