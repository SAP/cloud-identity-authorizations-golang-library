package dcn

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestLocalLoader(t *testing.T) {
	t.Run("broken data.json", func(t *testing.T) {
		errReceived := make(chan bool)
		errors := []error{}

		loader := NewLocalLoader("edgecases/broken-data-json")
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

		loader := NewLocalLoader("edgecases/broken-dcn")
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
		loader := NewLocalLoader(tmp)
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
		loader := NewLocalLoader(tmp)
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

		loader := NewLocalLoader("edgecases/non-existent-directory")
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
