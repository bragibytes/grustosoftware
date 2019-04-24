package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sambragge/go-software-solutions/core"
)

type VoteController struct {
	*core.Core
	_mux *mux.Router
}

func NewVoteController(core *core.Core) *VoteController {
	x := &VoteController{
		core,
		mux.NewRouter(),
	}

	x.InitMux()
	return x
}


func (x *VoteController) Create(w http.ResponseWriter, r *http.Request){

	var vote *core.Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		x.AddError(err)
		log.Printf("\n\n did not get the vote ")
		return
	}

	log.Printf("\n\n got the vote ")

	vote.ShowSelf()

	vote.Link(x.Core)
	vote.Save()
}


func (x *VoteController) ServeHTTP(w http.ResponseWriter, r *http.Request){

	x._mux.ServeHTTP(w, r)
}

func(x *VoteController) InitMux(){

	x._mux.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
}

