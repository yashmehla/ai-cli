package chat

type Message struct {
	Role    string
	Content string
}

type History struct {
	Messages []Message
}

func NewHistory() *History {
	return &History{
		Messages: []Message{},
	}
}

func (h *History) Add(role, content string) {

	h.Messages = append(h.Messages, Message{
		Role:    role,
		Content: content,
	})
}
