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
	errorHandler       []func(error)
}

func NewBundleLoader(

	targetURL *url.URL,
	client *http.Client,
	ticker time.Ticker,
	errorCallback func(error),
) *BundleLoader {

	result := BundleLoader{
		DCNChannel:         make(chan DcnContainer),
		AssignmentsChannel: make(chan Assignments),
		client:             client,
		url:                targetURL,
		ticker:             ticker,
		errorHandler:       []func(error){},
	}
	if errorCallback != nil {
		result.errorHandler = append(result.errorHandler, errorCallback)
	}

	return &result
}

func (b *BundleLoader) handleError(err error) {
	for _, handler := range b.errorHandler {
		handler(err)
	}
}

func (b *BundleLoader) Run(ctx context.Context) {
	go func() {
		retries := 20
		for range retries {
			err := b.BundleRequest(ctx, b.client)
			if err == nil {
				break
			}
			b.handleError(err)
			time.Sleep(time.Second * 1)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-b.ticker.C:
				err := b.BundleRequest(ctx, b.client)
				if err != nil {
					b.handleError(err)
				}
			}
		}
	}()
}

func (b *BundleLoader) SetHttpClient(ctx context.Context, client *http.Client) error {
	err := b.BundleRequest(ctx, client)
	if err != nil {
		return err
	}
	b.client = client
	return nil
}

func (b *BundleLoader) BundleRequest(ctx context.Context, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.url.String(), nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"If-None-Match": []string{b.lastEtag},
		"User-Agent":    []string{fmt.Sprintf("golang-dcn-%s", version)},
	}

	resp, err := client.Do(req)
	if err != nil {
		b.handleError(err)
		return err
	}
	if resp.StatusCode == http.StatusNotModified {
		return nil
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
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body)
	}
	b.lastEtag = resp.Header.Get("ETag")

	dcn, assignments, err := ReadBundleTarGz(resp.Body)
	if err != nil {
		return err
	}
	b.DCNChannel <- dcn
	b.AssignmentsChannel <- assignments
	return nil
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
				if dcnPart.Version > DCNVERSION || dcnPart.Version == 0 {
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
