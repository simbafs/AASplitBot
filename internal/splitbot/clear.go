package splitbot

import (
	"fmt"
	"log/slog"

	"aasplitbot/internal/group"
	"aasplitbot/internal/workerpool"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// for /clear
func (bot *AASplitBot) Clear() ext.Handler {
	return handlers.NewCommand("clear", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		data.Group.Clear()

		pool := workerpool.New()
		for msgID, r := range data.Records {
			pool.Do(func() error {
				_, _, err := b.EditMessageText(
					fmt.Sprintf("%s 出了 %d 元（已取消）", data.Group.Username(r.User), r.Amount),
					&gotgbot.EditMessageTextOpts{
						ChatId:      ctx.EffectiveChat.Id,
						MessageId:   msgID,
						ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
					},
				)
				if err != nil {
					slog.Error("clear inline keyboard", "error", err)
				}
				return nil
			})
		}
		pool.Wait()
		data.Records = make(map[int64]group.Record)

		if err := send("已清除所有分帳紀錄"); err != nil {
			return fmt.Errorf("send clear message: %w", err)
		}

		return nil
	})
}
