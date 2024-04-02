// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/log"
	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/naming"
	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/revoker"
	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/server"
)

var HandlerRegisterer = registerer{naming.PluginName, log.NoopLogger{}, nil, nil}

type registerer struct {
	name string

	logger  log.ILogger
	server  *server.Server
	revoker *revoker.Revoker
}

func (r *registerer) RegisterLogger(v interface{}) {
	l, ok := v.(log.ILogger)
	if !ok {
		return
	}
	r.logger = l
	r.logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer.String()))
}

func (r *registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(r.String(), r.registerHandlers)
}

func (r *registerer) String() string {
	return r.name
}

func (r *registerer) registerHandlers(_ context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	var (
		addr     = ":9000"
		key      = "Authorization"
		keyType  = "Header"
		response = map[string]interface{}{"status_code": http.StatusUnauthorized, "error": "Token is revoked"}
	)

	config, ok := extra[r.String()].(map[string]interface{})
	if ok {
		if cAddr, ok := config["addr"].(string); ok {
			addr = cAddr
		}
		if cKey, ok := config["key"].(string); ok {
			key = cKey
		}
		if cKeyType, ok := config["type"].(string); ok {
			keyType = cKeyType
		}
		if cResponse, ok := config["response"].(map[string]interface{}); ok {
			response = cResponse
		}
	}

	r.revoker = revoker.NewRevoker(key, keyType)
	r.server = server.NewServer(addr, r.revoker, r.logger)
	r.server.Run()

	return r.handler(response, handler), nil
}

func (r *registerer) handler(response map[string]interface{}, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := r.revoker.Middleware(req); err != nil {
			if errors.Is(err, revoker.ErrorTokenIsRevoked) {
				r.write(http.StatusUnauthorized, response, w)
			} else {
				r.write(http.StatusInternalServerError, []byte(err.Error()), w)
			}
			return
		}

		handler.ServeHTTP(w, req)
	})
}

func (r *registerer) write(status int, response any, writer http.ResponseWriter) {
	var bytes []byte
	if responseBytes, ok := response.([]byte); ok {
		bytes = responseBytes
	} else {
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			writer.Header().Set("Content-Type", "text/plain")
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error())) // nolint
			return
		}
		bytes = jsonBytes
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(bytes) // nolint
}

func main() {}
