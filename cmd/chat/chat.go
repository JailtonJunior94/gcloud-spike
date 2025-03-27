package chat

import (
	"context"
	"fmt"

	"google.golang.org/api/chat/v1"

	"google.golang.org/api/option"
)

type spike struct {
	service *chat.Service
}

const prefixLog = "upload: "

func NewChat(credentials []byte) (*spike, error) {
	service, err := chat.NewService(context.Background(), option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, fmt.Errorf("%s%v", prefixLog, err)
	}
	return &spike{service: service}, nil
}

func (s *spike) SendMessage(space string, message string) error {
	spaces, err := s.service.Spaces.Messages.List("spaces/" + space).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}
	if len(spaces.Messages) > 0 {
		return fmt.Errorf("%s%v", prefixLog, "space is not empty")
	}

	msg := &chat.Message{
		Text: message,
	}
	_, err = s.service.Spaces.Messages.Create(space, msg).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}
	return nil
}
