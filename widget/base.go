package widget

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-back/l2wbridge"

	"github.com/gorilla/websocket"
)

type Widget interface {
	W2ARequest(request string, data string) //widget -> app
	A2WRequest(request string, data string) //app -> widget
	Close()
}

type WidgetsHub struct {
	upgrader          websocket.Upgrader
	LogicEventHandler l2wbridge.W2LHandler
	widgets           map[Widget]bool

	rouletteWidgets int
}

func (wh *WidgetsHub) WidgetEventHandler(request string, data string) {
	for widget := range wh.widgets {
		widget.A2WRequest(request, data)
	}
}

func NewWidgetsHub() WidgetsHub {
	return WidgetsHub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		LogicEventHandler: nil,
		widgets:           map[Widget]bool{},
		rouletteWidgets:   0,
	}
}

func (wh *WidgetsHub) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	widgetType := r.URL.Query().Get("type")

	log.Printf("NEW CONNECTION TYPE %s\n", widgetType)

	if widgetType == "" {
		http.Error(w, "Widget type not specified", http.StatusBadRequest)
		return
	}

	conn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("HTTP to websocket upgrade error: ", err)
		return
	}
	defer conn.Close()

	var currentWidget Widget
	switch widgetType {
	case "roulette":
		{
			currentWidget = wh.NewRouletteWidget(conn, wh.LogicEventHandler)
		}
	default:
		{
			log.Fatalf("Unknown widget type %s", widgetType)
			return
		}
	}
	wh.widgets[currentWidget] = true

	for {
		msgType, payload, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if msgType != websocket.TextMessage {
			//используем только json, не бинари
			continue
		}
		var baseRequest struct {
			Request string `json:"request"`
		}
		err = json.Unmarshal(payload, &baseRequest)
		if err != nil {
			log.Fatalf("!!! Illformed widget request: %s", string(payload))
		}
		currentWidget.W2ARequest(baseRequest.Request, string(payload))
	}
	currentWidget.Close()
	delete(wh.widgets, currentWidget)
}

func (wh *WidgetsHub) Start(addr string) error {
	http.HandleFunc("/ws", wh.ConnectionHandler)
	return http.ListenAndServe(addr, nil)
}
