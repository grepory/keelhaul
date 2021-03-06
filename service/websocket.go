package service

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opsee/basic/schema"
	"github.com/opsee/keelhaul/bus"
	log "github.com/opsee/logrus"
	"github.com/opsee/vaper"
)

type websocketHandler struct {
	userChan  chan *schema.User
	closeChan chan struct{}
	ws        *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func (s *service) websocketHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("problem connecting to websocket")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler := websocketHandler{
		userChan:  make(chan *schema.User),
		closeChan: make(chan struct{}, 1),
		ws:        ws,
	}

	go handler.readLoop(s)
	handler.writeLoop(s)
}

func (handler websocketHandler) readLoop(s *service) {
	defer func() {
		handler.ws.Close()
		close(handler.userChan)
		close(handler.closeChan)
	}()

	for {
		_, reader, err := handler.ws.NextReader()
		if err != nil {
			log.WithError(err).Error("error reading from websocket, closing")
			handler.closeChan <- struct{}{}
			return
		}

		message := &bus.Message{}
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(message)
		if err != nil {
			log.WithError(err).Error("problem decoding incoming websocket message, discarding")
			continue
		}

		switch message.Command {
		case "authenticate":
			token, ok := message.Attributes["token"].(string)
			if !ok {
				log.Errorf("no token sent in authenticate request: %#v", message.Attributes)
				continue
			}

			decodedToken, err := vaper.Unmarshal(token)
			if err != nil {
				log.Errorf("couldn't unmarshal token: %s", token)
				continue
			}

			user := &schema.User{}
			err = decodedToken.Reify(user)
			if err != nil {
				log.WithError(err).Error("failed to decode bearer user token")
				continue
			}

			log.WithFields(log.Fields{
				"customer-id": user.CustomerId,
				"user-id":     user.Id,
			}).Info("authenticated user via websocket")

			handler.userChan <- user

		case "subscribe":
		default:
			log.Warnf("unrecognized command sent on websocket: %s", message.Command)
		}
	}
}

func (handler websocketHandler) writeLoop(s *service) {
	var (
		user *schema.User
		ok   bool
	)

	select {
	case user, ok = <-handler.userChan:
		if !ok {
			return
		}

	case <-handler.closeChan:
		log.Info("websocket received close from client")
		return

	case <-time.After(5 * time.Second):
		log.Warn("websocket closed waiting for authentication command")
		return
	}

	sub := s.bus.Subscribe(user)
	defer s.bus.Unsubscribe(sub)

	heartbeat := time.NewTicker(10 * time.Second)
	defer heartbeat.Stop()

	bastionBeat := time.NewTicker(30 * time.Second)
	defer bastionBeat.Stop()

	// initial messages
	err := handler.ws.WriteJSON(s.bastionMessage(user))
	if err != nil {
		log.WithError(err).Error("error sending to websocket")
		return
	}

	for {
		select {
		case <-handler.closeChan:
			log.Info("websocket received close from client")
			return

		case <-bastionBeat.C:
			err = handler.ws.WriteJSON(s.bastionMessage(user))
			if err != nil {
				log.WithError(err).Error("error sending to websocket")
				return
			}

		case bmsg := <-sub.Chan:
			err = handler.ws.WriteJSON(bmsg)
			if err != nil {
				log.WithError(err).Error("error sending to websocket")
				return
			}

		case t := <-heartbeat.C:
			err = handler.ws.WriteJSON(&bus.Message{
				Command:    "heartbeat",
				State:      "ok",
				CustomerID: user.CustomerId,
				Attributes: map[string]interface{}{
					"time": t,
				},
			})
			if err != nil {
				log.WithError(err).Error("error sending to websocket")
				return
			}

		case <-time.After(time.Minute):
			log.Warn("websocket timed out")
			return
		}
	}
}

func (s *service) bastionMessage(user *schema.User) *bus.Message {
	var msg *bus.Message

	bastionsResponse, err := s.ListBastions(user, &ListBastionsRequest{
		State: []string{"active"},
	})

	if err != nil {
		log.WithError(err).Error("failed listing bastions")

		msg = &bus.Message{
			Command:    "bastions",
			State:      "error",
			CustomerID: user.CustomerId,
			Message:    "error listing bastions",
		}
	} else {
		msg = &bus.Message{
			Command:    "bastions",
			State:      "ok",
			CustomerID: user.CustomerId,
			Attributes: map[string]interface{}{
				"bastions": bastionsResponse.Bastions,
			},
		}
	}

	return msg
}

func decodeBasicToken(token string, user *schema.User) error {
	jsonblob, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonblob, user)
	if err != nil {
		return err
	}

	err = user.Validate()
	if err != nil {
		return err
	}

	return nil
}
