package splitbot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// for /listbill
func (bot *AASplitBot) ListBill() ext.Handler {
	return handlers.NewCommand("listbill", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		if data.Group == nil {
			return send("目前還這個群組還沒有資料")
		}

		msg, err := data.Group.RecordsMsg()
		if err != nil {
			return send("無法取得紀錄: %s", err)
		}

		return send(msg)
	})
}
