package httpclient_test

import (
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func startTestServerContainer(scenario string) {
	stopTestServerContainer()
	repoRootAbs := absRepoPath()

	buildCmd := exec.Command( //nolint:gosec,noctx
		"docker",
		"build",
		"-f",
		path.Join(repoRootAbs, "http", "Dockerfile"),
		"-t",
		"ams-http-test",
		path.Join(repoRootAbs, "http"),
	)

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		panic(err)
	}

	dcnDir := path.Join(repoRootAbs, "pkg", "ams", "test", "scenarios", scenario)

	runCmd := exec.Command( //nolint:gosec,noctx
		"docker",
		"run",
		"-d",
		"--rm",
		"-p", "8099:8099",
		"-v", dcnDir+":/bundle",
		"-e", "AMS_DCN_ROOT=/bundle",
		"--name", "ams-http-test-container",
		"ams-http-test",
	)

	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		panic(err)
	}
	waitForTestServer()
}

func waitForTestServer() {
	interval := 10
	tries := 20

	for range tries {
		resp, err := http.Get("http://localhost:8099/v1/health") //nolint:noctx
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	panic("test server did not become ready in time")
}

func stopTestServerContainer() {
	exec.Command( //nolint:noctx
		"docker",
		"stop",
		"ams-http-test-container",
	).Run() //nolint:errcheck
}

func absRepoPath() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to determine current file path")
	}

	// pkg/ams/httpclient
	testDir := filepath.Dir(currentFile)

	// Repository root
	repoRoot := filepath.Join(testDir, "..", "..", "..")

	repoRootAbs, err := filepath.Abs(repoRoot)
	if err != nil {
		panic(err)
	}
	return repoRootAbs
}
