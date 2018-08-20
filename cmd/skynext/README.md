# Skygear next

This folder contains demo code for skygear next.

## Install dependencies

Please install dependencies at the top leve of `skygear-server`.

## Configuration

Configration can be set through environment variable.

- `IMPL`: `net/http`, `gorilla`, `httprouter`
- `GO_ENV`: `dev`, `production`

If env is `dev`, the following environment variable can also be set.

- `AUTH_PASSWORD_LENGTH`
- `RECORD_AUTO_MIGRATION`
- `ASSET_STORE_IMPL`
- `ASSET_STORE_SECRET`

This is list may not be the most updated, please see `/cmd/skynext-router/mode/config`.

To demo using old skygear handler, please configure the skynext server as if skygear-server.

## Run the program

Run `go run main.go`

## Project structure

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

## Gateway + Next structure

```
                                                                               HTTP Request
                                                                               /{module}/{path}
                                                                                     +
                                                                                     |
                                                                                     |
                                                                             +-------v-------+
                                                                             |               |
                                                                             |               |
                                                                             |    Gateway    |
                                                                             |               |
                                                                             |               |
                                                                             +-------+-------+
                                                                                     |
  Services                                                                           |
  (through Dependency Injection)                                                     |
+--------------------------------+                                         Header    |      Body
|                                |                                             +-----v-----+
|  Implementation     Interface  |                                 Request     |           |
|  +---------------+     +--+    |                                 Modifier    |           |
|  |               |     |  |    |                                  +--+       |           |
|  |               |     |  |    |  Function  +---------------+     |  |       |           |
|  | Configuration |     |  |    |  call      |               |     |  |Read   |           |
|  | Provider      |     |  <-----------------+               <-----+  <-------+           |
|  |               |     |  |    |            | Middleware(s) |     |  |       |           |
|  +---------------+     |  +----------------->               +----->  +------->           |
|                        |  |    |            |               |     |  |Write  |           |
|  +---------------+     |  |    |            +---------------+     |  |       |           |
|  |               |     |  |    |                                  +--+       |           |
|  |               |     |  |    |                                           +-v-----------v-+
|  |  XX Resource  |     |  |    |                                           |               | /record/{path}
|  |  Provider     |     |  |    |                                           |               +---------->
|  |               |     |  |    |                                           | Reverse Proxy | /cms/{path}
|  +---------------+     +--+    |                                           | (path rewrite)+---------->
|                                |                                           |               |
+--------------------------------+                                           +---------------+
                                                                                     |/auth/{path}
                                                                    +-------------------------------------+
                                                                                     |
                                                                                     |
                                                                    HTTP Request     |
                                                                                     |
                                                                                     |
                                                                    +-------------------------------------+
                                                                                     |/{path}
                                                                             +-------v-------+
                                                                             |               |
                                                                             |               |
                                                                             |  Auth Module  |
                                                                             |               |
                                                                             |               |
                                                                             +-------+-------+
                                                                                     |
                                                                                     |
                                                                                     |
                                                                                     |
                                                                             +-------v-------+
                                                                             |               | /signup
                                                                             |               +---------->
                                                                             |  App Router   | /logout
                                                                             |               +---------->
                                                                             |               |
                                                                             +---------------+
                                                                                     |/login
  Services                                                                           |
  (through Dependency Injection)                                                     |
+--------------------------------+                                         Header    |      Body
|                                |                                             +-----+-----+
|  Implementation     Interface  |                                 Request     |           |
|  +---------------+     +--+    |                                 Modifier    |           |
|  |               |     |  |    |                                  +--+       |           |
|  |               |     |  |    |  Function  +---------------+     |  |       |           |
|  |   DB Conn     |     |  |    |  call      |               |     |  |Read   |           |
|  |   Provider    |     |  <-----------------+               <-----+  <-------+           |
|  |               |     |  |    |            | Middleware(s) |     |  |       |           |
|  +---------------+     |  +----------------->               +----->  +------->           |
|                        |  |    |            |               |     |  |Write  |           |
|                        |  |    |            +---------------+     |  |       |           |
|                        |  |    |                                  |  |       |           |
|                        |  |    |                                  |  |       |           |
|                        |  |    |                            +-----+  <-------+           |
|                        |  |    |                            |     |  |Read               |
|                        |  |    |                            |     |  |                   |
|                        |  |    |                            |     |  |                   |
|                        |  |    |                            |     +--+                   |
|  +---------------+     |  |    |                            |                            |
|  |               |     |  |    |                       +----v----------------------------v------+
|  |               |     |  |    |                       |                                        |
|  |  XX Resource  |     |  <----------------------------+                                        |
|  |  Provider     |     |  |    |     Function call     |              Login Handler             |
|  |               |     |  +---------------------------->                                        |
|  +---------------+     +--+    |                       |                                        |
|                                |                       +--------------------+-------------------+
+--------------------------------+                                            |
                                                                              |
                                                                              |
                                                                              v
                                                                        HTTP Response

```

#### Notes

##### Request modifier (maybe not a good enough name)

It provides an interface for the middleware and handler to access and modify the http request and response.

For example,

- middleware in gateway can inject request id in request and response header.
- middleware can inject configuration to request header, and other middleware and handler later can access it from the request header. See [this](/cmd/skynext-router/model/config.go).

It acts like the interface to pass through states across the request pipeline. Since gateway and the module are two separate server, the request modifier can only modify the http request (instead of the golang request context).

## TODO

- ~~Demo using different golang web frameworks or router libraries~~
- Restrict the use of request context in middleware
- Generate model getter and setter functions
- Generate `http.handler` from a custom handler with life cycle (e.g. decode, validation, bussinuess logic, err handling, encode)
- ~~Demo using old skygear handler in new structure with minimal effort~~
- Create a package for shared library
