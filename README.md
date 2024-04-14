# echo - Go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Alvaroalonsobabbel/echo-go) ![Test](https://github.com/Alvaroalonsobabbel/echo-go/actions/workflows/test.yml/badge.svg)

Echo code challenge based on [these requirements](echo.md)

## Techincal information

This application was built using Go and the endpoints use [JSON:API v1.0](https://jsonapi.org/) as a format.

## Run locally

If you're using macOS you can just download the app binary [here](https://github.com/Alvaroalonsobabbel/echo-go/releases/latest/download/echo-go), give execution access to the file with `chmod +x echo-go` and start the server with `./echo-go`. You might also have to allow the app to run in the System Settings.

Otherwise you'll have to:

1. [Install Go](https://go.dev/doc/install)
2. Download dependencies using `go mod download`
3. Optionally you can run the tests using `go test -v ./...`
4. Run the server using `go run cmd/main.go`

Use cURL or Postman to send http requests to the server at `http://127.0.0.1:4567`
Server works using the exact API documentation specificed in the [requirements](echo.md#examples)

## Quick cURL commands to test the server

View endpoints:

```bash
curl -L -X GET 'http://127.0.0.1:4567/endpoints' 
```

Submit an endpoint:

```bash
curl -L -X POST 'http://127.0.0.1:4567/endpoints' \
-H 'Content-Type: application/vnd.api+json' \
-d '{
    "data": {
        "type": "endpoints",
        "attributes": {
            "verb": "GET",
            "path": "/revert_entropy",
            "response": {
              "code": 200,
              "headers": {},
              "body": "\"{ \"message\": \"INSUFFICIENT DATA FOR MEANINGFUL ANSWER\" }\""
            }
        }
    }
}'
```

Now you can run View endpoints again to check the enpoint is there.

Use the endpoint:

```bash
curl -L -X GET 'http://127.0.0.1:4567/revert_entropy' 
```
