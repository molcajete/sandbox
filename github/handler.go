package github

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Commit struct {
	Id  string
	Url string
}

type Webhook struct {
	Ref        string `json:"ref"`
	Id         string `json:"id"`
	Commits    []Commit
	Repository struct {
		Name           string `json:"name"`
		Url            string `json:"url"`
		Default_branch string `json:"default_branch"`
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	secret := []byte(os.Getenv("GITHUB_TOKEN"))
	hook, err := Parse(secret, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	wh := new(Webhook)
	decoder := json.NewDecoder(bytes.NewReader(hook.Payload))
	err = decoder.Decode(wh)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(wh); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Print(b)

	if !strings.HasSuffix(wh.Ref, wh.Repository.Default_branch) {
		http.Error(w,
			http.StatusText(http.StatusPreconditionFailed),
			http.StatusPreconditionFailed,
		)
		return
	}
	// everything seems to be OK so let's try to deploy
	w.WriteHeader(http.StatusAccepted)
}
