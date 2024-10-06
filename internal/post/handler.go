package post

import (
	"elasticsearch/internal/pkg/domain"
	"elasticsearch/internal/pkg/storage"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service
}

func New(storage storage.PostStorer) Handler {
	return Handler{
		service: service{storage: storage},
	}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}

	res, err := h.service.create(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	req.ID = httprouter.ParamsFromContext(r.Context()).ByName("id")
	res, err := h.service.update(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	res, err := h.service.delete(r.Context(), deleteRequest{ID: id})
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

func (h Handler) Find(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	res, err := h.service.find(r.Context(), findRequest{ID: id})
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}
