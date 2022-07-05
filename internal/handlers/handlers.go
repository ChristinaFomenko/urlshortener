package handlers

import (
	"encoding/json"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type service interface {
	Shorten(url models.CreateURL) (string, error)
	Expand(id string) (string, error)
	GetList() ([]models.UserURL, error)
	Ping() error
}

type handler struct {
	service service
}

func New(service service) *handler {
	return &handler{
		service: service,
	}
}

// Shorten Cut URL
func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to validate struct", 400)
		return
	}

	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := string(bytes)
	id, err := h.service.Shorten(models.CreateURL{
		User: idCookie.Value,
		URL:  url,
	})
	if err != nil {
		log.WithError(err).WithField("url", url).Error("shorten url error")
		http.Error(w, "url shortcut", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		log.WithError(err).WithField("id", id).Error("write response error")
		return
	}
	http.SetCookie(w, idCookie)
}

// Expand Returns full URL by ID of shorted one
func (h *handler) Expand(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, err := h.service.Expand(id)
	if err != nil {
		http.Error(w, "url not found", http.StatusNoContent)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) APIJSONShorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	req := models.ShortenRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	id, err := h.service.Shorten(models.CreateURL{
		User: idCookie.Value,
		URL:  req.URL,
	})
	if err != nil {
		log.WithError(err).WithField("url", req.URL).Error("shorten url error")
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.ShortenReply{ShortenURLResult: id}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("write response error")
		return
	}
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetList()
	if err != nil {
		log.WithError(err).Error("get urls error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")

	resp := toGetUrlsReply(urls)
	body, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("marshal urls response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(body)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("write response error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.service.Ping()
	if err != nil {
		log.Infof("DB not avalable %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}
