package publiclistener

const (
	ChatMessageSent = `App\Events\ChatMessageSentEvent`
)

type Message struct {
	Sender  string `json:"sender"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ReceivedEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// MESSAGE SENT
type ChatMessageSentData struct {
	Message ChatMessageSentDataMessage `json:"message"`
	User    ChatMessageSentDataUser    `json:"user"`
}

type ChatMessageSentDataMessage struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	ChatroomID string `json:"chatroom_id"`
}

type ChatMessageSentDataUser struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Role     *string `json:"role"`
}
