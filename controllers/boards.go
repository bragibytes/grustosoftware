package controllers

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"grustosoftware/core"
	"log"
	"net/http"
)

type BoardController struct {
	*core.Core
	_mux *mux.Router

}

func NewBoardController(c *core.Core)*BoardController{
	bc := &BoardController{
		c,
		mux.NewRouter(),
	}


	bc.initMux()
	return bc
}

func (bc *BoardController) Create(w http.ResponseWriter, r *http.Request){

	log.Print("getting to create method")
	if err := r.ParseForm();err!=nil {
		bc.AddError(err)
		return
	}

	var board core.Board
	if err := schema.NewDecoder().Decode(&board, r.PostForm);err!=nil{
		bc.AddError(err)
		return
	}

	board.Save(bc.C("boards"))
}

func(bc *BoardController) initMux(){
	bc._mux.Methods(http.MethodPost).Path("/").HandlerFunc(bc.Create)
}

func (bc *BoardController) ServeHTTP(w http.ResponseWriter, r *http.Request){
	bc._mux.ServeHTTP(w, r)
}

