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

## Dependency injection with multi-tenant

**TL;DR** See [Multi-tenant version (with request context and separate struct for dependency)](#multi-tenant-version-with-request-context-and-separate-struct-for-dependency).

#### (The good old) Single-tenant version

```golang
type AHandler struct {
  TokenStore TokenStore `inject:"TokenStore"`
}

type TokenStore interface {
  Get(tokenString string) Token
}

// 1. In main.go
func main() {
  configuration := GetConfigurationFromEnv()

  serveMux := http.NewServeMux()

  // 1.1 prepare handler
  handler := AHandler{}

  // 1.2 prepare dependency (use a DI library in real world)
  tokenStore := NewTokenStore(configuration)

  // 1.3 inject dependency (use a DI library in real world)
  handler.TokenStore = tokenStore

  // 1.4 register handler to a route
  // 1.5 start http server
}

// 2. in handler/a.go, using TokenStore in AHandler
func (h AHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  tokenStore := h.TokenStore
  token := tokenStore.Get(GetTokenString(r))

  // .....
}
```

#### Mutli-tenant version

```golang
type AHandler struct {
  TokenStoreProvider TokenStoreProvider `inject:"TokenStoreProvider"`
}

type TokenStoreProvider interface {
  GetStore(config Configuration) TokenStore
}

type TokenStore interface {
  Get(tokenString string) Token
}

// 1. In main.go
func main() {
  configuration := GetConfigurationFromEnv()

  serveMux := http.NewServeMux()

  // 1.1 prepare handler
  handler := AHandler{}

  // 1.2 prepare dependency (use a DI library in real world)
  tokenStoreProvider := NewTokenStoreProvider()

  // 1.3 inject dependency (use a DI library in real world)
  handler.TokenStore = tokenStore

  // 1.4 register handler to a route
  // 1.5 start http server
}

// 2. in handler/a.go, using TokenStore in AHandler
func (m AHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // 2.1 get configuration associated with current request
  configuration := GetConfiguration(r)

  // 2.2 create a token store with the configuration
  tokenStore := m.TokenStoreProvider.GetStore(configuration)

  token := tokenStore.Get(GetTokenString(r))

  // ...
}
```

Pros:

- Inject once when server start, in the `main.go`, its very similar to single-tenant version

Problems:

- An extra layer of abstraction (i.e. TokenStoreProvider) for creating the tenancy-unaware depedency (i.e. TokenStore) is introduced in each middleware and handler,

#### Multi-tenant version (with request context)

```golang
type AHandler struct {}

type TokenStoreProvider interface {
  GetStore(config Configuration) TokenStore
}

type TokenStore interface {
  Get(tokenString string) Token
}

// 1. In main.go
func main() {
  serveMux := http.NewServeMux()

  // 1.1 prepare handler and middleware
  handler := AHandler{}
  configMiddleware := ConfigurationMiddleware{}
  diMiddleware := DIMiddleware{
    TokenStoreProvider: TokenStoreProvider{},
  }
  otherMiddleware := XXXMiddleware{}

  // 1.2 wrap the handler with middleware
  // note: the configMiddleware and diMiddleware must come first
  wrappedHandler := configMiddleware.Handle(diMiddleware.Handle(otherMiddleware.Handle(handler)))

  // 1.3 register the handler to a route
  // 1.4 start http server
}

type DIMiddleware struct {
  TokenStoreProvider TokenStoreProvider
}

// 2. inject dependency in the middleware
func (m DIMiddleware) Handle(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // 2.1 get configuration from the request
    configuration := GetConfiguration(r)

    // 2.2 create a dependency instance with the configuration
    tokenStore := m.TokenStoreProvider.GetStore(configuration)

    // 2.3 set the depedenecy instance in the request
    SetTokenStore(r, tokenStore)

    // inject other dependencies here
    // ...

    next.ServeHTTP(w, r)
  }
}

// 3. in handler/a.go, using TokenStore in AHandler
func (m AHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  tokenStore := GetTokenStore(r)
  token := tokenStore.Get(GetTokenString(r))

  // ...
}
```

Pros:

- Middlewares and handlers do not have knowledge of tenants

Problems:

- Type-less (can be solved by getter and setter)
- Implicit dependency, dependencies cannot be declared per middleware and handler

#### Multi-tenant version (with request context and separate struct for dependency)

The following may solve the second problem metioned above.

```golang
type AHandler struct {}

// define the dependencies required by A
type AHandlerDependency struct {
  TokenStore service.token `inject:"TokenStore"`
}

// 1. In main.go (same as previous version)
func main() {
  serveMux := http.NewServeMux()

  // 1.1 prepare handler and middleware
  handler := AHandler{}
  configMiddleware := ConfigurationMiddleware{}
  diMiddleware := DIMiddleware{
    TokenStoreProvider: TokenStoreProvider{},
  }
  otherMiddleware := XXXMiddleware{}

  // 1.2 wrap the handler with middleware
  // note: the configMiddleware and diMiddleware must come first
  wrappedHandler := configMiddleware.Handle(diMiddleware.Handle(otherMiddleware.Handle(handler)))

  // 1.3 register the handler to a route
  // 1.4 start http server
}

type DIMiddleware struct {
  TokenStoreProvider TokenStoreProvider
}

// 2. inject dependency in the middleware (same as previous version)
func (m DIMiddleware) Handle(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // 2.1 get configuration from the request
    configuration := GetConfiguration(r)

    // 2.2 create a dependency instance with the configuration
    tokenStore := m.TokenStoreProvider.GetStore(configuration)

    // 2.3 set the depedenecy instance in the request
    SetTokenStore(r, tokenStore)

    // inject other dependencies here
    // ...

    next.ServeHTTP(w, r)
  }
}

// 3. in handler/a.go, inject dependency required by AHandler and use it in AHandler
func (m AHandler) Handle(w http.ResponseWriter, r *http.Request) {
  // 3.1 create an empty dependency struct
  dep := AHandlerDependency{}

  // 3.2 inject dependecies (associate with the request) to the struct
  InjectDependency(r, &dep)

  token := dep.tokenStore.Get(GetTokenString(r))

  // ...
}

// note: this function can be used in any middlewares and handlers
func InjectDependency(h *http.Request, dep interface{}) { /* ... */ }
```

## TODO

- ~~Demo using different golang web frameworks or router libraries~~
- Restrict the use of request context in middleware
- Generate model getter and setter functions
- Generate `http.handler` from a custom handler with life cycle (e.g. decode, validation, bussinuess logic, err handling, encode)
- ~~Demo using old skygear handler in new structure with minimal effort~~
- Create a package for shared library
