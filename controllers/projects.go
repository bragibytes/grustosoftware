package controllers

import (
	"github.com/gorilla/schema"
	"grustosoftware/core"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectController struct {
	*core.Core
	_mux *mux.Router
}

func NewProjectController(c *core.Core) *ProjectController {
	x := &ProjectController{
		c,
		mux.NewRouter(),
	}

	x.InitMux()
	return x
}

func (x *ProjectController) Create(w http.ResponseWriter, r *http.Request) {

	log.Print("\n Getting to the Project Create method \n")
	defer http.Redirect(w, r, x.Path, http.StatusSeeOther)

	if err := r.ParseForm(); err != nil {
		x.AddError(err)
		return
	}

	var project core.Project
	if err := schema.NewDecoder().Decode(&project, r.PostForm); err != nil {
		x.AddError(err)
		return
	}

	project.Link(x.Core)

	if ok := project.Validate(); !ok {
		return
	}

	project.Save()

}

func (x *ProjectController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x._mux.ServeHTTP(w, r)
}

func (x *ProjectController) InitMux() {
	x._mux.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
}
