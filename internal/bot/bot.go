package bot // TODO: renmae this package

import (
	"fmt"

	"splitbill/internal/group"
	"splitbill/internal/storage"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func sender(b *gotgbot.Bot, ctx *ext.Context) func(string, ...any) error {
	return func(format string, args ...any) error {
		_, err := ctx.EffectiveChat.SendMessage(b, fmt.Sprintf(format, args...), nil)
		return err
	}
}

func id(ctx *ext.Context) int64 {
	if ctx.EffectiveChat != nil {
		return ctx.EffectiveChat.Id
	}
	if ctx.EffectiveUser != nil {
		return ctx.EffectiveUser.Id
	}
	return 0
}

type ChatData struct {
	Group   *group.Group
	Records map[int64]group.Record
}

type Bot struct {
	storage storage.Storage[ChatData]
}

func New() (*Bot, error) {
	s, err := storage.NewMemory[ChatData]("data.gob")
	if err != nil {
		return nil, fmt.Errorf("creating storage: %w", err)
	}
	return &Bot{
		storage: s,
	}, nil
}
