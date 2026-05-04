package dcn

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

type mockLogger struct {
	errors         []string
	errorsReceived chan bool
}

func (l *mockLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (l *mockLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (l *mockLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (l *mockLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
	l.errorsReceived <- true
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		errors:         []string{},
		errorsReceived: make(chan bool),
	}
}

func TestLocalLoader(t *testing.T) {
	t.Run("on testfolder", func(t *testing.T) {
		errors := []string{}
		loader := NewLocalLoader("testfolder", &mockLogger{errors: errors})

		dcn := <-loader.DCNChannel
		assignments := <-loader.AssignmentsChannel
		if len(errors) != 0 {
			t.Fatalf("expected 0 errors, got %d", len(errors))
		}
		if len(dcn.Policies) != 1 {
			t.Fatalf("expected 1 policy, got %d", len(dcn.Policies))
		}
		if len(dcn.Schemas) != 1 {
			t.Fatalf("expected 1 schema, got %d", len(dcn.Schemas))
		}

		wantAssignments := Assignments{
			"tenant1": {
				"user1": []string{"cas.Base"},
			},
		}
		if !reflect.DeepEqual(assignments, wantAssignments) {
			t.Fatalf("expected %v, got %v", wantAssignments, assignments)
		}
	})
	t.Run("broken data.json", func(t *testing.T) {
		ml := newMockLogger()

		NewLocalLoader("edgecases/broken-data-json", ml)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("broken DCN file", func(t *testing.T) {
		ml := newMockLogger()

		NewLocalLoader("edgecases/broken-dcn", ml)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("unreadable data.json", func(t *testing.T) {
		ml := newMockLogger()
		tmp := createTempFolderWithUnreadableFile("data.json")
		defer os.RemoveAll(tmp) // Clean up
		NewLocalLoader(tmp, ml)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Errorf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("unreadable DCN", func(t *testing.T) {
		ml := newMockLogger()

		tmp := createTempFolderWithUnreadableFile("x.dcn")
		defer os.RemoveAll(tmp) // Clean up
		NewLocalLoader(tmp, ml)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Errorf("expected 1 request, got %d", len(ml.errors))
		}
	})

	t.Run("non existent directory", func(t *testing.T) {
		ml := newMockLogger()

		NewLocalLoader("edgecases/non-existent-directory", ml)

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(ml.errors))
		}
	})
}

func createTempFolderWithUnreadableFile(unreadableFileName string) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
		return ""
	}

	// Create an unreadable file in the temporary directory
	unreadableFilePath := path.Join(tempDir, unreadableFileName)
	err = os.WriteFile(unreadableFilePath, []byte("This is a test file."), 0o000)
	if err != nil {
		os.RemoveAll(tempDir)
		log.Fatal(err)
	}

	// Return the path to the temporary directory
	return tempDir
}
