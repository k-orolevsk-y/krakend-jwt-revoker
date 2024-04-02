package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/log"
	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/naming"
	"github.com/k-orolevsk-y/krakend-jwt-revoker/pkg/revoker"
)

type Server struct {
	router *mux.Router
	srv    *http.Server

	revoker *revoker.Revoker
	logger  log.ILogger
}

func NewServer(addr string, revoker *revoker.Revoker, logger log.ILogger) *Server {
	router := mux.NewRouter()
	server := &Server{
		router: router,
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},

		revoker: revoker,
		logger:  logger,
	}
	router.HandleFunc("/", server.AddRevokeToken).Methods("POST")

	return server
}

func (s *Server) AddRevokeToken(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewResponse(http.StatusBadRequest, "Invalid body")) // nolint
		return
	}

	items := make(map[string]string)
	for _, item := range body {
		if item.Key == "" || item.Value == "" {
			continue
		}

		items[item.Key] = item.Value
	}

	if len(items) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewResponse(http.StatusBadRequest, "Items not provided")) // nolint
		return
	}

	if err := s.revoker.Add(items); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewResponse(http.StatusInternalServerError, "Unknown error, try later...")) // nolint
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(NewResponse(http.StatusCreated, "")) // nolint
}

func (s *Server) Run() {
	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal(fmt.Sprintf("[PLUGIN %s] Failed listen and serve revoke server: %s", naming.PluginName, err))
		}
	}()
}
