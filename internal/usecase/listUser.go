package usecase

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// for /listUser
func (bot *AASplitBot) ListUser() ext.Handler {
	return handlers.NewCommand("listUser", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		if data.Group == nil {
			return send("目前還這個群組還沒有資料")
		}

		return send("%s", strings.Join(data.Group.Usernames(), "、"))
	})
}
