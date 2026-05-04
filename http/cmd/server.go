package main

import (
	"net/http"
	"os"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/http/logging"
	"github.com/sap/cloud-identity-authorizations-golang-library/http/server"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-security-client-go/env"
)

const envDCNPath = "AMS_DCN_ROOT"

func main() {
	var am *ams.AuthorizationManager
	var err error
	l := logging.PlainLogger{}
	if os.Getenv(envDCNPath) != "" {
		am = ams.NewAuthorizationManagerForFs(os.Getenv(envDCNPath), nil)
	} else {
		config, err := env.ParseIdentityConfig()
		if err != nil {
			panic(err)
		}
		am, err = ams.NewAuthorizationManagerForIAS(
			config.GetAuthorizationBundleURL(),
			config.GetAuthorizationInstanceID(),
			config.GetCertificate(),
			config.GetKey(),
			l,
		)
		// am, err = ams.NewAuthorizationManagerForIASConfig(
		// 	config,
		// 	l,
		// )

		if err != nil {
			panic(err)
		}
	}
	router := server.NewRouter(am, l)

	srv := &http.Server{
		Addr:         ":8099",
		Handler:      router.Mux(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
