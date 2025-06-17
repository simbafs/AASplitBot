package usecase // TODO: renmae this package

import (
	"fmt"
	"log/slog"

	"aasplitbot/internal/group"
	"aasplitbot/internal/storage"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func sender(b *gotgbot.Bot, ctx *ext.Context) func(string, ...any) error {
	return func(format string, args ...any) error {
		_, err := ctx.EffectiveChat.SendMessage(b, fmt.Sprintf(format, args...), nil)
		return err
	}
}

func id(ctx *ext.Context) int64 {
	if ctx.EffectiveChat != nil {
		return ctx.EffectiveChat.Id
	}
	if ctx.EffectiveUser != nil {
		return ctx.EffectiveUser.Id
	}
	return 0
}

type ChatData struct {
	Group   *group.Group
	Records map[int64]group.Record
}

type AASplitBot struct {
	storage storage.Storage[ChatData]
}

func New() (*AASplitBot, error) {
	s, err := storage.NewMemory[ChatData]("data.gob")
	if err != nil {
		return nil, fmt.Errorf("creating storage: %w", err)
	}
	return &AASplitBot{
		storage: s,
	}, nil
}

func (bot *AASplitBot) SetCommand(b *gotgbot.Bot, d *ext.Dispatcher) error {
	commands := []gotgbot.BotCommand{
		{Command: "help", Description: "顯示可用指令"},
		{Command: "join", Description: "加入分帳"},
		{Command: "result", Description: "顯示分帳結果"},
		{Command: "listuser", Description: "列出所有使用者"},
		{Command: "listbill", Description: "列出所有分帳紀錄"},
		{Command: "start", Description: "歡迎訊息"},
	}
	d.AddHandler(bot.Bill())
	d.AddHandler(bot.Start())        // /start
	d.AddHandler(bot.Join())         // /join
	d.AddHandler(bot.ListUser())     // listuser
	d.AddHandler(bot.ListBill())     // listbill
	d.AddHandler(bot.Result())       // result
	d.AddHandler(bot.Help(commands)) // help

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
