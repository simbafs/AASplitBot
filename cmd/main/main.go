package main

import (
	"fmt"
	"log/slog"
	"time"

	"aasplitbot/internal/config"
	"aasplitbot/internal/usecase"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func run() error {
	cfg := config.NewConfigWithEnv()

	b, err := gotgbot.NewBot(cfg.Token, nil)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	me, err := b.GetMe(nil)
	if err != nil {
		return fmt.Errorf("getting bot info: %w", err)
	}
	slog.Info("authorized", "username", me.Username, "id", me.Id)

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			slog.Error("bot dispatcher", "error", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	cases, err := usecase.New()
	if err != nil {
		return fmt.Errorf("creating my bot: %w", err)
	}

	if err := cases.SetCommand(b, dispatcher); err != nil {
		return fmt.Errorf("setting commands: %w", err)
	}

	updater := ext.NewUpdater(dispatcher, nil)

	if err := updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	}); err != nil {
		return fmt.Errorf("bot polling update: %w", err)
	}

	slog.Info("bot start")

	updater.Idle()

	return nil
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err := run(); err != nil {
		panic(err)
	}
}
