package redis

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/livebotapp/publiclistener"
)

type Pusher struct {
	connection               *asynq.Client
	confirmationAttemptQueue string
}

func NewPusher(ctx context.Context,
	redisAddr string,
	confirmationAttemptQueue string,
) (*Pusher, error) {
	connection := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &Pusher{
		connection:               connection,
		confirmationAttemptQueue: confirmationAttemptQueue,
	}, nil
}

func (c *Pusher) SendConfirmationAttempt(ctx context.Context, msg *publiclistener.ConfirmationAttempt) error {
	evt, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	task := asynq.NewTask(c.confirmationAttemptQueue, evt)

	_, err = c.connection.Enqueue(task, asynq.Queue(c.confirmationAttemptQueue))
	return err
}
