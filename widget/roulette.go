package widget

import (
	"encoding/json"
	"go-back/l2wbridge"
	"log"

	"github.com/gorilla/websocket"
)

type RouletteWidget struct {
	hub        *WidgetsHub
	connection *websocket.Conn
	lbridge    l2wbridge.W2LHandler
}

func (wh *WidgetsHub) NewRouletteWidget(connection *websocket.Conn, lbridge l2wbridge.W2LHandler) *RouletteWidget {
	lbridge.LogicEventHandler("rouletteConnected", "")
	wh.rouletteWidgets += 1
	return &RouletteWidget{wh, connection, lbridge}
}

func (rw *RouletteWidget) A2WRequest(request string, data string) {
	switch request {
	case "enqueue-spins":
		//do smth
		rw.connection.WriteMessage(websocket.TextMessage, []byte(data))
	case "reloadRoulette":
		var data struct {
			Request string `json:"request"`
		}
		data.Request = "reset"
		marshalledData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("json encoding error %v", err)
			return
		}
		rw.connection.WriteMessage(websocket.TextMessage, marshalledData)
	}
}

func (rw *RouletteWidget) W2ARequest(request string, data string) {
	switch request {
	case "spins-done":
		rw.lbridge.LogicEventHandler("spins-done", "")
	}
}

func (rw *RouletteWidget) Close() {
	//do nothing
	rw.hub.rouletteWidgets -= 1
	if rw.hub.rouletteWidgets == 0 {
		rw.lbridge.LogicEventHandler("rouletteDisconnected", "")
	}
}
