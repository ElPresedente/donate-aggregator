package widget

import (
	"encoding/json"
	"go-back/l2wbridge"
	"log"

	"github.com/gorilla/websocket"
)

type RewardWidget struct {
	hub        *WidgetsHub
	connection *websocket.Conn
	lbridge    l2wbridge.W2LHandler
}

func (wh *WidgetsHub) NewRewardWidget(connection *websocket.Conn, lbridge l2wbridge.W2LHandler) *RewardWidget {
	lbridge.LogicEventHandler("rewardConnected", "")
	wh.rewardWidgets += 1
	return &RewardWidget{wh, connection, lbridge}
}

func (rw *RewardWidget) A2WRequest(request string, data string) {
	switch request {
	case "reward-set-text":
		var sendingData struct {
			Request string `json:"request"`
			Text    string `json:"text"`
		}
		sendingData.Request = "set-text"
		sendingData.Text = data
		marshalledData, err := json.Marshal(sendingData)
		if err != nil {
			log.Fatalf("json encoding error %v", err)
			return
		}
		rw.connection.WriteMessage(websocket.TextMessage, marshalledData)
	case "reward-reset":
		var sendingData struct {
			Request string `json:"request"`
		}
		sendingData.Request = "reset"
		marshalledData, err := json.Marshal(sendingData)
		if err != nil {
			log.Fatalf("json encoding error %v", err)
			return
		}
		rw.connection.WriteMessage(websocket.TextMessage, marshalledData)
	}
}

func (rw *RewardWidget) W2ARequest(request string, data string) {

}

func (rw *RewardWidget) Close() {
	//do nothing
	rw.hub.rewardWidgets -= 1
	if rw.hub.rewardWidgets == 0 {
		rw.lbridge.LogicEventHandler("rewardDisconnected", "")
	}
}
