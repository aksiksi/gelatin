package emby

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aksiksi/gelatin/api"
	"github.com/google/go-cmp/cmp"
)

type mockEmbyServer struct {
	resp   []byte
	status int
}

func (s *mockEmbyServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if s.status != http.StatusOK {
		http.Error(resp, http.StatusText(s.status), s.status)
		return
	}

	resp.WriteHeader(s.status)
	resp.Header().Add("Content-Type", "application/json")
	resp.Write(s.resp)
}

func setUp(t *testing.T) (*EmbyApiClient, *httptest.Server, *mockEmbyServer) {
	t.Helper()

	s := &mockEmbyServer{}
	srv := httptest.NewServer(s)
	client := NewEmbyApiClient(srv.URL, srv.Client())

	return client, srv, s
}

func TestEmbySystemEndpoints(t *testing.T) {
	client, srv, s := setUp(t)
	defer srv.Close()

	apiKey := api.NewApiKey("test123")

	t.Run("SystemPing", func(t *testing.T) {
		s.status = http.StatusOK
		err := client.SystemPing()
		if err != nil {
			t.Errorf("failed to call ping endpoint")
		}
	})

	t.Run("SystemInfoPublic", func(t *testing.T) {
		s.status = http.StatusOK
		wantResp := &EmbySystemInfoPublicResponse{
			LocalAddress: "127.0.0.1",
			Version:      "4.4.4",
		}
		bytes, _ := json.Marshal(wantResp)
		s.resp = bytes

		resp, err := client.SystemInfoPublic()
		if err != nil {
			t.Errorf("failed to call SystemInfoPublic endpoint")
		}

		if diff := cmp.Diff(wantResp, resp); diff != "" {
			t.Errorf("+want,-got: %s", diff)
		}
	})

	t.Run("GetVersion", func(t *testing.T) {
		s.status = http.StatusOK
		wantResp := &EmbySystemInfoPublicResponse{
			LocalAddress: "127.0.0.1",
			Version:      "4.4.4",
		}
		bytes, _ := json.Marshal(wantResp)
		s.resp = bytes

		version, err := client.GetVersion()
		if err != nil {
			t.Errorf("failed to call SystemInfoPublic endpoint")
		}

		if version != wantResp.Version {
			t.Errorf("version mismatch")
		}
	})

	t.Run("SystemLogs", func(t *testing.T) {
		s.status = http.StatusOK
		wantResp := []byte("this is a log file")
		s.resp = wantResp

		logReader, err := client.SystemLogs(apiKey, "test")
		if err != nil {
			t.Errorf("failed to call SystemInfoPublic endpoint")
		}
		logData, _ := io.ReadAll(logReader)

		want, got := string(wantResp), string(logData)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("+want,-got:%s", diff)
		}
	})
}
