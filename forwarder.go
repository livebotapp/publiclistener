package publiclistener

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Forwarder struct {
	pusher Pusher
}

func NewForwarder(ctx context.Context, pusher Pusher) *Forwarder {
	return &Forwarder{
		pusher: pusher,
	}
}

func (f *Forwarder) ForwardMessage(ctx context.Context, rawMessage []byte) error {
	receivedEvent := new(ReceivedEvent)
	err := json.Unmarshal(rawMessage, receivedEvent)

	if err != nil {
		fmt.Printf("Err unmarshal : %+v\n", err)
		return err
	}

	if receivedEvent.Event == ChatMessageSent {
		data := new(ChatMessageSentData)

		err = json.Unmarshal([]byte(receivedEvent.Data), &data)
		if err != nil {
			fmt.Printf("Err unmarshal data : %+v\n", err)
			return err
		}

		forwardedMessage := &Message{
			Sender:  strings.ToLower(data.User.Username),
			Content: data.Message.Message,
		}

		if strings.HasPrefix(forwardedMessage.Content, "!verify") {
			fmt.Printf("Verify command detected : %+v\n", forwardedMessage)
			cutString := strings.Split(forwardedMessage.Content, " ")
			if len(cutString) < 2 {
				return nil
			}
			f.pusher.SendConfirmationAttempt(ctx, &ConfirmationAttempt{
				Code:         cutString[1],
				KickID:       fmt.Sprintf("%d", data.User.ID),
				KickUsername: forwardedMessage.Sender,
			})
			return nil
		}
	}
	return nil
}
