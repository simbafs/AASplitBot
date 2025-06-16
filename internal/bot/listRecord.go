package bot

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func (bot *Bot) ListRecordHandler() ext.Handler {
	return handlers.NewCommand("listRecord", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		if data.Group == nil {
			return send("目前還這個群組還沒有資料")
		}

		if len(data.Group.Bills) == 0 {
			return send("目前沒有任何分帳紀錄")
		}

		var records []string
		for _, r := range data.Group.Bills {
			usernames := make([]string, 0, len(r.Shared))
			for _, id := range r.Shared {
				usernames = append(usernames, data.Group.Username(id))
			}

			records = append(records, fmt.Sprintf("%s 代墊了 %d 元，%s 要付錢", data.Group.Username(r.User), r.Amount, strings.Join(usernames, "、")))
		}

		return send("目前的分帳紀錄有：\n%s", strings.Join(records, "\n"))
	})
}
