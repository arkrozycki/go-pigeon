package main

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

type server struct {
	conf *Config
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = JSONPayloadToProto(body, s.conf)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`:(`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}

}

func Serve(port string, conf *Config, wg *sync.WaitGroup) {
	log.Info().Msg("OK ... REST API")
	s := &server{conf}
	http.Handle("/", s)
	log.Fatal().Err(http.ListenAndServe(port, nil))
	defer wg.Done()
}
