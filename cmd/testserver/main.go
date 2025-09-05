package main

import (
	"net/http"
	"os"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/internal/testserver"
)

func main() {
	port, ok := os.LookupEnv("DCN_TEST_SERVER_PORT")
	if !ok {
		port = "8085"
	}

	r := testserver.Router{}

	s := http.Server{
		Addr:         ":" + port,
		Handler:      r.Mux(),
		ReadTimeout:  10 * time.Second, // 10 secods should be enough to upload a DCN
		WriteTimeout: 10 * time.Second,
	}
	defer s.Close()
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}
