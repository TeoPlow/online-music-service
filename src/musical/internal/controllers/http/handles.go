package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

func getID(req *http.Request) (uuid.UUID, error) {
	idStr := mux.Vars(req)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s *Server) addMusic(w http.ResponseWriter, req *http.Request) {
	var track models.Track
	d := json.NewDecoder(req.Body)
	if err := d.Decode(&track); err != nil {
		http.Error(w, fmt.Sprintf("json.Decode: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.service.AddMusic(req.Context(), track); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getMusic(w http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid ID: %s", err), http.StatusBadRequest)
		return
	}
	track, err := s.service.GetMusic(req.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		return
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(track); err != nil {
		http.Error(
			w, fmt.Sprintf("json.Encode: %s", err),
			http.StatusInternalServerError)
		return
	}
}
