package core

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type PotentialUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}

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
	if _, err := w.Write([]byte(r)); err != nil {
		log.Print("error writing json")
	}
}

func decodeUser(body io.Reader) (PotentialUser, error) {
	var pu PotentialUser

	log.Printf("going to try to decode the user json into a go struct. body is %s\n\n", body)
	if err := json.NewDecoder(body).Decode(&pu); err != nil {
		log.Print("\n\nerror decoding json into struct\n\n")
		return PotentialUser{}, err
	}
	return pu, nil
}
