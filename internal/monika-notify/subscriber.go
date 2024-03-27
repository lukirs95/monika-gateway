package monikanotify

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

var upgrader = websocket.Upgrader{}

type subscriber struct {
	conn        *websocket.Conn
	updateScope []types.DeviceId
	openErrors  map[int64]*types.PubError
}

func newSubscriber(openErrors map[int64]*types.PubError) *subscriber {
	return &subscriber{
		conn:        nil,
		updateScope: make([]types.DeviceId, 0),
		openErrors:  openErrors,
	}
}

func (sub *subscriber) upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
	}
	defer conn.Close()

	sub.conn = conn

	for _, openError := range sub.openErrors {
		sub.notifyError(openError)
	}
	sub.openErrors = nil

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("websocket read failed: ", err)
			break
		}
		if mt != websocket.TextMessage {
			break
		}

		if err := sub.scopeUpdate(message); err != nil {
			log.Print(err)
		}
	}

}

func (sub *subscriber) scopeUpdate(raw []byte) error {
	var newScope []types.DeviceId
	err := json.Unmarshal(raw, &newScope)
	if err != nil {
		return err
	}

	sub.updateScope = newScope
	return nil
}

func (sub *subscriber) notifyUpdate(device types.DeviceUpdate) error {
	found := false
	for _, deviceName := range sub.updateScope {
		if deviceName == device.Id {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	update := &report{
		Type:  ReportType_Update,
		Value: device,
	}
	return sub.conn.WriteJSON(update)
}

func (sub *subscriber) notifyError(newError *types.PubError) error {
	msg := &report{
		Type:  ReportType_Error,
		Value: newError,
	}

	return sub.conn.WriteJSON(msg)
}

func (sub *subscriber) deleteError(oldError *types.PubError) error {
	msg := &report{
		Type:  ReportType_ErrorRelease,
		Value: oldError.ErrorId,
	}

	return sub.conn.WriteJSON(msg)
}
