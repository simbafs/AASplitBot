package bot

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"splitbill/internal/group"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

var ChooseShared = "choose-shared"

func buildKeyboard(r group.Record, g *group.Group) gotgbot.InlineKeyboardMarkup {
	inlineKeyboard := [][]gotgbot.InlineKeyboardButton{{
		{Text: "全選", CallbackData: "all"},
		{Text: "全不選", CallbackData: "none"},
	}}

	for ids := range slices.Chunk(g.IDs(), 3) {
		line := make([]gotgbot.InlineKeyboardButton, 0, 3)
		for _, id := range ids {
			text := ""
			if slices.Contains(r.Shared, id) {
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

	inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
		{Text: "取消", CallbackData: "cancel"},
		{Text: "完成", CallbackData: "done"},
	})

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func (bot *Bot) RecordHandler() ext.Handler {
	return handlers.NewConversation([]ext.Handler{handlers.NewMessage(message.HasPrefix("$"), bot.startRecord)},
		map[string][]ext.Handler{
			ChooseShared: {handlers.NewCallback(callbackquery.All, bot.chooseSharedHandler)},
		}, &handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", bot.cancel)},
			AllowReEntry: true,
		})
}

func (bot *Bot) startRecord(b *gotgbot.Bot, ctx *ext.Context) error {
	send := sender(b, ctx)
	data := bot.storage.Get(id(ctx))
	defer bot.storage.Set(id(ctx), data)

	if ctx.EffectiveChat.Type != gotgbot.ChatTypeGroup {
		return send("目前只支援在群組中使用")
	}

	if data.Group == nil {
		// TODO: 不需要初始化也能使用，用 SelectUser 之類的功能
		return send("目前只支援在已經初始化的群組中使用")
	}

	if _, ok := data.Group.Users[ctx.EffectiveUser.Id]; !ok {
		return send("目前只支援在已經加入分帳的使用者中使用")
	}

	msg := ctx.EffectiveMessage

	amountStr := strings.TrimPrefix(msg.Text, "$")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return send("金額格式錯誤，請輸入以 $ 開頭的金額")
	}

	r := group.Record{
		User:   ctx.EffectiveUser.Id,
		Shared: []int64{},
		Amount: amount,
	}

	replyMarkup := buildKeyboard(r, data.Group)

	msg, err = b.SendMessage(
		ctx.EffectiveChat.Id,
		fmt.Sprintf("%s 出了 %d 元", data.Group.Username(ctx.EffectiveUser.Id), amount),
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

func (bot *Bot) chooseSharedHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	send := sender(b, ctx)
	data := bot.storage.Get(id(ctx))
	defer bot.storage.Set(id(ctx), data)
	defer ctx.CallbackQuery.Answer(b, nil)

	cb := ctx.CallbackQuery
	msgID := cb.Message.GetMessageId()

	r, ok := data.Records[msgID]
	if !ok {
		return send("找不到記錄，請重新開始")
	}

	switch cb.Data {
	case "all":
		r.Shared = slices.Clone(data.Group.IDs())
	case "none":
		r.Shared = []int64{}
	case "cancel":
		delete(data.Records, msgID)
		_, _, err := b.EditMessageText(
			fmt.Sprintf("%s 出了 %d 元（已取消）", data.Group.Username(ctx.EffectiveUser.Id), r.Amount),
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
			fmt.Sprintf("%s 出了 %d 元，%s 要分錢", data.Group.Username(ctx.EffectiveUser.Id), r.Amount, strings.Join(users, "、")),
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
		ReplyMarkup: buildKeyboard(r, data.Group),
	})
	if err != nil {
		return err
	}

	return handlers.NextConversationState(ChooseShared)
}

func (bot *Bot) cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	return handlers.EndConversation()
}
