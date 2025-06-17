package usecase

import (
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *AASplitBot) SetCommand(b *gotgbot.Bot, d *ext.Dispatcher) error {
	d.AddHandler(bot.Bill())
	d.AddHandler(bot.Start())    // /start
	d.AddHandler(bot.Join())     // /join
	d.AddHandler(bot.ListUser()) // listuser
	d.AddHandler(bot.ListBill()) // listbill
	d.AddHandler(bot.Result())   // result

	commands := []gotgbot.BotCommand{
		{Command: "start", Description: "歡迎訊息"},
		{Command: "join", Description: "加入分帳"},
		{Command: "listuser", Description: "列出所有使用者"},
		{Command: "listbill", Description: "列出所有分帳"},
		{Command: "result", Description: "顯示分帳結果"},
	}

	ok, err := b.SetMyCommands(commands, &gotgbot.SetMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeAllGroupChats{},
	})
	if err != nil {
		return fmt.Errorf("setting commands: %w", err)
	}
	if !ok {
		slog.Error("failed to set commands")
	} else {
		slog.Info("commands set successfully")
	}

	return nil
}
