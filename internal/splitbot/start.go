package splitbot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// for command /start
func (bot *AASplitBot) Start() ext.Handler {
	return handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		msg := `歡迎使用分帳機器人
開始使用之前每個人都要用指令 /join 加入分帳
`
		_, err := b.SendMessage(ctx.EffectiveChat.Id, msg, nil)
		return err
	})
}
