package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func WriteWebsocketMessage(conn *websocket.Conn, data any, logger *log.Logger) error {
	err := conn.WriteJSON(data)
	if err != nil {
		logger.Printf("ERROR: writing message: %v", err)
		return err
	}
	return nil
}
