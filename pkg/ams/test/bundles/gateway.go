package test

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
)

//go:embed data/simple.dcn.tar.gz
var SimpleDCN []byte

type BundleGatewayMock struct {
	server *httptest.Server
}

func NewBundleGatewayMock() *BundleGatewayMock {
	result := &BundleGatewayMock{}
	result.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.WriteHeader(200)
		w.Write(SimpleDCN)
	}))

	return result
}

func (b *BundleGatewayMock) GetAuthorizationBundleURL() string {
	return b.server.URL + "/bundle"
}

func (b *BundleGatewayMock) GetAuthorizationInstanceID() string {
	return "test-instance-id"
}

func (b *BundleGatewayMock) GetHttpClient() *http.Client {
	return b.server.Client()
}

func (b *BundleGatewayMock) Close() {
	b.server.Close()
}
