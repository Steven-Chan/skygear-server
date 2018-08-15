### Skygear next

This folder contains demo code for skygear next.

#### Install dependencies

Please install dependencies at the top leve of `skygear-server`.

#### Configuration

Configration can be set through environment variable.

- `IMPL`: `net/http`, `gorilla`, `httprouter`
- `GO_ENV`: `dev`, `production`

If env is `dev`, the following environment variable can also be set.

- `AUTH_PASSWORD_LENGTH`
- `RECORD_AUTO_MIGRATION`
- `ASSET_STORE_IMPL`
- `ASSET_STORE_SECRET`

This is list may not be the most updated, please see `/cmd/skynext-router/mode/config`.

#### Run the program

Run `go run main.go`

#### Project structure

This codebase is divided into the following parts:

- `handler`
  - Implements `http.Handler`, contains businuess logic
- `middleware`
  - Implements `func(http.Handler) http.Handler`, contains common logic across multiple api endpoints
- `model`
  - Defines model
  - Provide functions to get and set the model in a request
- `service`
  - Contains various logic which are meant to be injected in `handler` and `middleware`

#### TODO

- ~~Demo using different golang web frameworks or router libraries~~
- Restrict the use of request context in middleware
- Generate model getter and setter functions
- Generate `http.handler` from a custom handler with life cycle (e.g. decode, validation, bussinuess logic, err handling, encode)
