package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
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

	user, err := h.db.GetUserByAuthID(r.Context(), authID)
	if err != nil {
		redirect(w, "/user/new")
		return
	}

	_, err = session.GetUserID(r)
	if err != nil {
		session.SetUserID(w, r, user.ID)
	}

	pods := user.Edges.JoinedPods

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
		fmt.Fprintln(w, err)
		return
	}
	sess, err := session.GetSession(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sess.Values["authID"] = user.UserID
	sess.Values["email"] = user.Email
	sess.Save(r, w)
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

	user, err := h.db.CreateUser(r.Context(), authID, username, email)
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

// POST /{podID}/{beanID}/post
func (h *HandlerEnv) CreatePost(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	beanID, err := strconv.Atoi(chi.URLParam(r, "beanID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}

	p, err := h.db.CreatePost(r.Context(), userID, beanID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.pc.RenderFragment(w, "post.html", p)
}

// POST /pod/join
func (h *HandlerEnv) JoinPodHandler(w http.ResponseWriter, r *http.Request) {
	inviteCode := r.FormValue("invite")
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}

	_, err = h.db.JoinPod(r.Context(), userID, inviteCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", "/")
}

// GET /pod/{podID}
func (h *HandlerEnv) PodViewHandler(w http.ResponseWriter, r *http.Request) {
	podID := chi.URLParam(r, "podID")
	podIDInt, err := strconv.Atoi(podID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}
	beans, err := h.db.GetBeans(r.Context(), userID, podIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.pc.Render(w, "bean.html", beans[0])
}
