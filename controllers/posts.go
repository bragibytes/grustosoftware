package controllers

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"sambragge/go-software-solutions/core"
)

type PostController struct {
	core *core.Core
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

	if err := r.ParseForm(); err != nil {
		pc.core.AddError(err)
		return
	}
	var post core.Post
	if err := schema.NewDecoder().Decode(&post, r.PostForm); err != nil {
		pc.core.AddError(err)
		return
	}
	post.Link(pc.core)

	if ok := post.Validate(); !ok {
		return
	}

	post.Save()

}

func (pc *PostController) InitMux() {
	pc._mux.Methods(http.MethodPost).Path("/").HandlerFunc(pc.Create)
}

func (pc *PostController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pc._mux.ServeHTTP(w, r)
}
