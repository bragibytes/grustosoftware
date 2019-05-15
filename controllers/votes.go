package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"grustosoftware/core"
	"net/http"
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

func (x *VoteController) Create(w http.ResponseWriter, r *http.Request) {

	var vote *core.Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		x.AddError(err)
		return
	}

	vote.Link(x.Core)

	existing, v := vote.ExistingVote()
	if existing {
		v.UpdateValue(vote.Value)
	} else {
		vote.Save(x.C("votes"))
	}

}

func (x *VoteController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	x._mux.ServeHTTP(w, r)
}

func (x *VoteController) InitMux() {

	x._mux.Methods(http.MethodPost).Path("/").HandlerFunc(x.Create)
}
