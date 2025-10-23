package main

// ActionInputs represents the input parameters for the notification action
type ActionInputs struct {
	Action          string // Required: Send/Update
	SlackMsgID      string // Optional: ID of Slack message for update
	TelegramMsgID   string // Optional: ID of Telegram message for update
	Channel         string // Required: Channel (Telegram/Slack)
	Message         string // Required: Message to send
	SlackApiKey     string // Optional: Slack API key
	SlackChannel    string // Optional: Slack channel
	TelegramApiKey  string // Optional: Telegram API key
	TelegramChannel string // Optional: Telegram channel
	AddCommitInfo   bool   // Optional: Whether to add commit info
	ImageTag        string // Optional: Docker image tag
	CommitSha       string // Optional: Commit SHA
	Branch          string // Optional: Branch name
	Author          string // Optional: Commit author
	CommitTime      string // Optional: Commit time
	CommitMsg       string // Optional: Commit Message
	WorkflowName    string // Optional: WorkflowName
}
