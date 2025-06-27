package switchboard

import (
	"context"
	"log/slog"
	"sync"

	"github.com/coder/websocket"
	"github.com/segmentio/ksuid"
)

type UserSocket struct {
	conn   *websocket.Conn
	userID int64
}

type BeanManager struct {
	users map[ksuid.KSUID]*UserSocket
}

type PodManager struct {
	beans map[int64]*BeanManager
}

type Switchboard struct {
	pods map[int64]*PodManager
	rw   sync.RWMutex
}

func New() *Switchboard {
	return &Switchboard{
		pods: make(map[int64]*PodManager),
	}
}

func (sb *Switchboard) RegisterUser(podID, beanID, userID int64, uniqueID ksuid.KSUID, conn *websocket.Conn) {
	sb.rw.Lock()
	defer sb.rw.Unlock()

	pod, ok := sb.pods[podID]
	if !ok {
		pod = &PodManager{
			beans: make(map[int64]*BeanManager),
		}
		sb.pods[podID] = pod
	}

	bean, ok := pod.beans[beanID]
	if !ok {
		bean = &BeanManager{
			users: make(map[ksuid.KSUID]*UserSocket),
		}
		pod.beans[beanID] = bean
	}

	bean.users[uniqueID] = &UserSocket{
		conn:   conn,
		userID: userID,
	}

	slog.Info("registered user", "podID", podID, "beanID", beanID, "uniqueID", uniqueID)
}

func (sb *Switchboard) UnregisterUser(podID, beanID int64, uniqueID ksuid.KSUID) {
	sb.rw.Lock()
	defer sb.rw.Unlock()

	pod, ok := sb.pods[podID]
	if !ok {
		slog.Warn("unregistering user from non-existent pod", "podID", podID)
		return
	}

	bean, ok := pod.beans[beanID]
	if !ok {
		slog.Warn("unregistering user from non-existent bean", "beanID", beanID)
		return
	}

	delete(bean.users, uniqueID)

	slog.Info("unregistered user", "podID", podID, "beanID", beanID, "uniqueID", uniqueID)
}

// BroadcastMessage sends a message to every user in the bean
func (sb *Switchboard) BroadcastMessage(ctx context.Context, podID, beanID, userID int64, message string) {
	sb.rw.RLock()
	defer sb.rw.RUnlock()

	pod, ok := sb.pods[podID]
	if !ok {
		slog.Warn("broadcasting message to non-existent pod", "podID", podID)
		return
	}

	bean, ok := pod.beans[beanID]
	if !ok {
		slog.Warn("broadcasting message to non-existent bean", "beanID", beanID)
		return
	}

	slog.Info("broadcasting message", "podID", podID, "beanID", beanID, "userID", userID)

	for _, user := range bean.users {
		if err := user.conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
			slog.Error("failed to write message", "error", err)
		}
	}
}

// SendMessage sends a message to a specific user
func (sb *Switchboard) SendMessage(ctx context.Context, podID, beanID, userID int64, message string) {
	sb.rw.RLock()
	defer sb.rw.RUnlock()

	pod, ok := sb.pods[podID]
	if !ok {
		slog.Warn("sending message to non-existent pod", "podID", podID)
		return
	}

	bean, ok := pod.beans[beanID]
	if !ok {
		slog.Warn("sending message to non-existent bean", "beanID", beanID)
		return
	}

	for _, user := range bean.users {
		if user.userID != userID {
			continue
		}

		if err := user.conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
			slog.Error("failed to write message", "error", err)
		}
	}
}

// Sends message to everyone but one user
func (sb *Switchboard) SendMessageToOthers(ctx context.Context, podID, beanID, userID int64, message string) {
	sb.rw.RLock()
	defer sb.rw.RUnlock()

	pod, ok := sb.pods[podID]
	if !ok {
		slog.Warn("sending message to non-existent pod", "podID", podID)
		return
	}

	bean, ok := pod.beans[beanID]
	if !ok {
		slog.Warn("sending message to non-existent bean", "beanID", beanID)
		return
	}

	for _, user := range bean.users {
		if user.userID == userID {
			continue
		}

		if err := user.conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
			slog.Error("failed to write message", "error", err)
		}
	}
}
