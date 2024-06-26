package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/Alvaroalonsobabbel/echo-go/internal/types"
)

var baseURL string
var exmaple1id int
var exmaple2id int
var client = &http.Client{}

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

		if resp.Header["Content-Type"][0] != "application/vnd.api+json" {
			t.Errorf("Expected Content-Type to be 'application/vnd.api+json', got: %v", resp.Header["Content-Type"][0])
		}

		body, _ := io.ReadAll(resp.Body)
		var response types.EndpointsWrapper
		json.Unmarshal(body, &response)

		if len(response.Data) != 0 {
			t.Errorf("Expected Data array to be empty, got: %v", response.Data)
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
		var errorResponse types.ErrorResponse
		expectedCode := "not_found"
		expectedDetail := "Requested page `GET`, `/endpointss` does not exist"
		json.Unmarshal(body, &errorResponse)

		if errorResponse.Errors[0].Code != expectedCode {
			t.Errorf("Expected error code 'not_found', got: %v", errorResponse.Errors[0].Code)
		}
		if errorResponse.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected detail %v, got: %v", expectedDetail, errorResponse.Errors[0].Detail)
		}
	})

	t.Run("POST /endpoint example 1", func(t *testing.T) {
		endpointFile, err := os.ReadFile("main_test_examples/endpoint_example.json")
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

		body, _ := io.ReadAll(resp.Body)
		var response types.SingleEndpointWrapper
		expectedVerb := "GET"
		expectedPath := "/greeting"
		expectedCode := 200
		expectedBody := `"\"{ \"message\": \"Hello, world\" }\""`

		if resp.Header["Location"][0] != baseURL+expectedPath {
			t.Errorf("Expected header Location: %v%v, got: %v", baseURL, expectedPath, resp.Header["Location"][0])
		}

		json.Unmarshal(body, &response)
		exmaple1id = response.Data.ID

		if response.Data.ID == 0 {
			t.Errorf("Expected generated ID to not be 0, got: %v", response.Data.ID)
		}
		if response.Data.Attributes.Path != expectedPath {
			t.Errorf("Expected path %s, got: %v", expectedPath, response.Data.Attributes.Path)
		}
		if response.Data.Attributes.Response.Code != expectedCode {
			t.Errorf("Expected code %s, got: %v", expectedVerb, response.Data.Attributes.Response.Code)
		}
		if string(response.Data.Attributes.Response.Body) != (expectedBody) {
			t.Errorf("Expected body %s, got: %v", expectedBody, string(response.Data.Attributes.Response.Body))
		}
		if len(response.Data.Attributes.Response.Headers) != 0 {
			t.Errorf("Expected headers to be empty, but got: %v", response.Data.Attributes.Response.Headers)
		}
	})

	t.Run("POST /endpoint example 2", func(t *testing.T) {
		endpointFile, err := os.ReadFile("main_test_examples/endpoint_example2.json")
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

		body, _ := io.ReadAll(resp.Body)
		var response types.SingleEndpointWrapper
		expectedVerb := "POST"
		expectedPath := "/test"
		expectedCode := 404

		json.Unmarshal(body, &response)
		exmaple2id = response.Data.ID

		if response.Data.ID == 0 {
			t.Errorf("Expected generated ID to not be 0, got: %v", response.Data.ID)
		}
		if response.Data.Attributes.Path != expectedPath {
			t.Errorf("Expected path %s, got: %v", expectedPath, response.Data.Attributes.Path)
		}
		if response.Data.Attributes.Response.Code != expectedCode {
			t.Errorf("Expected code %s, got: %v", expectedVerb, response.Data.Attributes.Response.Code)
		}
		if string(response.Data.Attributes.Response.Body) != "" {
			t.Errorf("Expected an empty body, got: %v", string(response.Data.Attributes.Response.Body))
		}
		if response.Data.Attributes.Response.Headers["x-some-header"] != "header" {
			t.Errorf("Expected headers to be 'x-some-header': 'header'`, but got: %v", response.Data.Attributes.Response.Headers)
		}
	})

	t.Run("GET /endpoints to check endopints have been created", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/endpoints")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var response types.EndpointsWrapper
		json.Unmarshal(body, &response)

		if len(response.Data) != 2 {
			t.Errorf("Expected Data array to have 2 items, got: %v", response.Data)
		}
	})

	t.Run("calling example 1", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/greeting")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var response struct {
			Message string `json:"message"`
		}
		json.Unmarshal(body, &response)

		if response.Message != "Hello, world" {
			t.Errorf("Expected body to have a message key with 'Hello, world', got: %v", response.Message)
		}
	})

	t.Run("calling example 2", func(t *testing.T) {
		var emptyBody = "{[]}"
		resp, err := http.Post(baseURL+"/test", "application/json", bytes.NewBuffer([]byte(emptyBody)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status Not Found; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		if len(body) != 0 {
			t.Errorf("Expected body to be empty, got: %v", string(body))
		}
		const headerKey = "x-some-header"
		const headerValue = "header"

		if resp.Header[http.CanonicalHeaderKey(headerKey)][0] != headerValue {
			t.Errorf("Expected %v: %v, got: %v", headerKey, headerValue, resp.Header[http.CanonicalHeaderKey(headerKey)])
		}
	})

	t.Run("PATCH /endpoint/{id} example 2", func(t *testing.T) {
		endpointFile, err := os.ReadFile("main_test_examples/endpoint_example2_rename.json")
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		id := strconv.Itoa(exmaple2id)
		req, err := http.NewRequest(http.MethodPatch, baseURL+"/endpoints/"+id, bytes.NewBuffer(endpointFile))
		if err != nil {
			t.Fatalf("Failed to prepare the request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var response types.SingleEndpointWrapper
		expectedVerb := "GET"
		expectedPath := "/test"
		expectedCode := 404

		json.Unmarshal(body, &response)

		if response.Data.ID == 0 {
			t.Errorf("Expected generated ID to not be 0, got: %v", response.Data.ID)
		}
		if response.Data.Attributes.Path != expectedPath {
			t.Errorf("Expected path %s, got: %v", expectedPath, response.Data.Attributes.Path)
		}
		if response.Data.Attributes.Response.Code != expectedCode {
			t.Errorf("Expected code %s, got: %v", expectedVerb, response.Data.Attributes.Response.Code)
		}
		if string(response.Data.Attributes.Response.Body) != "" {
			t.Errorf("Expected an empty body, got: %v", string(response.Data.Attributes.Response.Body))
		}
		if response.Data.Attributes.Response.Headers["x-some-header"] != "header" {
			t.Errorf("Expected headers to be 'x-some-header': 'header'`, but got: %v", response.Data.Attributes.Response.Headers)
		}
	})

	t.Run("calling example 2 after edit", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/test")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status Not Found; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		if len(body) != 0 {
			t.Errorf("Expected body to be empty, got: %v", string(body))
		}
		const headerKey = "x-some-header"
		const headerValue = "header"

		if resp.Header[http.CanonicalHeaderKey(headerKey)][0] != headerValue {
			t.Errorf("Expected %v: %v, got: %v", headerKey, headerValue, resp.Header[http.CanonicalHeaderKey(headerKey)])
		}
	})

	t.Run("DELETE /endpoint/{id} example 2", func(t *testing.T) {
		id := strconv.Itoa(exmaple2id)
		fmt.Println(id)
		req, err := http.NewRequest(http.MethodDelete, baseURL+"/endpoints/"+id, bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatalf("Failed to prepare the request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status No Content; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)

		if body == nil {
			t.Errorf("Expected body to be empty, got: %v", body)
		}
	})

	t.Run("GET /endpoints to verify item has been deleted", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/endpoints")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var response types.EndpointsWrapper
		json.Unmarshal(body, &response)

		if len(response.Data) != 1 {
			t.Errorf("Expected Data array to have 2 items, got: %v", response.Data)
		}
	})

	t.Run("POST /endpoint with incorrect verb", func(t *testing.T) {
		endpointFile, err := os.ReadFile("main_test_examples/endpoint_example_error.json")
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		resp, err := http.Post(
			baseURL+"/endpoints", "application/json", bytes.NewBuffer(endpointFile))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status Bad Request; got %v", resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		var errorResponse types.ErrorResponse
		expectedCode := "bad_request"
		json.Unmarshal(body, &errorResponse)

		if errorResponse.Errors[0].Code != expectedCode {
			t.Errorf("Expected error code 'bad_request', got: %v", errorResponse.Errors[0].Code)
		}
		if errorResponse.Errors[0].Detail == "" {
			t.Errorf("Expected detail not to be empty")
		}
	})
}
