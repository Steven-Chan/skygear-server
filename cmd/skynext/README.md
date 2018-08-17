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
                                                                                                        /{path}
                                                                                                              +
                                                                                                              |
                                                                                                      +-------v-------+
                                                                                                      |               |
                                                                                                      |               |
                                                                                                      |  Auth Module  |
                                                                                                      |               |
                                                                                                      |               |
                                                                                                      +-------+-------+
                                                                                                              |
                                                            HTTP Request                                      |
                                                            /{module}/{path}                                  |
                                                                  +                                           |
                                                                  |                                           |
                                                          +-------v-------+                           +-------v-------+
                                                          |               |                           |               |
                                                          |               |                           |               |
                                                          |    Gateway    |                           |    Router     |
                                                          |               |                           |               |
                                                          |               |                           |               |
                                                          +-------+-------+                           +-------+-------+
                                                                  |                                           |
+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|                                                                 |                                           |                                                                              |
|                      Services                         Request   |                 Skygear Next              |        Request                          Services       Implementation        |
|  Implementation      Interface                        Modifier  |                                           |        Modifier                         Interface      +---------------+     |
|  +---------------+     +--+     +---------------+      +--+     |                                           |         +--+   +---------------+          +--+         |               |     |
|  |               |     |  |     |               |      |  |     |                                           |         |  |   |               |          |  +-------->+               |     |
|  |               +----->  +----->               +--------------->                                           +---------------->               +---------->  |         |   DB Conn     |     |
|  |  XX Resource  |     |  |     |  Middleware   |      |  |     |                                           |         |  |   |  Middleware   |          |  <---------+   Provider    |     |
|  |  Provider     <-----+  <-----+               <---------------+                                           <----------------+               <----------+  |         |               |     |
|  |               |     |  |     |               |      |  |     |                                           |         |  |   |               |          |  |         +---------------+     |
|  +---------------+     +--+     +---------------+      +--+     |                                           |         |  |   +---------------+          |  |         +---------------+     |
|                                                                 |                                           |         |  |                              |  |         |               |     |
|                                                                 |                                           |         |  |                              |  +-------->+               |     |
|                                                                 |                                           |         |  |                              |  |         | Configuration |     |
+--------------------------------------------------------------------------------------------------------+    |         |  |                              |  <---------+ Provider      |     |
                                                                  |                                      |    |         |  |                              |  |         |               |     |
                                                          +-------v-------+                              |    |         |  |                              |  |         +---------------+     |
                                                          |               |                              |    |         |  |                              |  |         +---------------+     |
                                                          |               |                              |    |         |  |           +------------------>  |         |               |     |
                                                          | Revrese Proxy |                              |    |         |  |           |                  |  +-------->+               |     |
                                                          | (path rewrite)|                              |    +------------------+     | +----------------+  |         |  XX Resource  |     |
                                                          |               |                              |              |  |     |     | |                |  <---------+  Provider     |     |
                                                          +-------+-------+                              |              |  |     |     | |                +^-+         |               |     |
                                                                  |                                      |              +--+     |     | |                 |           +---------------+     |
                                                                  |                                      |                       |     | |                 |                                 |
                                            /auth/{path}          |/record/{path}          /cms/{path}   +-----------------------------------------------------------------------------------+
                                           +----------------------------------------------+                                      |     | |                 |
                          HTTP Request     |                      |                       |                                      |     | |                 |
                                           |                      |                       |                                      +---------------------+   |
                                  ------------------------------------------------------------------                             |     | |             |   |
                                           |                      |                       |                                      |     | |             |   |
                                           |/{path}               |/{path}                |/{path}                               |     | |     +----------------------------+ skygear-server handler
                                           |                      |                       |                                      |     | |     |       |   |                | to
                                           |                      |                       |                                 +----v-----+-v--+  |    +--v---v-----------+    | Handler
                                   +-------v-------+      +-------v-------+       +-------v-------+                         |               |  |    |     Preprocessor |    |
                                   |               |      |               |       |               |                         |  Handler      |  |    +--+---------------+    |
                                   |               |      |               |       |               |                         |               |  |       |                    |
                                   |  Auth Module  |      | Record Module |       |   CMS Module  |                         |  Businuess    |  |    +--v---------------+    |
                                   |               |      |               |       |               |                         |  Logic        |  |    |     Handler      |    |
                                   |               |      |               |       |               |                         +--------+------+  |    +--+---------------+    |
                                   +---------------+      +---------------+       +---------------+                                  |         |       |                    |
                                     localhost:3001         localhost:3002          localhost:3003                                   |         +----------------------------+
                                                                                                                                     v                 v
                                                                                                                               HTTP Response        HTTP Response

```

## TODO

- ~~Demo using different golang web frameworks or router libraries~~
- Restrict the use of request context in middleware
- Generate model getter and setter functions
- Generate `http.handler` from a custom handler with life cycle (e.g. decode, validation, bussinuess logic, err handling, encode)
- ~~Demo using old skygear handler in new structure with minimal effort~~
- Create a package for shared library
