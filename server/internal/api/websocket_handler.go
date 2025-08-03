package api

import (
	"chat/internal/store"
	"chat/internal/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // For now, allow all origins
	},
}

type WebSocketHandler struct {
	messageStore store.MessageStore
	userStore    store.UserStore
	logger       *log.Logger
	clients      map[int]*websocket.Conn
	clientsMutex sync.RWMutex
}

func NewWebSocketHandler(messageStore store.MessageStore, userStore store.UserStore, logger *log.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		messageStore: messageStore,
		userStore:    userStore,
		logger:       logger,
		clients:      make(map[int]*websocket.Conn),
	}
}

type WSMessage struct {
	Type       string `json:"type"`
	SenderID   int    `json:"sender_id,omitempty"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	Content    string `json:"content,omitempty"`
	Error      string `json:"error,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Printf("ERROR: user_id is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User ID is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Printf("ERROR: invalid user_id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	_, err = h.userStore.GetUserByID(userID)
	if err != nil {
		h.logger.Printf("ERROR: user not found: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "User not found"})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Printf("ERROR: upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	h.clientsMutex.Lock()
	h.clients[userID] = conn
	h.clientsMutex.Unlock()
	defer func() {
		h.clientsMutex.Lock()
		delete(h.clients, userID)
		h.clientsMutex.Unlock()
	}()
	h.logger.Printf("INFO: client connected: %d", userID)

	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				h.logger.Printf("INFO: client disconnected: %d", userID)
			} else {
				h.logger.Printf("ERROR: reading message: %v", err)
			}
			return
		}

		h.handleMessage(userID, &msg)
	}
}

func (h *WebSocketHandler) handleMessage(senderID int, msg *WSMessage) {
	switch msg.Type {
	case "send_message":
		h.handleSendMessage(senderID, msg)
	case "get_history":
		h.handleGetMessages(senderID, msg)
	default:
		h.handleInvalidMessage(senderID)
	}
}

func (h *WebSocketHandler) handleSendMessage(senderID int, msg *WSMessage) {
	if msg.ReceiverID == 0 {
		h.logger.Printf("ERROR: receiver_id is required")
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Receiver ID is required",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	if msg.Content == "" {
		h.logger.Printf("ERROR: content is required")
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Content is required",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	_, err := h.userStore.GetUserByID(msg.ReceiverID)
	if err != nil {
		h.logger.Printf("ERROR: receiver user not found: %v", err)
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Receiver user not found",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	_, err = h.messageStore.CreateMessage(senderID, msg.ReceiverID, msg.Content)
	if err != nil {
		h.logger.Printf("ERROR: creating message: %v", err)
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Failed to send message",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	// Send new_message to recipient
	h.clientsMutex.RLock()
	recipientConn, recipientExists := h.clients[msg.ReceiverID]
	senderConn, senderExists := h.clients[senderID]
	h.clientsMutex.RUnlock()

	// Get current timestamp
	currentTime := time.Now().Format(time.RFC3339)

	if recipientExists {
		response := WSMessage{
			Type:       "new_message",
			SenderID:   senderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  currentTime,
		}
		err := utils.WriteWebsocketMessage(recipientConn, response, h.logger)
		if err != nil {
			h.logger.Printf("ERROR: failed to send message to recipient: %v", err)
		}
	}

	// Send new_message to sender as well so they see their own message
	if senderExists {
		response := WSMessage{
			Type:       "new_message",
			SenderID:   senderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  currentTime,
		}
		err = utils.WriteWebsocketMessage(senderConn, response, h.logger)
		if err != nil {
			h.logger.Printf("ERROR: failed to send message to sender: %v", err)
		}
	}
}

func (h *WebSocketHandler) handleGetMessages(senderID int, msg *WSMessage) {
	if msg.ReceiverID == 0 {
		h.logger.Printf("ERROR: receiver_id is required")
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Receiver ID is required",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	messages, err := h.messageStore.GetMessagesBetweenUsers(senderID, msg.ReceiverID)
	if err != nil {
		h.logger.Printf("ERROR: getting messages: %v", err)
		h.clientsMutex.RLock()
		senderConn, exists := h.clients[senderID]
		h.clientsMutex.RUnlock()
		if exists {
			response := WSMessage{
				Type:  "error",
				Error: "Failed to get messages",
			}
			err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
			if err != nil {
				return
			}
		}
		return
	}

	h.clientsMutex.RLock()
	senderConn, exists := h.clients[senderID]
	h.clientsMutex.RUnlock()
	if exists {
		response := map[string]interface{}{
			"type":        "messages_history",
			"sender_id":   senderID,
			"receiver_id": msg.ReceiverID,
			"messages":    messages,
		}
		err = utils.WriteWebsocketMessage(senderConn, response, h.logger)
		if err != nil {
			return
		}
	}
}

func (h *WebSocketHandler) handleInvalidMessage(senderID int) {
	h.clientsMutex.RLock()
	senderConn, exists := h.clients[senderID]
	h.clientsMutex.RUnlock()
	if exists {
		response := WSMessage{
			Type:  "error",
			Error: "Invalid message type",
		}
		err := utils.WriteWebsocketMessage(senderConn, response, h.logger)
		if err != nil {
			return
		}
	}
}
