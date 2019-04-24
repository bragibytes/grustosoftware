package controllers

import (
	"github.com/gorilla/mux"
	"sambragge/go-software-solutions/core"
)

type VoteController struct {
	*core.Core
	_mux *mux.Router
}



