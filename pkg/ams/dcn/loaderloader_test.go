package dcn

import (
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestLocalLoader(t *testing.T) {
	t.Run("on testfolder", func(t *testing.T) {
		errors := []error{}
		loader := NewLocalLoader("testfolder", nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
		})

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
		errReceived := make(chan bool)
		errors := []error{}

		loader := NewLocalLoader("edgecases/broken-data-json", nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
	})

	t.Run("broken DCN file", func(t *testing.T) {
		errReceived := make(chan bool)
		errors := []error{}

		loader := NewLocalLoader("edgecases/broken-dcn", nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
		}
	})

	t.Run("unreadable data.json", func(t *testing.T) {
		errReceived := make(chan bool)
		errors := []error{}
		tmp := createTempFolderWithUnreadableFile("data.json")
		defer os.RemoveAll(tmp) // Clean up
		loader := NewLocalLoader(tmp, nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Errorf("expected 1 request, got %d", len(errors))
		}
	})

	t.Run("unreadable DCN", func(t *testing.T) {
		errReceived := make(chan bool)
		errors := []error{}

		tmp := createTempFolderWithUnreadableFile("x.dcn")
		defer os.RemoveAll(tmp) // Clean up
		loader := NewLocalLoader(tmp, nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Errorf("expected 1 request, got %d", len(errors))
		}
	})

	t.Run("non existent directory", func(t *testing.T) {
		errReceived := make(chan bool)
		errors := []error{}

		loader := NewLocalLoader("edgecases/non-existent-directory", nil)
		loader.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			errReceived <- true
		})

		<-errReceived
		if len(errors) != 1 {
			t.Fatalf("expected 1 request, got %d", len(errors))
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
