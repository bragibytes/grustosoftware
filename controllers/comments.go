package controllers

import (
	"grustosoftware/core"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type CommentController struct {
	*core.Core
	r *mux.Router
}

func NewCommentController(core *core.Core) *CommentController {
	x := &CommentController{
		core,
		mux.NewRouter(),
	}

	x.InitMux()
	return x
}

func (x *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, x.Path, http.StatusSeeOther)

	if err := r.ParseForm(); err != nil {
		x.AddError(err)
		return
	}

	var comment core.Comment
	if err := schema.NewDecoder().Decode(&comment, r.PostForm); err != nil {
		x.AddError(err)
		return
	}

	comment.Link(x.Core)
	if ok := comment.Validate(); !ok {
		return
	}
	comment.Save()

}

func (x *CommentController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	x.r.ServeHTTP(w, r)
}

func (x *CommentController) InitMux() {
	x.r.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
}
