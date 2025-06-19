package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pilegoblin/garbanzo/db/sqlc"
	"github.com/pilegoblin/garbanzo/internal/session"
	"github.com/segmentio/ksuid"
)

// GET /ws/{podID}/{beanID}
func (h *HandlerEnv) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	podIDParam := chi.URLParam(r, "podID")
	beanIDParam := chi.URLParam(r, "beanID")

	podID, err := strconv.ParseInt(podIDParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	beanID, err := strconv.ParseInt(beanIDParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := session.GetUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Currently users in a pod can see all beans in the pod
	check, err := h.query.CheckUserInPod(r.Context(), sqlc.CheckUserInPodParams{
		UserID: userID,
		PodID:  podID,
	})

	if !check || err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to websocket", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	uniqueID := ksuid.New()
	h.switchboard.RegisterUser(podID, beanID, uniqueID, conn)
	defer h.switchboard.UnregisterUser(podID, beanID, uniqueID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]any
		if err := json.Unmarshal(message, &data); err != nil {
			slog.Error("failed to unmarshal message", "error", err)
			continue
		}

		slog.Info("received message on websocket", "message", data)
	}

}
