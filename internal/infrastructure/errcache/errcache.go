package errcache

import (
	"context"
	"fmt"
	"sync"

	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
)

type ErrorCache struct {
	mu       sync.RWMutex
	messages map[string]errmsg.ErrorMessage
	locale   string
}

func New(repo errmsg.Repository, locale string) (*ErrorCache, error) {
	msgs, err := repo.FindAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("errcache load: %w", err)
	}
	m := make(map[string]errmsg.ErrorMessage, len(msgs))
	for _, msg := range msgs {
		m[msg.Code] = msg
	}
	return &ErrorCache{messages: m, locale: locale}, nil
}

func (ec *ErrorCache) Error(code string) error {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	if msg, ok := ec.messages[code]; ok {
		text := msg.LocaleEN
		if ec.locale == "th" {
			text = msg.LocaleTH
		}
		return &errmsg.AppError{Code: code, Message: text}
	}
	return &errmsg.AppError{Code: code, Message: code}
}

func (ec *ErrorCache) Message(code string) string {
	if err := ec.Error(code); err != nil {
		return err.Error()
	}
	return code
}
