package healthcheck

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/austinjalexander/util/pkg/server"
)

const path = "/api/healthcheck"

var (
	headers     = []string{}
	methods     = []string{"GET"}
	queryParams = []string{}
)

type data struct {
	Status string `json:"status"`
}

// Configure configures a new handler and returns it.
func Configure() server.Handler {
	return server.Handler{
		Func:        handlerFunc,
		Path:        path,
		Headers:     headers,
		Methods:     methods,
		QueryParams: queryParams,
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	resp := server.JSONresponse{
		Data: data{
			Status: "OK",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("error: marshaling response: %s", err), http.StatusInternalServerError)
	} else {
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}
