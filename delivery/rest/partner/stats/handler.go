package stats

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gitlab.com/cpanova/excentral/domain/partnerstats"
	"gitlab.com/cpanova/excentral/domain/sender"
)

// StatsHandler ...
type StatsHandler interface {
	Stats(http.ResponseWriter, *http.Request)
}

type statsHandler struct {
	senderRepo       sender.Repo
	partnerstatsRepo partnerstats.Repo
}

// NewHandler ...
func NewHandler(
	senderRepo sender.Repo,
	partnerstatsRepo partnerstats.Repo,
) StatsHandler {
	return &statsHandler{
		senderRepo,
		partnerstatsRepo,
	}
}

type dailyReportResponse []partnerstats.DailyReport

func (h *statsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	// get dates
	// layout := "2006-01-02"

	// read sender_id
	senderIDStr := chi.URLParam(r, "senderID")
	// get pid by sender_id
	senderID, err := strconv.ParseUint(senderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Sender ID", http.StatusBadRequest)
		return
	}

	sender, err := h.senderRepo.Get(uint(senderID))
	if err != nil {
		http.Error(w, "Please contact your manager", http.StatusUnauthorized)
		return
	}

	dailyReport, err := h.partnerstatsRepo.ByDay(sender.PID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(dailyReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}
