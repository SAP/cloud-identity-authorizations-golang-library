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
	ctx                context.Context
	DCNChannel         chan DcnContainer
	AssignmentsChannel chan Assignments
	lastEtag           string
	client             *http.Client
	url                *url.URL
	ticker             time.Ticker
	l                  logging.Logger
	closed             chan bool
	cancel             context.CancelFunc
}

func NewBundleLoader(
	ctx context.Context,
	targetURL *url.URL,
	client *http.Client,
	ticker time.Ticker,
	log logging.Logger,
) *BundleLoader {
	ctx, cancel := context.WithCancel(ctx)
	result := BundleLoader{
		ctx:                ctx,
		cancel:             cancel,
		DCNChannel:         make(chan DcnContainer),
		AssignmentsChannel: make(chan Assignments),
		client:             client,
		url:                targetURL,
		ticker:             ticker,
		l:                  log,
		closed:             make(chan bool),
	}

	go result.start()
	return &result
}

func (b *BundleLoader) handleError(err error) {
	if b.l != nil {
		b.l.Errorf(b.ctx, "%v", err)
	}
}

func (b *BundleLoader) start() {
	b.bundleRequest()

	for {
		select {
		case <-b.closed:
			b.closed <- true
			return
		case <-b.ticker.C:
			b.bundleRequest()
		}
	}
}

func (b *BundleLoader) Close(ctx context.Context) error {
	b.ticker.Stop()
	b.cancel()

	select {
	case b.closed <- true:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (b *BundleLoader) bundleRequest() {
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

	dcn, assignments, err := ReadBundleTarGz(resp.Body)
	if err != nil {
		b.handleError(err)
		return
	}

	b.DCNChannel <- dcn
	b.AssignmentsChannel <- assignments
}

func ReadBundleTarGz(reader io.Reader) (DcnContainer, Assignments, error) {
	dcn := DcnContainer{
		Policies:  []Policy{},
		Schemas:   []Schema{},
		Functions: []Function{},
		Tests:     []Test{},
	}
	assignments := Assignments{}

	gz, err := gzip.NewReader(reader)
	if err != nil {
		return DcnContainer{}, nil, err
	}

	defer gz.Close()

	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return DcnContainer{}, nil, err
		}

		// If it's a regular file, read the content
		if header.Typeflag == tar.TypeReg {
			if strings.HasSuffix(header.Name, ".dcn") {
				content := make([]byte, header.Size)
				_, err := io.ReadFull(tarReader, content)
				if err != nil {
					return DcnContainer{}, nil, err
				}
				var dcnPart DcnContainer
				err = json.Unmarshal(content, &dcnPart)
				if err != nil {
					return DcnContainer{}, nil, err
				}
				if dcnPart.Version > DCNVERSION {
					return DcnContainer{}, nil, fmt.Errorf(
						"incompatible DCN version: bundle has version %d but loader supports up to %d",
						dcnPart.Version,
						DCNVERSION,
					)
				}
				dcn.Policies = append(dcn.Policies, dcnPart.Policies...)
				dcn.Functions = append(dcn.Functions, dcnPart.Functions...)
				dcn.Schemas = append(dcn.Schemas, dcnPart.Schemas...)
			}
			if header.Name == "data.json" {
				content := make([]byte, header.Size)
				_, err := io.ReadFull(tarReader, content)
				if err != nil {
					return DcnContainer{}, nil, err
				}
				var assignmentsC AssignmentsContainer
				err = json.Unmarshal(content, &assignmentsC)
				if err != nil {
					return DcnContainer{}, nil, err
				}
				assignments = assignmentsC.Assignments
			}
		}
	}
	return dcn, assignments, nil
}
