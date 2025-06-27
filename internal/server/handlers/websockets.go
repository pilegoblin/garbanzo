package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
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

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: h.origins,
	})

	if err != nil {
		slog.Error("failed to initialize websocket", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.CloseNow()

	uniqueID := ksuid.New()
	h.switchboard.RegisterUser(podID, beanID, userID, uniqueID, conn)
	defer h.switchboard.UnregisterUser(podID, beanID, uniqueID)

	for {
		_, message, err := conn.Read(r.Context())
		if err != nil {
			break
		}

		var data map[string]any
		if err := json.Unmarshal(message, &data); err != nil {
			slog.Error("failed to unmarshal message", "error", err)
			continue
		}

		content, ok := data["content"].(string)
		if !ok {
			continue
		}

		if len(content) == 0 || len(content) > 2048 {
			continue
		}

		messageID := ksuid.New()

		m, err := h.query.CreateMessage(r.Context(), sqlc.CreateMessageParams{
			ID:       messageID.String(),
			BeanID:   beanID,
			AuthorID: userID,
			Content:  content,
		})
		if err != nil {
			slog.Error("failed to create message", "error", err)
			continue
		}

		messageStringSender := h.pc.FragmentString("message.html", MessageData{
			ID:              m.ID,
			Content:         m.Content,
			AuthorUsername:  m.AuthorUsername,
			AuthorUserColor: m.AuthorUserColor,
			AuthorID:        m.AuthorID,
			CreatedAt:       m.CreatedAt,
			Action:          MessageActionNew,
			Editable:        true,
		})

		messageStringReceiver := h.pc.FragmentString("message.html", MessageData{
			ID:              m.ID,
			Content:         m.Content,
			AuthorUsername:  m.AuthorUsername,
			AuthorUserColor: m.AuthorUserColor,
			AuthorID:        m.AuthorID,
			CreatedAt:       m.CreatedAt,
			Action:          MessageActionNew,
			Editable:        false,
		})

		h.switchboard.SendMessageToOthers(r.Context(), podID, beanID, userID, messageStringReceiver)
		h.switchboard.SendMessage(r.Context(), podID, beanID, userID, messageStringSender)
	}

}
