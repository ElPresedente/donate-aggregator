package widget

import (
	"go-back/l2wbridge"

	"github.com/gorilla/websocket"
)

type RouletteWidget struct {
	connection *websocket.Conn
	lbridge    l2wbridge.W2LHandler
}

func NewRouletteWidget(connection *websocket.Conn, lbridge l2wbridge.W2LHandler) *RouletteWidget {
	return &RouletteWidget{connection, lbridge}
}

func (rw *RouletteWidget) A2WRequest(request string, data string) {
	switch request {
	case "enqueue-spins":
		//do smth
		rw.connection.WriteMessage(websocket.TextMessage, []byte(data))
	}
}

func (rw *RouletteWidget) W2ARequest(request string, data string) {
	switch request {
	case "spins-done":
		rw.lbridge.LogicEventHandler("spins-done", "")
	}
}

func (*RouletteWidget) Close() {
	//do nothing
}
