package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	"github.com/pilegoblin/garbanzo/db/sqlc"
	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
	"github.com/pilegoblin/garbanzo/internal/session"
)

// GET /
func (h *HandlerEnv) IndexViewHandler(w http.ResponseWriter, r *http.Request) {
	authID, ok := gbzocontext.GetAuthID(r.Context())
	if !ok {
		redirect(w, "/login")
		return
	}

	user, err := h.query.GetUserByAuthID(r.Context(), authID)
	if err != nil {
		redirect(w, "/user/new")
		return
	}

	session.SetUserID(w, r, user.ID)

	pods, err := h.query.ListPodsForUser(r.Context(), user.ID)
	if err != nil {
		slog.Error("failed to get pods for user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(pods) == 0 {
		h.pc.Render(w, "join_pod.html", nil)
		return
	}

	h.pc.Render(w, "pods.html", pods)
}

// GET /login
func (h *HandlerEnv) LoginViewHandler(w http.ResponseWriter, r *http.Request) {
	h.pc.Render(w, "login.html", nil)
}

// GET /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.Logout(w, r)
	redirect(w, "/login")
}

// GET /auth
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// GET /auth/callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error("failed to complete user auth", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess, err := session.GetSession(r)
	if err != nil {
		slog.Error("failed to get session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess.Values["authID"] = user.UserID
	sess.Values["email"] = user.Email

	err = sess.Save(r, w)
	if err != nil {
		slog.Error("failed to save session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirect(w, "/")
}

// POST /user/new
func (h *HandlerEnv) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		redirect(w, "/")
		return
	}
	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authID, err := session.GetAuthID(r)
	if err != nil {
		slog.Error("failed to get authID", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email, err := session.GetEmail(r)
	if err != nil {
		slog.Error("failed to get email", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.query.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username:  username,
		AuthID:    authID,
		Email:     email,
		UserColor: createUserColor(username),
	})

	if err != nil {
		slog.Error("failed to create user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.SetUserID(w, r, user.ID)
	w.Header().Set("HX-Location", "/")
}

// GET /user/new
func (h *HandlerEnv) NewUserViewHandler(w http.ResponseWriter, r *http.Request) {
	h.pc.Render(w, "new_user.html", nil)
}

// POST /pod/join
func (h *HandlerEnv) JoinPodHandler(w http.ResponseWriter, r *http.Request) {
	inviteCode := r.FormValue("invite")
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}

	pod, err := h.query.GetPodByInviteCode(r.Context(), pgtype.Text{String: inviteCode, Valid: true})
	if err != nil {
		slog.Error("failed to get pod by invite code", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//_, err = h.db.JoinPod(r.Context(), userID, inviteCode)
	_, err = h.query.AddPodMember(r.Context(), sqlc.AddPodMemberParams{
		UserID: userID,
		PodID:  pod.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", "/")
}

// GET /pod/{podID}
func (h *HandlerEnv) PodViewHandler(w http.ResponseWriter, r *http.Request) {
	podIDParam := chi.URLParam(r, "podID")
	podID, err := strconv.ParseInt(podIDParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}

	check, err := h.query.CheckUserInPod(r.Context(), sqlc.CheckUserInPodParams{
		UserID: userID,
		PodID:  podID,
	})

	if !check || err != nil {
		redirect(w, "/")
		return
	}

	beans, err := h.query.ListBeansForPodFull(r.Context(), podID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []BeanWithMessages

	// Doing some shenanigans because the sql query returns json for messages
	for _, bean := range beans {
		var messages []MessageData
		bytes, err := json.Marshal(bean.Messages)
		if err != nil {
			slog.Error("failed to marshal messages", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(bytes, &messages)
		if err != nil {
			slog.Error("failed to unmarshal messages", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = append(result, BeanWithMessages{
			ID:       bean.ID,
			Name:     bean.Name,
			PodID:    bean.PodID,
			PodName:  bean.PodName,
			Messages: messages,
		})

	}
	bean := result[0]

	for i := range bean.Messages {
		bean.Messages[i].Editable = bean.Messages[i].AuthorID == userID
	}

	// TODO: Handle multiple beans
	h.pc.Render(w, "bean.html", bean)
}

// GET /messages/edit/{messageID}
// Retreives the template for editing a message
func (h *HandlerEnv) EditMessageViewHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	message, err := h.query.GetMessage(r.Context(), messageID)
	if err != nil {
		slog.Error("failed to get message", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.pc.RenderFragment(w, "edit_message.html", message)
}

// POST /messages/edit/{messageID}
// Updates a message
func (h *HandlerEnv) EditMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	m, err := h.query.UpdateMessage(r.Context(), sqlc.UpdateMessageParams{
		ID:      messageID,
		Content: content,
	})
	if err != nil {
		slog.Error("failed to update message", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageStringSender := h.pc.FragmentString("message.html", MessageData{
		ID:              m.ID,
		Content:         m.Content,
		AuthorUsername:  m.AuthorUsername,
		AuthorUserColor: m.AuthorUserColor,
		AuthorID:        m.AuthorID,
		CreatedAt:       m.CreatedAt,
		Action:          MessageActionEdit,
		Editable:        true,
	})

	messageStringReceiver := h.pc.FragmentString("message.html", MessageData{
		ID:              m.ID,
		Content:         m.Content,
		AuthorUsername:  m.AuthorUsername,
		AuthorUserColor: m.AuthorUserColor,
		AuthorID:        m.AuthorID,
		CreatedAt:       m.CreatedAt,
		Action:          MessageActionEdit,
		Editable:        false,
	})

	h.switchboard.SendMessageToOthers(r.Context(), m.PodID, m.BeanID, m.AuthorID, messageStringReceiver)
	h.switchboard.SendMessage(r.Context(), m.PodID, m.BeanID, m.AuthorID, messageStringSender)
}
