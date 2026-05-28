package message

type MessageRequest struct {
	FromUID string `json:"from_uid"`
	ToUID   string `json:"to_uid"`
	Content string `json:"content"`
}
