package main

import (
	"net/http"
	"os"

	"github.com/sap/cloud-identity-authorizations-golang-library/internal/testserver"
)

func main() {
	port, ok := os.LookupEnv("DCN_TEST_SERVER_PORT")
	if !ok {
		port = "8085"
	}

	r := testserver.Router{}

	http.ListenAndServe(":"+port, r.Mux())

}
