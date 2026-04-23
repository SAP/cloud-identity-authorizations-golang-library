package dcn

import (
	"archive/tar"
	"compress/gzip"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/logging"
)

//go:embed VERSION
var version string

const DCNVERSION = 1

type BundleLoader struct {
	DCNChannel         chan DcnContainer
	AssignmentsChannel chan Assignments
	lastEtag           string
	client             *http.Client
	url                *url.URL
	ticker             time.Ticker
	l                  logging.Logger
}

func NewBundleLoader(targetURL *url.URL,
	client *http.Client,
	ticker time.Ticker,
	log logging.Logger,
) *BundleLoader {
	result := BundleLoader{
		DCNChannel:         make(chan DcnContainer),
		AssignmentsChannel: make(chan Assignments),
		client:             client,
		url:                targetURL,
		ticker:             ticker,
		l:                  log,
	}

	go result.start()
	return &result
}

func (b *BundleLoader) handleError(err error) {
	if b.l != nil {
		b.l.Error(context.Background(), fmt.Sprintf("%v", err))
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
			"User-Agent":    []string{fmt.Sprintf("golang-dcn-%s", version)},
		},
	}

	resp, err := b.client.Do(req)
	if err != nil {
		b.handleError(err)
		return
	}
	if resp.StatusCode == http.StatusNotModified {
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
				_, err := io.ReadFull(tarReader, content)
				if err != nil {
					b.handleError(err)
					return
				}
				var dcnPart DcnContainer
				err = json.Unmarshal(content, &dcnPart)
				if err != nil {
					b.handleError(err)
					return
				}
				if dcnPart.Version > DCNVERSION {
					b.handleError(fmt.Errorf(
						"incompatible DCN version: bundle has version %d but loader supports up to %d",
						dcnPart.Version,
						DCNVERSION,
					))
					return
				}
				dcn.Policies = append(dcn.Policies, dcnPart.Policies...)
				dcn.Functions = append(dcn.Functions, dcnPart.Functions...)
				dcn.Schemas = append(dcn.Schemas, dcnPart.Schemas...)
			}
			if header.Name == "data.json" {
				content := make([]byte, header.Size)
				_, err := io.ReadFull(tarReader, content)
				if err != nil {
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
