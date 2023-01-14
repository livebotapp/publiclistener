package kick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/livebotapp/publiclistener"
)

const (
	VerifbotChannelID  = 320238
	VerifbotChatroomID = 320236
)

type Consumer struct {
	forwarder *publiclistener.Forwarder
	wsAddr    string
	wsPath    string
}

func NewConsumer(ctx context.Context, wsAddr string, wsPath string,
	forwarder *publiclistener.Forwarder) *Consumer {
	return &Consumer{
		forwarder: forwarder,
		wsAddr:    wsAddr,
		wsPath:    wsPath,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	u := url.URL{Scheme: "wss", Host: c.wsAddr, Path: c.wsPath, RawQuery: "protocol%3D7"}

	co, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.runSubscribeChannelFlow(ctx, co)

	go func() {
		for {
			_, message, err := co.ReadMessage()
			if err != nil {
				fmt.Println("Read error : ", err)
				return
			}
			fmt.Printf("Received message : %s\n", message)
			c.forwarder.ForwardMessage(ctx, message)
		}
	}()

	return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
	return nil
}

type RegisterEvent struct {
	Event string            `json:"event"`
	Data  RegisterEventData `json:"data"`
}

type RegisterEventData struct {
	Auth    string `json:"auth"`
	Channel string `json:"channel"`
}

func (c *Consumer) runSubscribeChannelFlow(ctx context.Context, co *websocket.Conn) error {
	channelLivebotRegisterMsg, err := c.GenerateRegisterChannelPayload(ctx, VerifbotChannelID)
	err = co.WriteMessage(websocket.TextMessage, channelLivebotRegisterMsg)
	if err != nil {
		fmt.Printf("[KICK CONSUMER] - Error registering channel : %+v\n", err)
		return err
	}
	chatroomLivebotRegisterMsg, err := c.GenerateRegisterChatroomPayload(ctx, VerifbotChatroomID)
	err = co.WriteMessage(websocket.TextMessage, chatroomLivebotRegisterMsg)
	if err != nil {
		fmt.Printf("[KICK CONSUMER] - Error registering chatroom : %+v\n", err)
		return err
	}
	return nil
}

func (c *Consumer) GenerateRegisterChannelPayload(ctx context.Context, channelID int) ([]byte, error) {
	evt := &RegisterEvent{
		Event: "pusher:subscribe",
		Data: RegisterEventData{
			Auth:    "",
			Channel: fmt.Sprintf("channel.%d", channelID),
		},
	}

	return json.Marshal(evt)
}

func (c *Consumer) GenerateRegisterChatroomPayload(ctx context.Context, chatroomID int) ([]byte, error) {
	evt := &RegisterEvent{
		Event: "pusher:subscribe",
		Data: RegisterEventData{
			Auth:    "",
			Channel: fmt.Sprintf("chatrooms.%d", chatroomID),
		},
	}

	return json.Marshal(evt)
}
