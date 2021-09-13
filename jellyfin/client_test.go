package jellyfin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	gelatin "github.com/aksiksi/gelatin/lib"
	"github.com/google/go-cmp/cmp"
)

type mockJellyfinServer struct {
	resp   []byte
	status int
}

func (s *mockJellyfinServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if s.status != http.StatusOK {
		http.Error(resp, http.StatusText(s.status), s.status)
		return
	}

	resp.WriteHeader(s.status)
	resp.Header().Add("Content-Type", "application/json")
	resp.Write(s.resp)
}

func readTestFile(t *testing.T) []byte {
	t.Helper()
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	path := fmt.Sprintf("testdata/%s.json", testName)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open test file: %s", path)
	}
	data, _ := io.ReadAll(f)
	return data
}

func setUp(t *testing.T) (*JellyfinApiClient, *httptest.Server, *mockJellyfinServer) {
	t.Helper()

	s := &mockJellyfinServer{}
	srv := httptest.NewServer(s)
	client := NewJellyfinApiClient(srv.URL)

	return client, srv, s
}

func TestJellyfinSystemEndpoints(t *testing.T) {
	client, srv, s := setUp(t)
	defer srv.Close()

	apiKey := NewApiKey("test123")

	s.status = http.StatusOK

	t.Run("SystemPing", func(t *testing.T) {
		err := client.Ping()
		if err != nil {
			t.Errorf("failed to call ping endpoint")
		}
	})

	t.Run("SystemInfoPublic", func(t *testing.T) {
		wantResp := []byte(`{
			"ServerName": "abc",
			"Version": "4.6.4.0",
			"Id": "ec68c767780f485d9fd4b3d58594f5ff"
		}`)

		s.resp = wantResp

		want := &gelatin.GelatinSystemInfo{}
		json.Unmarshal(wantResp, want)

		got, err := client.Info(nil, true)
		if err != nil {
			t.Errorf("failed to call SystemInfoPublic endpoint")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("+want,-got: %s", diff)
		}
	})

	t.Run("GetVersion", func(t *testing.T) {
		wantResp := []byte(`{
			"ServerName": "abc",
			"Version": "4.6.4.0",
			"Id": "ec68c767780f485d9fd4b3d58594f5ff"
		}`)
		s.resp = wantResp

		want := &JellyfinSystemInfoPublicResponse{}
		json.Unmarshal(wantResp, want)

		version, err := client.Version()
		if err != nil {
			t.Errorf("failed to call GetVersion")
		}

		if version != want.Version {
			t.Errorf("version mismatch: want %q != got %q", want, version)
		}
	})

	t.Run("SystemLogsName", func(t *testing.T) {
		wantResp := []byte("this is a log file")
		s.resp = wantResp

		logReader, err := client.GetLogFile(apiKey, "test")
		if err != nil {
			t.Errorf("failed to call SystemLogs")
		}
		logData, _ := io.ReadAll(logReader)

		want, got := string(wantResp), string(logData)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("+want,-got:%s", diff)
		}
	})

	t.Run("SystemLogs", func(t *testing.T) {
		data := readTestFile(t)

		s.resp = data

		var want []gelatin.GelatinSystemLog
		json.Unmarshal(data, &want)

		got, err := client.GetLogs(apiKey)
		if err != nil {
			t.Errorf("failed to call SystemLogsQuery")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("+want,-got:%s", diff)
		}
	})
}
