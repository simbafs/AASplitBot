package splitbot

import (
	"aasplitbot/internal/group"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// for command /join
func (bot *AASplitBot) Join() ext.Handler {
	return handlers.NewCommand("join", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		if data.Group == nil {
			data.Group = group.New()
		}

		data.Group.AddUser(ctx.EffectiveUser.Id, ctx.EffectiveSender.Username())

		return send("已將 %s 加入分帳", ctx.EffectiveUser.Username)
	})
}
