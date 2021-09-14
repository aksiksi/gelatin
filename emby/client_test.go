package emby

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

func setUp(t *testing.T) (*EmbyApiClient, *httptest.Server, *mockEmbyServer) {
	t.Helper()

	s := &mockEmbyServer{}
	srv := httptest.NewServer(s)
	client := NewEmbyApiClient(srv.URL)

	return client, srv, s
}

func TestEmbySystemEndpoints(t *testing.T) {
	client, srv, s := setUp(t)
	defer srv.Close()

	apiKey := NewApiKey("test123")

	s.status = http.StatusOK

	t.Run("Ping", func(t *testing.T) {
		err := client.Ping()
		if err != nil {
			t.Errorf("failed to call ping endpoint")
		}
	})

	t.Run("Info_public", func(t *testing.T) {
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
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("Version", func(t *testing.T) {
		wantResp := []byte(`{
			"ServerName": "abc",
			"Version": "4.6.4.0",
			"Id": "ec68c767780f485d9fd4b3d58594f5ff"
		}`)
		s.resp = wantResp

		want := &gelatin.GelatinSystemInfo{}
		json.Unmarshal(wantResp, want)

		version, err := client.Version()
		if err != nil {
			t.Errorf("failed to call GetVersion")
		}

		if version != want.Version {
			t.Errorf("version mismatch: want %q != got %q", want, version)
		}
	})

	t.Run("GetLogFile", func(t *testing.T) {
		wantResp := []byte("this is a log file")
		s.resp = wantResp

		logReader, err := client.GetLogFile(apiKey, "test")
		if err != nil {
			t.Errorf("failed to call SystemLogs")
		}
		logData, _ := io.ReadAll(logReader)

		want, got := string(wantResp), string(logData)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("-want,+got:%s", diff)
		}
	})

	t.Run("GetLogs", func(t *testing.T) {
		data := readTestFile(t)
		s.resp = data

		want := &EmbySystemLogsQueryResponse{}
		json.Unmarshal(data, want)

		got, err := client.GetLogs(apiKey)
		if err != nil {
			t.Errorf("failed to call SystemLogsQuery")
		}

		if diff := cmp.Diff(want.Items, got); diff != "" {
			t.Errorf("-want,+got:%s", diff)
		}
	})
}

func TestEmbyUserEndpoints(t *testing.T) {
	client, srv, s := setUp(t)
	defer srv.Close()

	apiKey := NewApiKey("test123")

	s.status = http.StatusOK

	t.Run("GetUsers_public", func(t *testing.T) {
		wantResp := []byte(`[
			{
				"Name": "test",
				"Id": "100000x00000",
				"HasConfiguredPassword": true
			}
		]`)

		s.resp = wantResp

		var want []gelatin.GelatinUser
		json.Unmarshal(wantResp, &want)

		got, err := client.GetUsers(nil, true)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("GetUsers", func(t *testing.T) {
		wantResp := []byte(`{
			"Items": [
				{
					"Name": "test",
					"Id": "100000x00000",
					"HasConfiguredPassword": true
				}
			],
			"TotalRecordCount": 1
		}`)

		s.resp = wantResp

		var want *EmbyUserQueryResponse
		json.Unmarshal(wantResp, &want)

		got, err := client.GetUsers(apiKey, false)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}

		if diff := cmp.Diff(want.Items, got); diff != "" {
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		wantResp := []byte(`{
			"Name": "test",
			"Id": "100000x00000",
			"HasConfiguredPassword": true
		}`)

		s.resp = wantResp

		var want *gelatin.GelatinUser
		json.Unmarshal(wantResp, &want)

		got, err := client.GetUser(apiKey, want.Id)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("UserUpdate", func(t *testing.T) {
		user := &gelatin.GelatinUser{Id: "abcd123"}
		err := client.UpdateUser(apiKey, user.Id, user)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}
	})

	t.Run("UserNew", func(t *testing.T) {
		wantResp := []byte(`{
			"Name": "test",
			"Id": "100000x00000",
			"HasConfiguredPassword": true
		}`)

		s.resp = wantResp

		var want *gelatin.GelatinUser
		json.Unmarshal(wantResp, &want)

		got, err := client.CreateUser(apiKey, want.Name)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("UserDelete", func(t *testing.T) {
		err := client.DeleteUser(apiKey, "test123")
		if err != nil {
			t.Errorf("failed to call endpoint")
		}
	})

	t.Run("UserPassword", func(t *testing.T) {
		err := client.UpdatePassword(apiKey, "1000x1000", "", "test123", true)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}
	})

	t.Run("UserAuth", func(t *testing.T) {
		wantToken := "12345"
		wantUserAuthResp := &EmbyUserAuthResponse{
			AccessToken: wantToken,
		}
		wantResp, _ := json.Marshal(wantUserAuthResp)

		s.resp = wantResp

		key, err := client.Authenticate("abcd", "test123")
		if err != nil {
			t.Errorf("failed to call endpoint")
		}

		token := key.ToString()

		if diff := cmp.Diff(wantToken, token); diff != "" {
			t.Errorf("-want,+got: %s", diff)
		}
	})

	t.Run("UserPolicy", func(t *testing.T) {
		policy := &gelatin.GelatinUserPolicy{}
		err := client.UpdatePolicy(apiKey, "abcd", policy)
		if err != nil {
			t.Errorf("failed to call endpoint")
		}
	})
}
