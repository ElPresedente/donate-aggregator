package sources

import (
  "fmt"
  "net/http"
  "github.com/gorilla/websocket"
  "encoding/json"
  "sync"
)

//Тут пока что насрано

type WebSocketHub struct { //чо такое хаб?
	clients	map[*websocket.Conn]bool
	clientsMu sync.Mutex
	upgrader websocket.Upgrader
}

func NewWebSocketHub() *WebSocketHub {
  return &WebSocketHub{
    clients: make(map[*websocket.Conn]bool),
    upgrader: websocket.Upgrader{
      CheckOrigin: func(r *http.Request) bool {
        return true
      },
    },
  }
}

func (hub *WebSocketHub) Start() {
  http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    conn, err := hub.upgrader.Upgrade(w, r, nil)
    if err != nil {
      fmt.Println("Ошибка апгрейда:", err)
      return
    }

    hub.clientsMu.Lock()
    hub.clients[conn] = true
    hub.clientsMu.Unlock()

    defer func() {
      hub.clientsMu.Lock()
      delete(hub.clients, conn)
      hub.clientsMu.Unlock()
      conn.Close()
    }()

    for {
      if _, _, err := conn.ReadMessage(); err != nil {
        break
      }
    }
  })

  go http.ListenAndServe(":8080", nil)
}


type Message struct {
  Type    string `json:"type"`
  Payload string `json:"payload"`
}

func (hub *WebSocketHub) SendToAll(msg Message) {
  data, _ := json.Marshal(msg)

  hub.clientsMu.Lock()
  defer hub.clientsMu.Unlock()

  for conn := range hub.clients {
    if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
      fmt.Println("Ошибка при отправке:", err)
      conn.Close()
      delete(hub.clients, conn)
    }
  }
}