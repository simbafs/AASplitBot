package splitbot

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"aasplitbot/internal/group"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

var ChooseShared = "choose-shared"

type buildKeyboardOpt struct {
	NoBatch  bool // 是否不使用批次處理
	NoCancel bool // 是否不顯示取消按鈕
}

func buildKeyboard(choosed []int64, g *group.Group, opt *buildKeyboardOpt) gotgbot.InlineKeyboardMarkup {
	inlineKeyboard := [][]gotgbot.InlineKeyboardButton{}
	if opt == nil || !opt.NoBatch {
		inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
			{Text: "全選", CallbackData: "all"},
			{Text: "全不選", CallbackData: "none"},
		})
	}

	for ids := range slices.Chunk(g.IDs(), 3) {
		line := make([]gotgbot.InlineKeyboardButton, 0, 3)
		for _, id := range ids {
			text := ""
			if slices.Contains(choosed, id) {
				text = "✅ " + g.Username(id)
			} else {
				text = g.Username(id)
			}
			line = append(line, gotgbot.InlineKeyboardButton{
				Text:         text,
				CallbackData: strconv.FormatInt(id, 10),
			})
		}
		inlineKeyboard = append(inlineKeyboard, line)
	}

	if opt == nil || !opt.NoCancel {
		inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
			{Text: "取消", CallbackData: "cancel"},
			{Text: "完成", CallbackData: "done"},
		})
	} else {
		inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
			{Text: "完成", CallbackData: "done"},
		})
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func (bot *AASplitBot) Bill() ext.Handler {
	return handlers.NewConversation([]ext.Handler{handlers.NewMessage(message.HasPrefix("$"), bot.billStart)},
		map[string][]ext.Handler{
			ChooseShared: {handlers.NewCallback(callbackquery.All, bot.billChooseShared)},
		}, &handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", bot.billCancel)},
			AllowReEntry: true,
		})
}

func (bot *AASplitBot) billStart(b *gotgbot.Bot, ctx *ext.Context) error {
	send := sender(b, ctx)
	data := bot.storage.Get(id(ctx))
	defer bot.storage.Set(id(ctx), data)

	if data.Group == nil {
		// TODO: 不需要初始化也能使用，用 SelectUser 之類的功能
		return send("目前只支援在已經初始化的群組中使用")
	}

	if _, ok := data.Group.Users[ctx.EffectiveUser.Id]; !ok {
		return send("目前只支援在已經加入分帳的使用者中使用")
	}

	msg := ctx.EffectiveMessage

	r := data.Group.ParseRecord(msg.Text)
	slog.Debug("parse record", "record", r)

	if r.User == 0 {
		r.User = ctx.EffectiveUser.Id
	}

	if len(r.Shared) == 0 {
		r.Shared = keysOfMap(data.Default)
	}

	replyMarkup := buildKeyboard(r.Shared, data.Group, nil)

	msg, err := sendMessage(
		b, ctx,
		fmt.Sprintf("%s 出了 %d 元", data.Group.Username(r.User), r.Amount),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &replyMarkup,
		},
	)
	if err != nil {
		return err
	}

	if data.Records == nil {
		data.Records = make(map[int64]group.Record)
	}

	data.Records[msg.MessageId] = r

	return handlers.NextConversationState(ChooseShared)
}

func (bot *AASplitBot) billChooseShared(b *gotgbot.Bot, ctx *ext.Context) error {
	send := sender(b, ctx)
	data := bot.storage.Get(id(ctx))
	defer bot.storage.Set(id(ctx), data)
	defer ctx.CallbackQuery.Answer(b, nil)

	cb := ctx.CallbackQuery
	msgID := cb.Message.GetMessageId()

	r, ok := data.Records[msgID]
	if !ok {
		if err := send("找不到記錄，請重新開始"); err != nil {
			return err
		}

		return handlers.EndConversation()
	}

	switch cb.Data {
	case "all":
		r.Shared = slices.Clone(data.Group.IDs())
	case "none":
		r.Shared = []int64{}
	case "cancel":
		delete(data.Records, msgID)
		_, _, err := b.EditMessageText(
			fmt.Sprintf("%s 出了 %d 元（已取消）", data.Group.Username(r.User), r.Amount),
			&gotgbot.EditMessageTextOpts{
				ChatId:      ctx.EffectiveChat.Id,
				MessageId:   msgID,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
			},
		)
		if err != nil {
			return fmt.Errorf("clear inline keyboard: %w", err)
		}

		return handlers.EndConversation()
	case "done":
		delete(data.Records, msgID)
		users := make([]string, 0, len(r.Shared))
		for _, id := range r.Shared {
			users = append(users, data.Group.Username(id))
		}
		_, _, err := b.EditMessageText(
			fmt.Sprintf("%s 出了 %d 元，%s 要分錢", data.Group.Username(r.User), r.Amount, strings.Join(users, "、")),
			&gotgbot.EditMessageTextOpts{
				ChatId:      ctx.EffectiveChat.Id,
				MessageId:   msgID,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
			},
		)
		if err != nil {
			return fmt.Errorf("clear inline keyboard: %w", err)
		}

		data.Group.AddRecord(r.User, r.Shared, r.Amount)

		return handlers.EndConversation()
	default:
		id, err := strconv.ParseInt(cb.Data, 10, 64)
		if err != nil {
			return send("無效的使用者 ID: %s", cb.Data)
		}

		slog.Debug("callback query", "record", r, "id", id)
		if slices.Contains(r.Shared, id) {
			r.Shared = slices.DeleteFunc(r.Shared, func(i int64) bool {
				return i == id
			})
		} else {
			r.Shared = append(r.Shared, id)
		}
	}

	data.Records[msgID] = r

	_, _, err := b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   msgID,
		ReplyMarkup: buildKeyboard(r.Shared, data.Group, nil),
	})
	if err != nil {
		return err
	}

	return handlers.NextConversationState(ChooseShared)
}

func (bot *AASplitBot) billCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	return handlers.EndConversation()
}
