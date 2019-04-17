package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"sambragge/go-software-solutions/core"
)

type CommentController struct {
	core *core.Core
	r    *mux.Router
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

}

func (x *CommentController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	x.r.ServeHTTP(w, r)
}

func (x *CommentController) InitMux() {
	x.r.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
}
