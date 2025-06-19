package splitbot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

const ChooseDefault = "choose-default"

func keysOfMap[T any, K comparable](m map[K]T) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (bot *AASplitBot) SetDefault() ext.Handler {
	return handlers.NewConversation([]ext.Handler{bot.startSetDefault()}, map[string][]ext.Handler{
		ChooseDefault: {bot.chooseDefault()},
	}, &handlers.ConversationOpts{
		AllowReEntry: false,
	})
}

func (bot *AASplitBot) startSetDefault() ext.Handler {
	return handlers.NewCommand("setdefault", func(b *gotgbot.Bot, ctx *ext.Context) error {
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)

		if data.Group == nil {
			if err := send("目前還這個群組還沒有資料"); err != nil {
				return err
			}
			return handlers.EndConversation()
		}

		if len(data.Group.IDs()) == 0 {
			if err := send("目前這個群組沒有使用者"); err != nil {
				return err
			}
			return handlers.EndConversation()
		}

		inlineMarkup := buildKeyboard(keysOfMap(data.Default), data.Group, &buildKeyboardOpt{
			NoCancel: true,
		})

		_, err := sendMessage(b, ctx, "請選擇", &gotgbot.SendMessageOpts{
			ReplyMarkup: &inlineMarkup,
		})
		if err != nil {
			return err
		}

		return handlers.NextConversationState(ChooseDefault)
	})
}

func (bot *AASplitBot) chooseDefault() ext.Handler {
	return handlers.NewCallback(callbackquery.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		slog.Info("chooseDefault state")
		send := sender(b, ctx)
		data := bot.storage.Get(id(ctx))
		defer bot.storage.Set(id(ctx), data)
		defer ctx.CallbackQuery.Answer(b, nil)

		if data.Group == nil {
			if err := send("目前還這個群組還沒有資料"); err != nil {
				return err
			}
			return handlers.EndConversation()
		}

		if data.Default == nil {
			data.Default = make(map[int64]struct{})
		}

		switch ctx.CallbackQuery.Data {
		case "all":
			for _, id := range data.Group.IDs() {
				data.Default[id] = struct{}{}
			}
		case "none":
			data.Default = make(map[int64]struct{})
		case "done":
			usernames := make([]string, 0, len(data.Default))
			for id := range data.Default {
				usernames = append(usernames, data.Group.Username(id))
			}
			_, _, err := ctx.EffectiveMessage.EditText(
				b,
				fmt.Sprintf("預設 %s 要分帳", strings.Join(usernames, "、")),
				&gotgbot.EditMessageTextOpts{
					ChatId:      ctx.EffectiveChat.Id,
					MessageId:   ctx.EffectiveMessage.GetMessageId(),
					ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
				},
			)
			if err != nil {
				return fmt.Errorf("update reply keyboard: %w", err)
			}
			return handlers.EndConversation()
		default:
			id, err := strconv.ParseInt(ctx.CallbackQuery.Data, 10, 64)
			if err != nil {
				return send("無效的選擇: %s", ctx.CallbackQuery.Data)
			}

			if _, ok := data.Default[id]; ok {
				delete(data.Default, id)
			} else {
				data.Default[id] = struct{}{}
			}

		}
		_, _, err := ctx.EffectiveMessage.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
			ChatId:    ctx.EffectiveChat.Id,
			MessageId: ctx.EffectiveMessage.GetMessageId(),
			ReplyMarkup: buildKeyboard(keysOfMap(data.Default), data.Group, &buildKeyboardOpt{
				NoCancel: true,
			}),
		})
		if err != nil {
			return fmt.Errorf("update reply keyboard: %w", err)
		}

		return handlers.NextConversationState(ChooseDefault)
	})
}
