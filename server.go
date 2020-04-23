package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

// server
type server struct {
	conf *Config
}

// ServeHTTP inbound handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// get message and routing key
	var temp map[string]interface{}
	json.Unmarshal([]byte(body), &temp)

	// check we have all the json keys required
	if temp["routingKey"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`routingKey required`))
		return
	}

	if temp["proto"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`proto required`))
		return
	}

	if temp["msg"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`msg required`))
		return
	}

	// Marshal just the msg payload back into JSON
	msg, err := json.Marshal(temp["msg"])

	// emit the message to the bus, whoop whoop
	err = Emit(temp["routingKey"].(string), temp["proto"].(string), msg, s.conf)

	// all good, or all bad
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`:(`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}

}

// Serve
func Serve(port string, conf *Config, wg *sync.WaitGroup) {

	s := &server{conf}
	http.Handle("/", s)
	log.Fatal().Err(http.ListenAndServe(port, nil))
	defer wg.Done()
}
