package usecase

import (
	"splitbill/internal/group"

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

		if ctx.EffectiveChat.Type != gotgbot.ChatTypeGroup {
			return send("目前只支援在群組中使用")
		}

		if data.Group == nil {
			data.Group = group.New()
		}

		data.Group.AddUser(ctx.EffectiveUser.Id, ctx.EffectiveSender.Username())

		return send("已將 %s 加入分帳", ctx.EffectiveUser.Username)
	})
}
