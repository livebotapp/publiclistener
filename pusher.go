package publiclistener

import "context"

type Pusher interface {
	SendConfirmationAttempt(ctx context.Context, msg *ConfirmationAttempt) error
}
