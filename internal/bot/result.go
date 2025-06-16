package bot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func (bot *Bot) ResultHandler() ext.Handler {
	return handlers.NewCommand("result", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))

		if data.Group == nil {
			return send("還沒有紀錄")
		}

		result, err := data.Group.ResultMsg()
		if err != nil {
			return send("無法取得結果: %s", err)
		}

		send(result)

		return nil
	})
}
