package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHasAllRequiredInputs(t *testing.T) {
	tests := []struct {
		name     string
		channel  string
		inputs   ActionInputs
		expected bool
	}{
		{
			name:    "Slack send with all required inputs",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "send",
				Message:      "test message",
				SlackApiKey:  "test_key",
				SlackChannel: "test_channel",
			},
			expected: true,
		},
		{
			name:    "Slack send missing message",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "send",
				SlackApiKey:  "test_key",
				SlackChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Slack send missing api key",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "send",
				Message:      "test message",
				SlackChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Slack send missing channel",
			channel: "slack",
			inputs: ActionInputs{
				Action:      "send",
				Message:     "test message",
				SlackApiKey: "test_key",
			},
			expected: false,
		},
		{
			name:    "Slack update with all required inputs",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "update",
				Message:      "test message",
				SlackApiKey:  "test_key",
				SlackChannel: "test_channel",
				SlackMsgID:   "123",
			},
			expected: true,
		},
		{
			name:    "Slack update missing msgid",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "update",
				Message:      "test message",
				SlackApiKey:  "test_key",
				SlackChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Telegram send with all required inputs",
			channel: "telegram",
			inputs: ActionInputs{
				Action:          "send",
				Message:         "test message",
				TelegramApiKey:  "test_key",
				TelegramChannel: "test_channel",
			},
			expected: true,
		},
		{
			name:    "Telegram send missing message",
			channel: "telegram",
			inputs: ActionInputs{
				Action:          "send",
				TelegramApiKey:  "test_key",
				TelegramChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Telegram send missing api key",
			channel: "telegram",
			inputs: ActionInputs{
				Action:          "send",
				Message:         "test message",
				TelegramChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Telegram send missing channel",
			channel: "telegram",
			inputs: ActionInputs{
				Action:         "send",
				Message:        "test message",
				TelegramApiKey: "test_key",
			},
			expected: false,
		},
		{
			name:    "Telegram update with all required inputs",
			channel: "telegram",
			inputs: ActionInputs{
				Action:          "update",
				Message:         "test message",
				TelegramApiKey:  "test_key",
				TelegramChannel: "test_channel",
				TelegramMsgID:   "123",
			},
			expected: true,
		},
		{
			name:    "Telegram update missing msgid",
			channel: "telegram",
			inputs: ActionInputs{
				Action:          "update",
				Message:         "test message",
				TelegramApiKey:  "test_key",
				TelegramChannel: "test_channel",
			},
			expected: false,
		},
		{
			name:    "Missing message fails",
			channel: "slack",
			inputs: ActionInputs{
				Action:       "send",
				SlackApiKey:  "test_key",
				SlackChannel: "test_channel",
				// Message is empty
			},
			expected: false,
		},
		{
			name:     "Unsupported channel",
			channel:  "unsupported",
			inputs:   ActionInputs{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global ParsedInputs
			ParsedInputs = tt.inputs

			err := hasAllRequiredInputs(tt.channel)
			if (err == nil) != tt.expected {
				t.Errorf("hasAllRequiredInputs() error = %v, expected no error = %v", err, tt.expected)
			}
		})
	}
}

func TestReadInputs(t *testing.T) {
	// Set up environment variables
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "test message")
	os.Setenv("INPUT_SLACK_API_KEY", "slack_key")
	os.Setenv("INPUT_SLACK_CHANNEL", "slack_chan")
	os.Setenv("INPUT_TELEGRAM_API_KEY", "telegram_key")
	os.Setenv("INPUT_TELEGRAM_CHANNEL", "telegram_chan")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "true")
	os.Setenv("INPUT_COMMIT_SHA", "abc123")
	os.Setenv("INPUT_BRANCH", "main")
	os.Setenv("INPUT_AUTHOR", "testuser")
	os.Setenv("INPUT_COMMIT_TIME", "2023-01-01")
	os.Setenv("INPUT_COMMIT_MSG", "test commit")
	os.Setenv("INPUT_WORKFLOW_NAME", "test workflow")
	os.Setenv("INPUT_IMAGE_TAG", "v1.0")
	os.Setenv("INPUT_SLACK_MSGID", "msg123")
	os.Setenv("INPUT_TELEGRAM_MSGID", "tmsg123")
	defer func() {
		// Clean up
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_SLACK_API_KEY")
		os.Unsetenv("INPUT_SLACK_CHANNEL")
		os.Unsetenv("INPUT_TELEGRAM_API_KEY")
		os.Unsetenv("INPUT_TELEGRAM_CHANNEL")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
		os.Unsetenv("INPUT_COMMIT_SHA")
		os.Unsetenv("INPUT_BRANCH")
		os.Unsetenv("INPUT_AUTHOR")
		os.Unsetenv("INPUT_COMMIT_TIME")
		os.Unsetenv("INPUT_COMMIT_MSG")
		os.Unsetenv("INPUT_WORKFLOW_NAME")
		os.Unsetenv("INPUT_IMAGE_TAG")
		os.Unsetenv("INPUT_SLACK_MSGID")
		os.Unsetenv("INPUT_TELEGRAM_MSGID")
	}()

	inputs := readInputs()

	expected := map[string]string{
		"action":           "send",
		"message":          "test message",
		"slack_api_key":    "slack_key",
		"slack_channel":    "slack_chan",
		"telegram_api_key": "telegram_key",
		"telegram_channel": "telegram_chan",
		"channel":          "slack",
		"add_commit_info":  "true",
		"commit_sha":       "abc123",
		"branch":           "main",
		"author":           "testuser",
		"commit_time":      "2023-01-01",
		"commit_msg":       "test commit",
		"workflow_name":    "test workflow",
		"image_tag":        "v1.0",
		"slack_msgid":      "msg123",
		"telegram_msgid":   "tmsg123",
	}

	for key, expectedValue := range expected {
		if actualValue, ok := inputs[key]; !ok {
			t.Errorf("readInputs() missing key %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("readInputs() key %s = %s, expected %s", key, actualValue, expectedValue)
		}
	}
}

func TestAddOutput(t *testing.T) {
	// Reset global outputs
	outputs = make(map[string]string)

	addOutput("key1", "value1")
	addOutput("key2", "value2")

	if outputs["key1"] != "value1" {
		t.Errorf("addOutput() key1 = %s, expected value1", outputs["key1"])
	}
	if outputs["key2"] != "value2" {
		t.Errorf("addOutput() key2 = %s, expected value2", outputs["key2"])
	}
}

func TestTemplateCommitInfo(t *testing.T) {
	ParsedInputs = ActionInputs{
		AddCommitInfo: true,
		CommitSha:     "abc123",
		Branch:        "main",
		WorkflowName:  "CI",
		CommitMsg:     "fix bug",
		Author:        "user",
		ImageTag:      "v1.0",
		CommitTime:    "2023-01-01",
	}

	result := templateCommitInfo()

	expected := "📦 *Github Workflow*\n\n📌 *Commit:* `abc123`\n🔖 *Branch:* `main`\n🛠️ *Workflow:* `CI`\n📝 *Message:* fix bug\n👤 *Author:* user\n🐳 *Image Tag:* v1.0\n🕗 *Commit Time:* 2023-01-01\n"
	if result != expected {
		t.Errorf("templateCommitInfo() = %q, expected %q", result, expected)
	}
}

func TestReadInputsTemplateCommitInfoAndAddMessage(t *testing.T) {
	// Set up environment variables
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Deployment completed successfully")
	os.Setenv("INPUT_SLACK_API_KEY", "slack_key")
	os.Setenv("INPUT_SLACK_CHANNEL", "slack_chan")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "true")
	os.Setenv("INPUT_COMMIT_SHA", "def456")
	os.Setenv("INPUT_BRANCH", "develop")
	os.Setenv("INPUT_AUTHOR", "developer")
	os.Setenv("INPUT_COMMIT_TIME", "2023-01-02")
	os.Setenv("INPUT_COMMIT_MSG", "update feature")
	os.Setenv("INPUT_WORKFLOW_NAME", "Deploy")
	os.Setenv("INPUT_IMAGE_TAG", "v2.0")
	defer func() {
		// Clean up
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_SLACK_API_KEY")
		os.Unsetenv("INPUT_SLACK_CHANNEL")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
		os.Unsetenv("INPUT_COMMIT_SHA")
		os.Unsetenv("INPUT_BRANCH")
		os.Unsetenv("INPUT_AUTHOR")
		os.Unsetenv("INPUT_COMMIT_TIME")
		os.Unsetenv("INPUT_COMMIT_MSG")
		os.Unsetenv("INPUT_WORKFLOW_NAME")
		os.Unsetenv("INPUT_IMAGE_TAG")
	}()

	// Read inputs
	inputs := readInputs()

	// Simulate init logic
	ParsedInputs = ActionInputs{
		Action:       inputs["action"],
		Message:      inputs["message"],
		SlackApiKey:  inputs["slack_api_key"],
		SlackChannel: inputs["slack_channel"],
		Channel:      inputs["channel"],
	}

	// Parse AddCommitInfo
	if addCommitStr := inputs["add_commit_info"]; addCommitStr != "" {
		if parsed, err := strconv.ParseBool(addCommitStr); err == nil {
			ParsedInputs.AddCommitInfo = parsed
		}
	}

	// Set other fields
	ParsedInputs.CommitSha = inputs["commit_sha"]
	ParsedInputs.Branch = inputs["branch"]
	ParsedInputs.Author = inputs["author"]
	ParsedInputs.CommitTime = inputs["commit_time"]
	ParsedInputs.CommitMsg = inputs["commit_msg"]
	ParsedInputs.WorkflowName = inputs["workflow_name"]
	ParsedInputs.ImageTag = inputs["image_tag"]

	// Template commit info if enabled
	var commitMsg string
	if ParsedInputs.AddCommitInfo {
		commitMsg = templateCommitInfo()
	} else {
		commitMsg = "📦 *Github Workflow*\n\n"
	}

	// Add the message field at the end
	finalMessage := commitMsg + ParsedInputs.Message

	expected := "📦 *Github Workflow*\n\n📌 *Commit:* `def456`\n🔖 *Branch:* `develop`\n🛠️ *Workflow:* `Deploy`\n📝 *Message:* update feature\n👤 *Author:* developer\n🐳 *Image Tag:* v2.0\n🕗 *Commit Time:* 2023-01-02\nDeployment completed successfully"
	if finalMessage != expected {
		t.Errorf("Final message = %q, expected %q", finalMessage, expected)
	}
}

func TestSlackSendFlow(t *testing.T) {
	// Set up environment variables for Slack send
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Test message")
	os.Setenv("INPUT_SLACK_API_KEY", "test_slack_key")
	os.Setenv("INPUT_SLACK_CHANNEL", "#test")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "false")
	defer func() {
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_SLACK_API_KEY")
		os.Unsetenv("INPUT_SLACK_CHANNEL")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
	}()

	// Simulate init logic
	inputs := readInputs()
	ParsedInputs = ActionInputs{
		Action:          strings.ToLower(inputs["action"]),
		SlackMsgID:      inputs["slack_msgid"],
		TelegramMsgID:   inputs["telegram_msgid"],
		Message:         inputs["message"],
		SlackApiKey:     inputs["slack_api_key"],
		SlackChannel:    inputs["slack_channel"],
		TelegramApiKey:  inputs["telegram_api_key"],
		TelegramChannel: inputs["telegram_channel"],
		ImageTag:        inputs["image_tag"],
		CommitSha:       inputs["commit_sha"],
		Branch:          inputs["branch"],
		Author:          inputs["author"],
		CommitTime:      inputs["commit_time"],
		CommitMsg:       inputs["commit_msg"],
		WorkflowName:    inputs["workflow_name"],
	}
	if channelStr := inputs["channel"]; channelStr != "" {
		ParsedInputs.Channel = strings.TrimSpace(channelStr)
	}
	if addCommitStr := inputs["add_commit_info"]; addCommitStr != "" {
		if parsed, err := strconv.ParseBool(addCommitStr); err == nil {
			ParsedInputs.AddCommitInfo = parsed
		}
	}

	// Check ParsedInputs
	if ParsedInputs.Action != "send" {
		t.Errorf("ParsedInputs.Action = %s, expected send", ParsedInputs.Action)
	}
	if ParsedInputs.Message != "Test message" {
		t.Errorf("ParsedInputs.Message = %s, expected Test message", ParsedInputs.Message)
	}
	if ParsedInputs.Channel != "slack" {
		t.Errorf("ParsedInputs.Channel = %s, expected slack", ParsedInputs.Channel)
	}
	if ParsedInputs.SlackApiKey != "test_slack_key" {
		t.Errorf("ParsedInputs.SlackApiKey = %s, expected test_slack_key", ParsedInputs.SlackApiKey)
	}
	if ParsedInputs.SlackChannel != "#test" {
		t.Errorf("ParsedInputs.SlackChannel = %s, expected #test", ParsedInputs.SlackChannel)
	}
	if ParsedInputs.AddCommitInfo != false {
		t.Errorf("ParsedInputs.AddCommitInfo = %v, expected false", ParsedInputs.AddCommitInfo)
	}

	// Check validation passes
	err := hasAllRequiredInputs("slack")
	if err != nil {
		t.Errorf("hasAllRequiredInputs failed: %v", err)
	}

	// Simulate the send flow (without API call)
	var commitMsg string
	if ParsedInputs.AddCommitInfo {
		commitMsg = templateCommitInfo()
	} else {
		commitMsg = "📦 *Github Workflow*\n\n"
	}
	finalMessage := commitMsg + ParsedInputs.Message

	expectedMessage := "📦 *Github Workflow*\n\nTest message"
	if finalMessage != expectedMessage {
		t.Errorf("Final message = %q, expected %q", finalMessage, expectedMessage)
	}
}

func TestSlackUpdateFlow(t *testing.T) {
	// Set up environment variables for Slack update
	os.Setenv("INPUT_ACTION", "update")
	os.Setenv("INPUT_MESSAGE", "Updated message")
	os.Setenv("INPUT_SLACK_API_KEY", "test_slack_key")
	os.Setenv("INPUT_SLACK_CHANNEL", "#test")
	os.Setenv("INPUT_SLACK_MSGID", "123456")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "true")
	os.Setenv("INPUT_COMMIT_SHA", "abc123")
	os.Setenv("INPUT_BRANCH", "main")
	os.Setenv("INPUT_AUTHOR", "user")
	os.Setenv("INPUT_COMMIT_TIME", "2023-01-01")
	os.Setenv("INPUT_COMMIT_MSG", "fix")
	os.Setenv("INPUT_WORKFLOW_NAME", "CI")
	os.Setenv("INPUT_IMAGE_TAG", "v1.0")
	defer func() {
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_SLACK_API_KEY")
		os.Unsetenv("INPUT_SLACK_CHANNEL")
		os.Unsetenv("INPUT_SLACK_MSGID")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
		os.Unsetenv("INPUT_COMMIT_SHA")
		os.Unsetenv("INPUT_BRANCH")
		os.Unsetenv("INPUT_AUTHOR")
		os.Unsetenv("INPUT_COMMIT_TIME")
		os.Unsetenv("INPUT_COMMIT_MSG")
		os.Unsetenv("INPUT_WORKFLOW_NAME")
		os.Unsetenv("INPUT_IMAGE_TAG")
	}()

	// Simulate init logic
	inputs := readInputs()
	ParsedInputs = ActionInputs{
		Action:          strings.ToLower(inputs["action"]),
		SlackMsgID:      inputs["slack_msgid"],
		TelegramMsgID:   inputs["telegram_msgid"],
		Message:         inputs["message"],
		SlackApiKey:     inputs["slack_api_key"],
		SlackChannel:    inputs["slack_channel"],
		TelegramApiKey:  inputs["telegram_api_key"],
		TelegramChannel: inputs["telegram_channel"],
		ImageTag:        inputs["image_tag"],
		CommitSha:       inputs["commit_sha"],
		Branch:          inputs["branch"],
		Author:          inputs["author"],
		CommitTime:      inputs["commit_time"],
		CommitMsg:       inputs["commit_msg"],
		WorkflowName:    inputs["workflow_name"],
	}
	if channelStr := inputs["channel"]; channelStr != "" {
		ParsedInputs.Channel = strings.TrimSpace(channelStr)
	}
	if addCommitStr := inputs["add_commit_info"]; addCommitStr != "" {
		if parsed, err := strconv.ParseBool(addCommitStr); err == nil {
			ParsedInputs.AddCommitInfo = parsed
		}
	}

	// Check ParsedInputs
	if ParsedInputs.Action != "update" {
		t.Errorf("ParsedInputs.Action = %s, expected update", ParsedInputs.Action)
	}
	if ParsedInputs.SlackMsgID != "123456" {
		t.Errorf("ParsedInputs.SlackMsgID = %s, expected 123456", ParsedInputs.SlackMsgID)
	}
	if !ParsedInputs.AddCommitInfo {
		t.Errorf("ParsedInputs.AddCommitInfo = %v, expected true", ParsedInputs.AddCommitInfo)
	}

	// Check validation passes
	err := hasAllRequiredInputs("slack")
	if err != nil {
		t.Errorf("hasAllRequiredInputs failed: %v", err)
	}

	// Simulate the update flow
	var commitMsg string
	if ParsedInputs.AddCommitInfo {
		commitMsg = templateCommitInfo()
	} else {
		commitMsg = "📦 *Github Workflow*\n\n"
	}
	finalMessage := commitMsg + ParsedInputs.Message

	expectedMessage := "📦 *Github Workflow*\n\n📌 *Commit:* `abc123`\n🔖 *Branch:* `main`\n🛠️ *Workflow:* `CI`\n📝 *Message:* fix\n👤 *Author:* user\n🐳 *Image Tag:* v1.0\n🕗 *Commit Time:* 2023-01-01\nUpdated message"
	if finalMessage != expectedMessage {
		t.Errorf("Final message = %q, expected %q", finalMessage, expectedMessage)
	}
}

func TestTelegramSendFlow(t *testing.T) {
	// Set up environment variables for Telegram send
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Telegram test")
	os.Setenv("INPUT_TELEGRAM_API_KEY", "test_telegram_key")
	os.Setenv("INPUT_TELEGRAM_CHANNEL", "@testchannel")
	os.Setenv("INPUT_CHANNEL", "telegram")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "false")
	defer func() {
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_TELEGRAM_API_KEY")
		os.Unsetenv("INPUT_TELEGRAM_CHANNEL")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
	}()

	// Simulate init logic
	inputs := readInputs()
	ParsedInputs = ActionInputs{
		Action:          strings.ToLower(inputs["action"]),
		SlackMsgID:      inputs["slack_msgid"],
		TelegramMsgID:   inputs["telegram_msgid"],
		Message:         inputs["message"],
		SlackApiKey:     inputs["slack_api_key"],
		SlackChannel:    inputs["slack_channel"],
		TelegramApiKey:  inputs["telegram_api_key"],
		TelegramChannel: inputs["telegram_channel"],
		ImageTag:        inputs["image_tag"],
		CommitSha:       inputs["commit_sha"],
		Branch:          inputs["branch"],
		Author:          inputs["author"],
		CommitTime:      inputs["commit_time"],
		CommitMsg:       inputs["commit_msg"],
		WorkflowName:    inputs["workflow_name"],
	}
	if channelStr := inputs["channel"]; channelStr != "" {
		ParsedInputs.Channel = strings.TrimSpace(channelStr)
	}
	if addCommitStr := inputs["add_commit_info"]; addCommitStr != "" {
		if parsed, err := strconv.ParseBool(addCommitStr); err == nil {
			ParsedInputs.AddCommitInfo = parsed
		}
	}

	// Check ParsedInputs
	if ParsedInputs.Channel != "telegram" {
		t.Errorf("ParsedInputs.Channel = %s, expected telegram", ParsedInputs.Channel)
	}
	if ParsedInputs.TelegramApiKey != "test_telegram_key" {
		t.Errorf("ParsedInputs.TelegramApiKey = %s, expected test_telegram_key", ParsedInputs.TelegramApiKey)
	}
	if ParsedInputs.TelegramChannel != "@testchannel" {
		t.Errorf("ParsedInputs.TelegramChannel = %s, expected @testchannel", ParsedInputs.TelegramChannel)
	}

	// Check validation passes
	err := hasAllRequiredInputs("telegram")
	if err != nil {
		t.Errorf("hasAllRequiredInputs failed: %v", err)
	}

	// Simulate the send flow
	var commitMsg string
	if ParsedInputs.AddCommitInfo {
		commitMsg = templateCommitInfo()
	} else {
		commitMsg = "📦 *Github Workflow*\n\n"
	}
	finalMessage := commitMsg + ParsedInputs.Message

	expectedMessage := "📦 *Github Workflow*\n\nTelegram test"
	if finalMessage != expectedMessage {
		t.Errorf("Final message = %q, expected %q", finalMessage, expectedMessage)
	}
}
