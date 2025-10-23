package slack

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/slack-go/slack"
)

// Client wraps the slack client
type SlackClient struct {
	*slack.Client
}

// NewClient creates a new Slack client with the given token
func NewClient(token string) (*SlackClient, error) {
	client := slack.New(token)
	return &SlackClient{Client: client}, nil
}

// InitClient initializes the Slack client with the provided token
func InitClient(token string) (*SlackClient, error) {
	return NewClient(token)
}

func (c *SlackClient) Send(slackChannel, msg string) (string, string, error) {
	chId, ts, err := c.PostMessageContext(context.Background(), slackChannel, slack.MsgOptionText(msg, false))
	if err != nil {
		slog.Error("Failed to Post Slack Message", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed to Post Slack Message- %s", err.Error())
	}
	return chId, ts, nil
}
func (c *SlackClient) Update(slackChannel, msgId, newMsg string) (string, string, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: slackChannel,
		Latest:    msgId,
		Oldest:    msgId,
		Inclusive: true,
		Limit:     1,
	}
	history, err := c.GetConversationHistory(params)
	if err != nil {
		slog.Error("Failed to get slack message", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed to get slack message- %s", err.Error())
	}
	if !(len(history.Messages) > 0) {
		slog.Error("Couldnt locate slack message")
		return "", "", fmt.Errorf("couldnt locate slack message")
	}
	existingMsg := history.Messages[0].Text
	if existingMsg == "" {
		slog.Warn("Found Slack message but its empty")
	}
	_, _, err = c.DeleteMessageContext(context.Background(), slackChannel, msgId)
	if err != nil {
		slog.Error("Failed to Delete Slack Message", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed To delete Slack Message err= %s", err)
	}
	existingMsg += newMsg
	chId, ts, err := c.Send(slackChannel, existingMsg)
	if err != nil {
		return "", "", err
	}
	return chId, ts, nil
}

func (c *SlackClient) GetMsgContent(chId, msgId string) (string, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: chId,
		Latest:    msgId,
		Oldest:    msgId,
		Inclusive: true,
		Limit:     1,
	}
	history, err := c.GetConversationHistory(params)
	if err != nil {
		slog.Error("Failed to get slack message", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to get slack message- %s", err.Error())
	}
	if !(len(history.Messages) > 0) {
		slog.Error("Couldnt locate slack message")
		return "", fmt.Errorf("couldnt locate slack message")
	}
	existingMsg := history.Messages[0].Text
	if existingMsg == "" {
		slog.Warn("Found Slack message but its empty")
	}
	return existingMsg, nil
}
func (c *SlackClient) Delete(chId, msgId string) error {
	_, _, err := c.DeleteMessageContext(context.Background(), chId, msgId)
	if err != nil {
		slog.Error("Failed to Delete Slack Message", slog.String("error", err.Error()))
		return fmt.Errorf("failed To delete Slack Message err= %s", err)
	}
	return nil
}
