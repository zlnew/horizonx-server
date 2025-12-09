package rest

import (
	"net/http"
	"strconv"

	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"
)

type UserHandler struct {
	svc domain.UserService
	cfg *config.Config
}

func NewUserHandler(svc domain.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		svc: svc,
		cfg: cfg,
	}
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	search := q.Get("search")
	isPaginate := q.Get("paginate") == "true"

	opts := domain.ListOptions{
		Page:       page,
		Limit:      limit,
		Search:     search,
		IsPaginate: isPaginate,
	}

	result, err := h.svc.Get(r.Context(), opts)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    result.Data,
		Meta:    result.Meta,
	})
}

func (h *UserHandler) Store(w http.ResponseWriter, r *http.Request) {}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *UserHandler) Destroy(w http.ResponseWriter, r *http.Request) {}
