package types

import "encoding/json"

type EndpointsWrapper struct {
	Data []Endpoint `json:"data"`
}

type SingleEndpointWrapper struct {
	Data Endpoint `json:"data"`
}

type Endpoint struct {
	Type       string `json:"type,omitempty"`
	ID         int    `json:"id,omitempty"`
	Attributes struct {
		Verb     string `json:"verb,omitempty"`
		Path     string `json:"path,omitempty"`
		Response struct {
			Code    int               `json:"code,omitempty"`
			Headers map[string]string `json:"headers,omitempty"`
			Body    json.RawMessage   `json:"body,omitempty"`
		} `json:"response,omitempty"`
	} `json:"attributes,omitempty"`
}

type ErrorResponse struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}
