package controllers

import (
	"grustosoftware/core"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type PostController struct {
	*core.Core
	_mux *mux.Router
}

func NewPostController(mc *core.Core) *PostController {
	pc := &PostController{
		mc,
		mux.NewRouter(),
	}
	pc.InitMux()

	return pc
}

func (pc *PostController) Create(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, pc.Path, http.StatusSeeOther)

	if err := r.ParseForm(); err != nil {
		pc.AddError(err)
		return
	}

	var post core.Post
	if err := schema.NewDecoder().Decode(&post, r.PostForm); err != nil {
		pc.AddError(err)
		return
	}

	post.Link(pc.Core)

	if ok := post.Validate(); !ok {
		return
	}

	post.Save()

}

func (pc *PostController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pc._mux.ServeHTTP(w, r)
}

func (pc *PostController) InitMux() {
	pc._mux.Methods(http.MethodPost).Path("/").HandlerFunc(pc.Create)
}
