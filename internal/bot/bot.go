package bot // TODO: renmae this package

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"splitbill/internal/group"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type Bot struct {
	groups map[int64]*group.Group
}

func New() *Bot {
	return &Bot{
		groups: make(map[int64]*group.Group),
	}
}

func (b *Bot) StartHandler() ext.Handler {
	return handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "hello", nil)
		return err
	})
}

func (bot *Bot) RecordHandler() ext.Handler {
	return handlers.NewMessage(func(msg *gotgbot.Message) bool {
		if msg.Text[0] != '$' {
			return false
		}
		if _, ok := bot.groups[msg.Chat.Id]; !ok {
			return true
		}
		return true
	}, func(b *gotgbot.Bot, ctx *ext.Context) error {
		logger := slog.With("chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)
		log := func(msg string) error {
			logger.Info(msg)
			_, err := b.SendMessage(ctx.EffectiveChat.Id, msg, nil)
			return err
		}

		if ctx.EffectiveChat.Type != gotgbot.ChatTypeGroup {
			return log("目前只支援在群組中使用")
		}

		_, ok := bot.groups[ctx.EffectiveChat.Id]
		if !ok {
			return log("目前只支援在已經初始化的群組中使用")
		}

		msg := ctx.EffectiveMessage
		logger.Info("msg", "text", msg.Text)

		amountStr := strings.TrimPrefix(msg.Text, "$")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return log("金額格式錯誤，請輸入以 $ 開頭的金額")
		}

		_, err = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("%d 出了 %d 元", ctx.EffectiveUser.Id, amount), nil)
		return err
	})
}

func (bot *Bot) Inithandler() ext.Handler {
	return handlers.NewCommand("init", func(b *gotgbot.Bot, ctx *ext.Context) error {
		logger := slog.With("chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)

		if ctx.EffectiveChat.Type != gotgbot.ChatTypeGroup {
			_, err := b.SendMessage(ctx.EffectiveChat.Id, "目前只支援在群組中使用", nil)
			return err
		}

		if _, ok := bot.groups[ctx.EffectiveChat.Id]; ok {
			_, err := b.SendMessage(ctx.EffectiveChat.Id, "群組已經初始化過了", nil)
			return err
		}

		bot.groups[ctx.EffectiveChat.Id] = group.New()

		logger.Info("group initialized")
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "群組已經初始化完成", nil)
		return err
	})
}
