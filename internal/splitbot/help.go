package splitbot

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func (bot *AASplitBot) Help(cmds []gotgbot.BotCommand) ext.Handler {
	return handlers.NewCommand("help", func(b *gotgbot.Bot, ctx *ext.Context) error {
		msg := strings.Builder{}

		msg.WriteString("可用指令:\n")
		for _, cmd := range cmds {
			msg.WriteString("/" + cmd.Command + " - " + cmd.Description + "\n")
		}

		msg.WriteString("\n輸入以錢號（`$`）開頭、接續著一個數字的訊息可以新增分帳紀錄，例如 `$10`。\n")

		_, err := ctx.EffectiveChat.SendMessage(b, msg.String(), nil)
		return err
	})
}
