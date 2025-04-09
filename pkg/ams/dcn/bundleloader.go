package dcn

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BundleLoader struct {
	DCNChannel         chan DcnContainer
	AssignmentsChannel chan Assignments
	errHandler         []func(error)
	lastEtag           string
	client             *http.Client
	url                *url.URL
	ticker             time.Ticker
}

func NewBundleLoader(targetURL *url.URL, client *http.Client, ticker time.Ticker) *BundleLoader {

	result := BundleLoader{
		DCNChannel:         make(chan DcnContainer),
		AssignmentsChannel: make(chan Assignments),
		errHandler:         []func(error){},
		client:             client,
		url:                targetURL,
		ticker:             ticker,
	}
	go result.start()
	return &result
}

func (b *BundleLoader) RegisterErrorHandler(handler func(error)) {
	b.errHandler = append(b.errHandler, handler)
}

func (b *BundleLoader) handleError(err error) {
	for _, handler := range b.errHandler {
		handler(err)
	}
}

func (b *BundleLoader) start() {
	b.bundleRequest()

	for {
		<-b.ticker.C
		b.bundleRequest()
	}
}

func (b *BundleLoader) bundleRequest() {
	dcn := DcnContainer{
		Policies:  []Policy{},
		Schemas:   []Schema{},
		Functions: []Function{},
		Tests:     []Test{},
	}
	assignments := Assignments{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    b.url,
		Header: http.Header{
			"If-None-Match": []string{b.lastEtag},
		},
	}

	resp, err := b.client.Do(req)
	if err != nil {
		b.handleError(err)
		return
	}
	if resp.StatusCode == http.StatusNotModified {
		b.handleError(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body string
		bodylen := 1024
		if resp.ContentLength < 1024 {
			bodylen = int(resp.ContentLength)
		}
		bodyBytes := make([]byte, bodylen)
		_, err := resp.Body.Read(bodyBytes)
		if err == nil || err == io.EOF {
			body = string(bodyBytes)
		}

		b.handleError(fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body))
		return
	}
	b.lastEtag = resp.Header.Get("ETag")

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		b.handleError(err)
		return
	}

	defer gz.Close()

	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			b.handleError(err)
			return
		}

		// If it's a regular file, read the content
		if header.Typeflag == tar.TypeReg {
			if strings.HasSuffix(header.Name, ".dcn") {
				content := make([]byte, header.Size)
				_, err := tarReader.Read(content)
				if err != nil && err != io.EOF {
					b.handleError(err)
					return
				}
				var dcnPart DcnContainer
				err = json.Unmarshal(content, &dcnPart)
				if err != nil {
					b.handleError(err)
					return
				}
				dcn.Policies = append(dcn.Policies, dcnPart.Policies...)
				dcn.Functions = append(dcn.Functions, dcnPart.Functions...)
				dcn.Schemas = append(dcn.Schemas, dcnPart.Schemas...)
			}
			if header.Name == "data.json" {
				content := make([]byte, header.Size)
				_, err := tarReader.Read(content)
				if err != nil && err != io.EOF {
					b.handleError(err)
					return
				}
				var assignmentsC AssignmentsContainer
				err = json.Unmarshal(content, &assignmentsC)
				if err != nil {
					b.handleError(err)
					return
				}
				assignments = assignmentsC.Assignments
			}
		}
	}

	b.DCNChannel <- dcn
	b.AssignmentsChannel <- assignments

}
