package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sambragge/go-software-solutions/core"
)

type response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

func jsonResponse(w http.ResponseWriter, success bool, payload interface{}) {

	res := response{
		Success: success,
		Payload: payload,
	}

	r, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(r))
}

func decodeUser(body io.Reader) (*core.PotentialUser, error) {
	var pu *core.PotentialUser

	log.Printf("going to try to decode the user json into a go struct. body is %s\n\n", body)
	err := json.NewDecoder(body).Decode(&pu)
	if err != nil {
		log.Printf("there was an error decoding : %v\n\n", err.Error())
		return nil, err
	}

	return pu, nil
}
