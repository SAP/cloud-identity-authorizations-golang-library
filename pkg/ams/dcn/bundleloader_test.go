package dcn

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

//go:embed bundles/original_bundle.tar.gz
var bundle []byte

//go:embed bundles/big_data_json.tar.gz
var bigDataJson []byte

const testetag = "test-etag"

func TestBundleLoader(t *testing.T) { //nolint:maintidx
	var recordedRequests []http.Request

	serveBundle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.Header().Set("Etag", testetag)
		if r.Header.Get("If-None-Match") == testetag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bundle) //nolint:errcheck
	})

	serveBigDataJSONBundle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.Header().Set("Etag", testetag)
		if r.Header.Get("If-None-Match") == testetag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bigDataJson) //nolint:errcheck
	})

	serveError := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.WriteHeader(http.StatusInternalServerError)
		res := []byte("Internal Server Error")
		w.Write(res) //nolint:errcheck
	})

	serveErrorNoBody := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.WriteHeader(http.StatusInternalServerError)
	})

	serveNonGzip := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.Header().Set("Etag", testetag)
		w.Write([]byte("asdf")) //nolint:errcheck
	})

	serveNonTar := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.Header().Set("Etag", "test-etag")
		gzip.NewWriter(w).Write([]byte("asdf")) //nolint:errcheck
	})

	serveUnparseableDCN := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz := gzip.NewWriter(w)
		defer gz.Close()
		tarWriter := tar.NewWriter(gz)
		defer tarWriter.Close()
		content := []byte("this is not a dcn")
		err := tarWriter.WriteHeader(&tar.Header{
			Name: "broken.dcn",
			Size: int64(len(content)),
		})
		if err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		_, err = tarWriter.Write(content)
		if err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	})

	serveUnparseableDataJSON := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz := gzip.NewWriter(w)
		defer gz.Close()
		tarWriter := tar.NewWriter(gz)
		defer tarWriter.Close()
		content := []byte("this is not a dcn")
		err := tarWriter.WriteHeader(&tar.Header{
			Name: "data.json",
			Size: int64(len(content)),
		})
		if err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		_, err = tarWriter.Write(content)
		if err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	})

	serveBrokenDataJSON := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := bytes.Buffer{}
		gz := gzip.NewWriter(w)
		defer gz.Close()
		tarWriter := tar.NewWriter(&result)
		defer tarWriter.Close()
		err := tarWriter.WriteHeader(&tar.Header{
			Name: "data.json",
			Size: 10,
		})
		if err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		tarWriter.Write([]byte("asdfghjkla")) //nolint:errcheck
		truncated := result.Bytes()[:len(result.Bytes())-2]
		gz.Write(truncated) //nolint:errcheck
	})

	serveBrokenDCN := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := bytes.Buffer{}
		gz := gzip.NewWriter(w)
		defer gz.Close()
		tarWriter := tar.NewWriter(&result)
		defer tarWriter.Close()
		err := tarWriter.WriteHeader(&tar.Header{
			Name: "broken.dcn",
			Size: 10,
		})
		if err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		tarWriter.Write([]byte("asdfghjkla")) //nolint:errcheck
		truncated := result.Bytes()[:len(result.Bytes())-2]
		gz.Write(truncated) //nolint:errcheck
	})

	tickerC := make(chan time.Time)
	ticker := time.Ticker{
		C: tickerC,
	}

	t.Run("Normal case", func(t *testing.T) {
		ts := httptest.NewServer(serveBundle)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		bundleLoader := NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, nil)

		gotDCN := <-bundleLoader.DCNChannel
		gotAssignments := <-bundleLoader.AssignmentsChannel

		if len(gotDCN.Policies) != 6 {
			t.Fatalf("expected 6 policies, got %d", len(gotDCN.Policies))
		}
		if gotAssignments == nil {
			t.Fatalf("expected assignments to be non-nil")
		}
		if len(gotAssignments) == 0 {
			t.Fatalf("expected assignments to be non-empty")
		}
		_, ok := gotAssignments["cb0bd09b-f94b-4678-b172-b82e5d42fb37"]
		if !ok {
			t.Fatalf("expected assignment to be present")
		}
		if len(gotDCN.Schemas) != 1 {
			t.Fatalf("expected 1 schema, got %d", len(gotDCN.Schemas))
		}

		if len(recordedRequests) != 1 {
			t.Fatalf("expected 1 request, got %d", len(recordedRequests))
		}

		tickerC <- time.Now()
		time.Sleep(time.Millisecond)
		if len(recordedRequests) != 2 {
			t.Fatalf("expected 2 request, got %d", len(recordedRequests))
		}
		want := "test-etag"
		got := recordedRequests[1].Header.Get("If-None-Match")
		if got != want {
			t.Fatalf("expected If-None-Match header to be '%s', got '%s'", want, got)
		}

		want = fmt.Sprintf("golang-dcn-%s", version)
		got = recordedRequests[1].Header.Get("User-Agent")
		if got != want {
			t.Fatalf("expected User-Agent header to be '%s', got '%s'", want, got)
		}
	})

	t.Run("Big data JSON bundle", func(t *testing.T) {
		ts := httptest.NewServer(serveBigDataJSONBundle)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		bundleLoader := NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		select {
		case <-ml.errorsReceived:
			t.Fatalf("expected no error, got an error")
		case <-bundleLoader.DCNChannel:
			assignments := <-bundleLoader.AssignmentsChannel
			if len(assignments) < 200 {
				t.Fatalf("expected at least 200 assignments, got %d", len(assignments))
			}
		}
	})

	t.Run("broken server", func(t *testing.T) {
		ml := newErrorHandler()
		ts := httptest.NewServer(serveError)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("broken server no body", func(t *testing.T) {
		ts := httptest.NewServer(serveErrorNoBody)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("broken http client", func(t *testing.T) {
		targetURL, _ := url.Parse("http://127.0.0.1:1234")

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, &http.Client{}, ticker, ml.Callback)

		time.Sleep(2 * time.Millisecond)
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("server not gzip", func(t *testing.T) {
		ts := httptest.NewServer(serveNonGzip)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}

	})

	t.Run("server not tar", func(t *testing.T) {
		ts := httptest.NewServer(serveNonTar)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("unparseable dcn", func(t *testing.T) {
		ts := httptest.NewServer(serveUnparseableDCN)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("unparseable data.json", func(t *testing.T) {
		ts := httptest.NewServer(serveUnparseableDataJSON)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		ml := newErrorHandler()

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("broken data.json filebody", func(t *testing.T) {
		ml := newErrorHandler()
		ts := httptest.NewServer(serveBrokenDataJSON)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("broken dcn filebody", func(t *testing.T) {
		ml := newErrorHandler()
		ts := httptest.NewServer(serveBrokenDCN)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		NewBundleLoader(context.Background(), targetURL, ts.Client(), ticker, ml.Callback)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})
}
