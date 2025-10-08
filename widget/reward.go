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

type LogStr struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Value  string `json:"value"`
	Pinned bool   `json:"pinned"`
}

func (wh *WidgetsHub) NewRewardWidget(connection *websocket.Conn, lbridge l2wbridge.W2LHandler) *RewardWidget {
	lbridge.LogicEventHandler("rewardConnected", "")
	wh.rewardWidgets += 1
	return &RewardWidget{wh, connection, lbridge}
}

func (rw *RewardWidget) A2WRequest(request string, data string) {
	switch request {
	case "reward-set":
		var sendingData struct {
			Request string `json:"request"`
			Data    LogStr `json:"pin"`
		}
		var logData LogStr
		err := json.Unmarshal([]byte(data), &logData)
		if err != nil {
			log.Fatalf("Ошибка декодирования:", err)
			return
		}

		sendingData.Request = "set-text"
		sendingData.Data = logData
		marshalledData, err := json.Marshal(sendingData)
		if err != nil {
			log.Fatalf("json encoding error %v", err)
			return
		}
		rw.connection.WriteMessage(websocket.TextMessage, marshalledData)
	case "reward-reset":
		var sendingData struct {
			Request string `json:"request"`
			Data    LogStr `json:"pin"`
		}
		var logData LogStr
		err := json.Unmarshal([]byte(data), &logData)
		if err != nil {
			log.Fatalf("Ошибка декодирования:", err)
			return
		}

		sendingData.Request = "reset-text"
		sendingData.Data = logData
		marshalledData, err := json.Marshal(sendingData)
		if err != nil {
			log.Fatalf("json encoding error %v", err)
			return
		}
		rw.connection.WriteMessage(websocket.TextMessage, marshalledData)
	case "widget-reload":
		var itemsData []LogStr

		if len(data) > 0 {
			if err := json.Unmarshal([]byte(data), &itemsData); err != nil {
				log.Printf("Ошибка парсинга JSON: %v", err)
				return
			}
		}
		log.Printf("Получено %d элементов", len(itemsData))

		var sendingData struct {
			Request string   `json:"request"`
			Data    []LogStr `json:"pin"`
		}
		sendingData.Request = "widget-reload"
		sendingData.Data = itemsData

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
