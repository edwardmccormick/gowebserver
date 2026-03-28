package realtime

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type NotificationType string

const (
	NotificationTypeMessage NotificationType = "chat_message"
	NotificationTypeMatch   NotificationType = "match_update"
)

type NotificationEvent struct {
	Type      NotificationType `json:"type"`
	MatchID   uint             `json:"match_id"`
	Count     int              `json:"count"`
	Timestamp time.Time        `json:"timestamp"`
}

type NotificationCenter struct {
	clients     map[string]map[chan NotificationEvent]bool
	clientsLock sync.RWMutex
}

func NewNotificationCenter() *NotificationCenter {
	return &NotificationCenter{
		clients: make(map[string]map[chan NotificationEvent]bool),
	}
}

func (nc *NotificationCenter) Register(userID string, ch chan NotificationEvent) {
	nc.clientsLock.Lock()
	defer nc.clientsLock.Unlock()

	if nc.clients[userID] == nil {
		nc.clients[userID] = make(map[chan NotificationEvent]bool)
	}
	nc.clients[userID][ch] = true
}

func (nc *NotificationCenter) Unregister(userID string, ch chan NotificationEvent) {
	nc.clientsLock.Lock()
	defer nc.clientsLock.Unlock()

	if nc.clients[userID] != nil {
		delete(nc.clients[userID], ch)
		if len(nc.clients[userID]) == 0 {
			delete(nc.clients, userID)
		}
	}
}

func (nc *NotificationCenter) Broadcast(userID string, event NotificationEvent) {
	nc.clientsLock.RLock()
	defer nc.clientsLock.RUnlock()

	if nc.clients[userID] == nil {
		return
	}

	for clientCh := range nc.clients[userID] {
		select {
		case clientCh <- event:
		default:
			go func(ch chan NotificationEvent) {
				nc.clientsLock.Lock()
				delete(nc.clients[userID], ch)
				if len(nc.clients[userID]) == 0 {
					delete(nc.clients, userID)
				}
				nc.clientsLock.Unlock()
				close(ch)
			}(clientCh)
		}
	}
}

type Hub struct {
	rooms                 map[string]map[*websocket.Conn]bool
	roomsLock             sync.Mutex
	activeConnections     map[uint]*websocket.Conn
	activeConnectionsLock sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:             make(map[string]map[*websocket.Conn]bool),
		activeConnections: make(map[uint]*websocket.Conn),
	}
}

func (h *Hub) AddRoomConnection(roomID string, conn *websocket.Conn) {
	h.roomsLock.Lock()
	defer h.roomsLock.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[roomID][conn] = true
}

func (h *Hub) RemoveRoomConnection(roomID string, conn *websocket.Conn) bool {
	h.roomsLock.Lock()
	defer h.roomsLock.Unlock()

	if h.rooms[roomID] == nil {
		return true
	}

	delete(h.rooms[roomID], conn)
	return len(h.rooms[roomID]) == 0
}

func (h *Hub) RoomConnections(roomID string) map[*websocket.Conn]bool {
	h.roomsLock.Lock()
	defer h.roomsLock.Unlock()

	connections := make(map[*websocket.Conn]bool, len(h.rooms[roomID]))
	for conn, present := range h.rooms[roomID] {
		connections[conn] = present
	}
	return connections
}

func (h *Hub) RemoveRoomClient(roomID string, conn *websocket.Conn) {
	h.roomsLock.Lock()
	defer h.roomsLock.Unlock()

	if h.rooms[roomID] != nil {
		delete(h.rooms[roomID], conn)
	}
}

func (h *Hub) SetActiveConnection(userID uint, conn *websocket.Conn) {
	h.activeConnectionsLock.Lock()
	defer h.activeConnectionsLock.Unlock()
	h.activeConnections[userID] = conn
}

func (h *Hub) RemoveActiveConnection(userID uint) {
	h.activeConnectionsLock.Lock()
	defer h.activeConnectionsLock.Unlock()
	delete(h.activeConnections, userID)
}

func (h *Hub) HasActiveConnection(userID uint) bool {
	h.activeConnectionsLock.RLock()
	defer h.activeConnectionsLock.RUnlock()
	_, exists := h.activeConnections[userID]
	return exists
}
