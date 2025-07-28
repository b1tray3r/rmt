package handler

import (
	"log/slog"
	"net/http"
)

type RedmineIssueHandler struct {
	logger *slog.Logger
}

func NewRedmineIssueHandler(logger *slog.Logger) *RedmineIssueHandler {
	return &RedmineIssueHandler{
		logger: logger,
	}
}

func (h *RedmineIssueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "...."}`))
}
