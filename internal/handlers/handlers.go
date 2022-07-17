package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

var (
	ErrURLNotFound = errors.New("url not found error")
	//	ErrURLConflict  = errors.New("urls conflict")
	//ErrNotUniqueURL = errors.New("not unique url")
)

type service interface {
	Shorten(ctx context.Context, url string, userID string) (string, error)
	Expand(ctx context.Context, id string) (string, error)
	FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error)
	ShortenBatch(ctx context.Context, originalURLs []models.OriginalURL, userID string) ([]models.UserURL, error)
	GetOriginURL(ctx context.Context, originURL string) (models.UserURL, error)
}

type auth interface {
	UserID(ctx context.Context) string
}

type pingService interface {
	Ping(ctx context.Context) bool
}

type handler struct {
	service     service
	auth        auth
	pingService pingService
}

func New(service service, userAuth auth, pingServ pingService) *handler {
	return &handler{
		service:     service,
		auth:        userAuth,
		pingService: pingServ,
	}
}

// Shorten Cut URL
func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	//var statusCode int
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to validate struct", 400)
		return
	}

	userID := h.auth.UserID(r.Context())

	url := string(bytes)
	shortcut, err := h.service.Shorten(r.Context(), url, userID)
	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) && pge.Code == pgerrcode.UniqueViolation {
			createURL, err := h.service.GetOriginURL(r.Context(), url)
			if err != nil {
				log.WithError(err).WithField("create url", createURL).Error("response error unique violation")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortcut))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//if errors.Is(err, ErrNotUniqueURL) {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if err != nil {
	//	statusCode = http.StatusConflict
	//} else {
	//	statusCode = http.StatusCreated
	//}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortcut))
	if err != nil {
		log.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}

// Expand Returns full URL by ID of shorted one
func (h *handler) Expand(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, err := h.service.Expand(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrURLNotFound) {
			http.Error(w, "url not found", http.StatusNoContent)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func (h *handler) APIJSONShorten(w http.ResponseWriter, r *http.Request) {
	//var statusCode int
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	req := ShortenRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	userID := h.auth.UserID(r.Context())

	shortcut, err := h.service.Shorten(r.Context(), req.URL, userID)
	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) && pge.Code == pgerrcode.UniqueViolation {
			createURL, err := h.service.GetOriginURL(r.Context(), req.URL)
			if err != nil {
				log.WithError(err).WithField("create url", createURL).Error("response error unique violation")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortcut))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//if err != nil {
	//	statusCode = http.StatusConflict
	//} else {
	//	statusCode = http.StatusCreated
	//}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := ShortenReply{ShortenURLResult: shortcut}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *handler) FetchURLs(w http.ResponseWriter, r *http.Request) {
	userID := h.auth.UserID(r.Context())
	urls, err := h.service.FetchURLs(r.Context(), userID)
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
	w.WriteHeader(http.StatusOK)

	resp := toGetUrlsReply(urls)
	body, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("marshal urls response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(body)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("write response error")
		http.Error(w, err.Error(), 500)
		return
	}

}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	if success := h.pingService.Ping(r.Context()); !success {
		http.Error(w, "ping database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req []ShortenBatchRequest
	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	userID := h.auth.UserID(r.Context())
	originalUrls := toShortenBatchRequest(req)

	urls, err := h.service.ShortenBatch(r.Context(), originalUrls, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := toShortenBatchReply(urls)
	marshal, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).WithField("urls", urls).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
