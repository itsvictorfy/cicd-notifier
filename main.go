package main

import (
	"cicd-notifier/pkg/slack"
	"cicd-notifier/pkg/telegram"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	ParsedInputs ActionInputs
)

func init() {
	initDev()

	// Read inputs and initialize the struct
	inputs := readInputs()
	ParsedInputs = ActionInputs{
		Action:       strings.ToLower(inputs["action"]),
		MsgID:        inputs["msgid"],
		Message:      inputs["message"],
		ApiKey:       inputs["api_key"],
		ChannelId:    inputs["channel_id"],
		ImageTag:     inputs["image_tag"],
		CommitSha:    inputs["commit_sha"],
		Branch:       inputs["branch"],
		Author:       inputs["author"],
		CommitTime:   inputs["commit_time"],
		CommitMsg:    inputs["commit_msg"],
		WorkflowName: inputs["workflow_name"],
	}

	// Parse Channel
	if channelStr := inputs["channel"]; channelStr != "" {
		ParsedInputs.Channel = strings.TrimSpace(channelStr)
	}

	// Parse AddCommitInfo bool
	if addCommitStr := inputs["add_commit_info"]; addCommitStr != "" {
		if parsed, err := strconv.ParseBool(addCommitStr); err == nil {
			ParsedInputs.AddCommitInfo = parsed
		}
	}
}
func initDev() {
	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file doesn't exist")
	}
}
func main() {
	// Validate inputs
	validateInputs()

	var msg string
	if ParsedInputs.Action == "send" {
		if ParsedInputs.AddCommitInfo {
			msg = templateCommitInfo()
		}
		msg += fmt.Sprintf("* - %s*- %s \n", ParsedInputs.Message, time.Now())
		switch strings.ToLower(ParsedInputs.Channel) {
		case "slack":
			c, err := slack.InitClient(ParsedInputs.ApiKey)
			if err != nil {
				slog.Error("Failed to initialize slack client", slog.String("error", err.Error()))
				os.Exit(1)
			}
			chId, msgId, err := c.Send(ParsedInputs.ChannelId, msg)
			if err != nil {
				slog.Error("Failed To Post Slack message", slog.String("error", err.Error()))
				os.Exit(1)
			}
			addOutput("message_id", msgId)
			addOutput("channel_id", chId)
		case "telegram":
			c, err := telegram.InitClient(ParsedInputs.ApiKey)
			if err != nil {
				slog.Error("Failed to initialize telegram client", slog.String("error", err.Error()))
				os.Exit(1)
			}
			msgId, err := c.Send(ParsedInputs.ChannelId, msg)
			if err != nil {
				slog.Error("Failed To Post Telegram message", slog.String("error", err.Error()))
				os.Exit(1)
			}
			addOutput("message_id", msgId)
		}
	}
	if ParsedInputs.Action == "update" {
		switch strings.ToLower(ParsedInputs.Channel) {
		case "slack":
			c, err := slack.InitClient(ParsedInputs.ApiKey)
			if err != nil {
				slog.Error("Failed to initialize slack client", slog.String("error", err.Error()))
				os.Exit(1)
			}
			msg, err = c.GetMsgContent(ParsedInputs.ChannelId, ParsedInputs.MsgID)
			if err != nil {
				slog.Error("Failed get slack message Content", slog.String("error", err.Error()))
				os.Exit(1)
			}
			msg += fmt.Sprintf(" - *%s*- %s \n", ParsedInputs.Message, time.Now())
			err = c.Delete(ParsedInputs.ChannelId, ParsedInputs.MsgID)
			if err != nil {
				slog.Error("Failed to delete slack message", slog.String("error", err.Error()))
				os.Exit(1)
			}
			chId, msgId, err := c.Send(ParsedInputs.ChannelId, msg)
			if err != nil {
				slog.Error("Failed To send Slack message", slog.String("error", err.Error()))
				os.Exit(1)
			}
			addOutput("message_id", msgId)
			addOutput("channel_id", chId)
		case "telegram":
			slog.Error("Update message is not supported in Telegram")
			os.Exit(1)
		}
	}
	// Set outputs for GitHub Actions
	setOutputs()
	os.Exit(0)
}
