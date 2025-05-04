// Package http реализует HTTP-обработчики и маршруты для взаимодействия с внешними клиентами.
package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers"
)

type Server struct {
	http.Server
	service controllers.MusicService
}

func NewServer(s controllers.MusicService) *Server {
	serv := &Server{
		Server: http.Server{
			Addr:    config.Config.HTTPPort,
			Handler: nil,
		},
		service: s,
	}
	serv.setupRoutes()
	return serv
}

func (s *Server) setupRoutes() {
	r := mux.NewRouter()
	r.Use(logMW)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ping pong\n")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.HandleFunc("/music", s.addMusic).Methods(http.MethodPost)
	r.HandleFunc("/music/{id}", s.getMusic).Methods(http.MethodGet)

	s.Server.Handler = r
}
