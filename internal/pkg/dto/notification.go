package dto

type FCMMessage struct {
	To      string         `json:"to"`
	Message FCMMessageData `json:"notification"`
}

type FCMMessageData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
