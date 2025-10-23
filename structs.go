package main

// ActionInputs represents the input parameters for the notification action
type ActionInputs struct {
	Action        string // Required: Send/Update
	Channel       string // Required: Channel (Telegram/Slack)
	Message       string // Required: Message to send
	ApiKey        string // Required: API key
	ChannelId     string // Required: channel/chat used in slack/telegram
	MsgID         string // Optional: ID of the message to update
	AddCommitInfo bool   // Optional: Whether to add commit info
	ImageTag      string // Optional: Docker image tag
	CommitSha     string // Optional: Commit SHA
	Branch        string // Optional: Branch name
	Author        string // Optional: Commit author
	CommitTime    string // Optional: Commit time
	CommitMsg     string // Optional: Commit Message
	WorkflowName  string // Optional: WorkflowName
}
