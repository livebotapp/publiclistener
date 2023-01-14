package publiclistener

type ConfirmationAttempt struct {
	Code         string `json:"code"`
	KickID       string `json:"kick_id"`
	KickUsername string `json:"kick_username"`
}
