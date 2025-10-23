package telegram

import (
	"fmt"
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Client wraps the telegram bot client
type TelegramClient struct {
	*tgbotapi.BotAPI
}

// NewClient creates a new Telegram client with the given token
func NewClient(token string) (*TelegramClient, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramClient{BotAPI: bot}, nil
}

// InitClient initializes the Telegram client with the provided token
func InitClient(token string) (*TelegramClient, error) {
	return NewClient(token)
}

func (c *TelegramClient) Send(telegramChatId, msg string) (string, error) {
	intTelegramChatId, err := strconv.ParseInt(telegramChatId, 10, 64)
	if err != nil {
		slog.Error("Failed to parse telegramChatId to int64", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to parse telegramChatId to int64- %s", err.Error())
	}
	msgConfig := tgbotapi.NewMessage(intTelegramChatId, msg)
	msgConfig.ParseMode = "Markdown"
	tgMsg, err := c.BotAPI.Send(msgConfig)
	if err != nil {
		return "", fmt.Errorf("failed To send Telegram Message err= %s", err)
	}
	messageIdstr := strconv.Itoa(tgMsg.MessageID)
	return messageIdstr, nil
}
