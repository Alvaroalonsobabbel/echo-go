package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/Alvaroalonsobabbel/echo-go/internal/storage"
	"github.com/Alvaroalonsobabbel/echo-go/internal/types"
	"github.com/Alvaroalonsobabbel/echo-go/internal/validator"
)

const port = ":4567"

type Storer interface {
	Read() *types.EndpointsWrapper
	Create(types.Endpoint)
	Update(int, types.Endpoint) bool
	Delete(int) bool
	Find(string, string) (*types.Endpoint, bool)
}

var data Storer

func main() {
	log.Printf("Server started on port %s!", port)
	mux := SetupServer()
	startServer(mux)
}

func SetupServer() *http.ServeMux {
	data = storage.NewMemoryStorage() // In memory storage using an interface. Can be replaced by any other storage.
	mux := http.NewServeMux()
	mux.HandleFunc("/endpoints/", endpointsProxyHandler)
	mux.HandleFunc("/endpoints", endpointsHandler)
	mux.HandleFunc("/", handleAll)
	return mux
}

func startServer(mux *http.ServeMux) {
	WrappedMux := SetCommonHeaders(mux)
	log.Fatal(http.ListenAndServe(port, WrappedMux))
}

func SetCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		next.ServeHTTP(w, r)
	})
}

func endpointsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data.Read())
	case http.MethodPost:
		createEndpoint(w, r)
	default:
		handleAll(w, r)
	}
}

func endpointsProxyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		updateEndpoint(w, r)
	case http.MethodDelete:
		deleteEndpoint(w, r)
	default:
		handleAll(w, r)
	}
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	endpoint, ok := data.Find(r.Method, r.RequestURI)
	if ok {
		for k, v := range endpoint.Attributes.Response.Headers {
			w.Header().Add(k, v)
		}
		var decodedBody string
		if err := json.Unmarshal([]byte(endpoint.Attributes.Response.Body), &decodedBody); err != nil {
			log.Println(err)
		}
		w.WriteHeader(endpoint.Attributes.Response.Code)
		w.Write([]byte(decodedBody))
	} else {
		detail := fmt.Sprintf("Requested page `%s`, `%s` does not exist", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newError("not_found", detail))
	}
}

func createEndpoint(w http.ResponseWriter, r *http.Request) {
	wrapper, err := checkBody(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(newError("bad_request", err.Error()))
		return
	}
	defer r.Body.Close()

	id := generateID()
	wrapper.Data.ID = id
	data.Create(wrapper.Data)

	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	w.Header().Set("Location", scheme+r.Host+wrapper.Data.Attributes.Path)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wrapper)
}

func updateEndpoint(w http.ResponseWriter, r *http.Request) {
	wrapper, err := checkBody(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(newError("bad_request", err.Error()))
		return
	}
	defer r.Body.Close()

	// Ignoring this error since failing to convert the ID to into will send the request to handleAll
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/endpoints/"))

	if ok := data.Update(id, wrapper.Data); ok {
		wrapper.Data.ID = id
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(wrapper)
	} else {
		handleAll(w, r)
	}
}

func deleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/endpoints/"))

	if ok := data.Delete(id); ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		handleAll(w, r)
	}
}

func checkBody(requestBody io.Reader) (types.SingleEndpointWrapper, error) {
	var wrapper types.SingleEndpointWrapper
	body, err := io.ReadAll(requestBody)
	if err != nil {
		return types.SingleEndpointWrapper{}, err
	}

	err = validator.Validate(string(body))
	if err != nil {
		return types.SingleEndpointWrapper{}, err
	}

	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return types.SingleEndpointWrapper{}, err
	}

	return wrapper, nil
}

func newError(code string, detail string) types.ErrorResponse {
	return types.ErrorResponse{
		Errors: []types.ErrorDetail{
			{Code: code, Detail: detail},
		},
	}
}

func generateID() int {
	var id int
	var i int

IDGenerator:
	for i = 0; i < 9999; i++ {
		id = rand.Intn(10000)
		if id == 0 {
			continue
		} else {
			for _, d := range data.Read().Data {
				if d.ID == id {
					continue IDGenerator
				}
			}
		}
		return id
	}
	log.Fatalln("Not enough IDs")
	return 0
}
