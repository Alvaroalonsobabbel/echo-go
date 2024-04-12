package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Alvaroalonsobabbel/echo-go/internal/types"
)

var baseURL string
var dbCheck types.EndpointsWrapper

func TestMain(m *testing.M) {
	mux := SetupServer()
	wrappedMux := SetCommonHeaders(mux)
	server := httptest.NewServer(wrappedMux)
	defer server.Close()

	baseURL = server.URL

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestServer(t *testing.T) {
	t.Run("GET /endpoints", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/endpoints")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}
		expectedBody := `{"data":[]}`
		body, _ := io.ReadAll(resp.Body)
		body = bytes.TrimSpace(body)

		if string(body) != expectedBody {
			t.Errorf("Expected body %s; got %s", expectedBody, body)
		}
	})

	t.Run("GET /endpointss", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/endpointss")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status Not Found; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var errorResponse map[string][]types.ErrorResponse
		expectedCode := "not_found"
		expectedDetail := "Requested page `GET`, `/endpointss` does not exist"
		json.Unmarshal(body, &errorResponse)

		if errorResponse["errors"][0].Code != expectedCode {
			t.Errorf("Expected error code 'not_found', got: %v", errorResponse["errors"][0].Code)
		}
		if errorResponse["errors"][0].Detail != expectedDetail {
			t.Errorf("Expected error code 'not_found', got: %v", errorResponse["errors"][0].Detail)
		}
	})

	t.Run("POST /endpoint", func(t *testing.T) {
		endpointFile, err := os.ReadFile("test_examples/endpoint_example.json")
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		resp, err := http.Post(
			baseURL+"/endpoints", "application/json", bytes.NewBuffer(endpointFile))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}
	})
}
