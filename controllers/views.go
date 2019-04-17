package controllers

import (
	"github.com/gorilla/mux"
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

func (x *ViewController) BlogPage(w http.ResponseWriter, r *http.Request)      { x.View(w) }
func (x *ViewController) TeamPage(w http.ResponseWriter, r *http.Request)      { x.View(w) }
func (x *ViewController) PortfolioPage(w http.ResponseWriter, r *http.Request) { x.View(w) }
func (x *ViewController) ContactPage(w http.ResponseWriter, r *http.Request)   { x.View(w) }
func (x *ViewController) RegisterPage(w http.ResponseWriter, r *http.Request)  { x.View(w) }
func (x *ViewController) ProfilePage(w http.ResponseWriter, r *http.Request)   { x.View(w) }
func (x *ViewController) SplashPage(w http.ResponseWriter, r *http.Request)    { x.View(w) }

func (x *ViewController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		x.Path = r.URL.Path
	}

	x.CheckState(w, r)

	x._mux.ServeHTTP(w, r)
}

func (x *ViewController) InitMux() {

	x._mux.Methods(http.MethodGet).Path("/").HandlerFunc(x.SplashPage)
	x._mux.Methods(http.MethodGet).Path("/blog").HandlerFunc(x.BlogPage)
	x._mux.Methods(http.MethodGet).Path("/team").HandlerFunc(x.TeamPage)
	x._mux.Methods(http.MethodGet).Path("/portfolio").HandlerFunc(x.PortfolioPage)
	x._mux.Methods(http.MethodGet).Path("/contact").HandlerFunc(x.ContactPage)
	x._mux.Methods(http.MethodGet).Path("/register").HandlerFunc(x.RegisterPage)
	x._mux.Methods(http.MethodGet).Path("/profile").HandlerFunc(x.ProfilePage)

	// Session
	x._mux.Methods(http.MethodPost).Path("/login").HandlerFunc(x.Login)
	x._mux.Methods(http.MethodDelete).Path("/logout").HandlerFunc(x.Logout)

}
