package dcn

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

//go:embed bundles/original_bundle.tar.gz
var bundle []byte

func TestBundleLoader(t *testing.T) { //nolint:maintidx
	var recordedRequests []http.Request

	serveBundle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordedRequests = append(recordedRequests, *r)
		w.Header().Set("Etag", "test-etag")
		if r.Header.Get("If-None-Match") == "test-etag" {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bundle) //nolint:errcheck
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
		w.Header().Set("Etag", "test-etag")
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

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)

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
	})

	t.Run("broken server", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveError)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected status code 500: Internal Server Error"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("broken server no body", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveErrorNoBody)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errReceived <- true
			errors = append(errors, err)
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected status code 500: "
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("broken http client", func(t *testing.T) {
		targetURL, _ := url.Parse("http://127.0.0.1:1234")

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, &http.Client{}, ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
		})

		time.Sleep(2 * time.Millisecond)
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}

		errors = []error{}
	})

	t.Run("server not gzip", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveNonGzip)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected EOF" //nolint:goconst
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("server not tar", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveNonTar)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected EOF"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("unparseable dcn", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveUnparseableDCN)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "invalid character 'h' in literal true (expecting 'r')"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("unparseable data.json", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveUnparseableDataJSON)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "invalid character 'h' in literal true (expecting 'r')"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("broken data.json filebody", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveBrokenDataJSON)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected EOF"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})

	t.Run("broken dcn filebody", func(t *testing.T) {
		errReceived := make(chan bool)
		ts := httptest.NewServer(serveBrokenDCN)
		defer ts.Close()

		targetURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %v", err)
		}

		errors := []error{}

		bundleLoader := NewBundleLoader(targetURL, ts.Client(), ticker)
		bundleLoader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
		want := "unexpected EOF"
		got := errors[0].Error()
		if got != want {
			t.Fatalf("expected error to be '%s', got %v", want, got)
		}
		errors = []error{}
	})
}
