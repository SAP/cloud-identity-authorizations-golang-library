package test

import (
	"net/http"
	"testing"
)

func TestBundleGatewayMock(t *testing.T) {
	mock := NewBundleGatewayMock()
	defer mock.server.Close()

	client := mock.GetHttpClient()

	resp, err := client.Get(mock.GetAuthorizationBundleURL() + "/" + mock.GetAuthorizationInstanceID() + ".dcn.tar.gz")

	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
}
