// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// pipeline encapsulates a transformation which a request will come through
// from preprocessors to the actual handler. (and postprocessor later)
type pipeline struct {
	Tag           string
	Preprocessors []Processor
	Handler
}

// Router to dispatch HTTP request to respective handler
type Router struct {
	commonRouter
	actions struct {
		sync.RWMutex
		m map[string]pipeline
	}
}

// NewRouter is factory for Router
func NewRouter() *Router {
	r := &Router{
		actions: struct {
			sync.RWMutex
			m map[string]pipeline
		}{
			m: map[string]pipeline{},
		},
	}
	r.commonRouter.payloadFunc = NewPayload
	r.commonRouter.matchHandlerFunc = r.matchHandler
	return r
}

// Map to register action to handle mapping
func (r *Router) Map(action, tag string, handler Handler, preprocessors ...Processor) {
	r.actions.Lock()
	defer r.actions.Unlock()
	if len(preprocessors) == 0 {
		preprocessors = handler.GetPreprocessors()
	}
	r.actions.m[action] = pipeline{
		Tag:           tag,
		Preprocessors: preprocessors,
		Handler:       handler,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.commonRouter.ServeHTTP(w, req)
}

func (r *Router) matchHandler(p *Payload) (routeConfig, error) {
	r.actions.RLock()
	defer r.actions.RUnlock()

	// matching using URL
	action := p.Meta["path"].(string)
	if strings.HasPrefix(action, "/") {
		action = action[1:]
	}

	action = strings.Replace(action, "/", ":", -1)
	var matchedPipeline *pipeline
	if len(action) > 0 { // prevent matching HomeHandler
		if pipeline, ok := r.actions.m[action]; ok {
			matchedPipeline = &pipeline
		}
	}

	if matchedPipeline == nil {
		if pipeline, ok := r.actions.m[p.RouteAction()]; ok {
			matchedPipeline = &pipeline
		}
	}

	if matchedPipeline == nil {
		return routeConfig{}, errors.New("route unmatched")
	}

	return routeConfig{
		Tag:           matchedPipeline.Tag,
		Preprocessors: matchedPipeline.Preprocessors,
		Handler:       matchedPipeline.Handler,
	}, nil
}

func NewPayload(req *http.Request) (p *Payload, err error) {
	reqBody := req.Body
	if reqBody == nil {
		reqBody = ioutil.NopCloser(bytes.NewReader(nil))
	}

	data := map[string]interface{}{}
	if jsonErr := json.NewDecoder(reqBody).Decode(&data); jsonErr != nil && jsonErr != io.EOF {
		err = jsonErr
		return
	}

	p = &Payload{
		Data: data,
		Meta: map[string]interface{}{},
	}
	p.SetContext(req.Context())

	if apiKey := req.Header.Get("X-Skygear-Api-Key"); apiKey != "" {
		p.Data["api_key"] = apiKey
	}
	if accessToken := req.Header.Get("X-Skygear-Access-Token"); accessToken != "" {
		p.Data["access_token"] = accessToken
	}

	p.Meta["path"] = req.URL.Path
	p.Meta["method"] = req.Method
	p.Meta["remote_addr"] = req.RemoteAddr
	if xff := req.Header.Get("x-forwarded-for"); xff != "" {
		p.Meta["x_forwarded_for"] = xff
	}
	if xri := req.Header.Get("x-real-ip"); xri != "" {
		p.Meta["x_real_ip"] = xri
	}
	if forwarded := req.Header.Get("forwarded"); forwarded != "" {
		p.Meta["forwarded"] = forwarded
	}

	return
}
