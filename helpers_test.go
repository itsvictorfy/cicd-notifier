package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestReadInputs(t *testing.T) {
	// Set up environment variables
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "test message")
	os.Setenv("INPUT_API_KEY", "api_key")
	os.Setenv("INPUT_CHANNEL_ID", "channel_id")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "true")
	os.Setenv("INPUT_COMMIT_SHA", "abc123")
	os.Setenv("INPUT_BRANCH", "main")
	os.Setenv("INPUT_AUTHOR", "testuser")
	os.Setenv("INPUT_COMMIT_TIME", "2023-01-01")
	os.Setenv("INPUT_COMMIT_MSG", "test commit")
	os.Setenv("INPUT_WORKFLOW_NAME", "test workflow")
	os.Setenv("INPUT_IMAGE_TAG", "v1.0")
	os.Setenv("INPUT_MSGID", "msg123")
	defer func() {
		// Clean up
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_API_KEY")
		os.Unsetenv("INPUT_CHANNEL_ID")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
		os.Unsetenv("INPUT_COMMIT_SHA")
		os.Unsetenv("INPUT_BRANCH")
		os.Unsetenv("INPUT_AUTHOR")
		os.Unsetenv("INPUT_COMMIT_TIME")
		os.Unsetenv("INPUT_COMMIT_MSG")
		os.Unsetenv("INPUT_WORKFLOW_NAME")
		os.Unsetenv("INPUT_IMAGE_TAG")
		os.Unsetenv("INPUT_MSGID")
	}()

	inputs := readInputs()

	expected := map[string]string{
		"action":          "send",
		"message":         "test message",
		"api_key":         "api_key",
		"channel_id":      "channel_id",
		"channel":         "slack",
		"add_commit_info": "true",
		"commit_sha":      "abc123",
		"branch":          "main",
		"author":          "testuser",
		"commit_time":     "2023-01-01",
		"commit_msg":      "test commit",
		"workflow_name":   "test workflow",
		"image_tag":       "v1.0",
		"msgid":           "msg123",
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

	expected := "📦 *Github Workflow*\n\n📌 *Commit:* `abc123`\n🔖 *Branch:* `main`\n🛠️ *Workflow:* `CI`\n📝 *Message:* fix bug\n👤 *Author:* user\n🐳 *Image Tag:* v1.0\n🕗 *Commit Time:* 2023-01-01\n\n"
	if result != expected {
		t.Errorf("templateCommitInfo() = %q, expected %q", result, expected)
	}
}

func TestReadInputsTemplateCommitInfoAndAddMessage(t *testing.T) {
	// Set up environment variables
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Deployment completed successfully")
	os.Setenv("INPUT_API_KEY", "api_key")
	os.Setenv("INPUT_CHANNEL_ID", "channel_id")
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
		os.Unsetenv("INPUT_API_KEY")
		os.Unsetenv("INPUT_CHANNEL_ID")
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
		Action:    inputs["action"],
		Message:   inputs["message"],
		ApiKey:    inputs["api_key"],
		ChannelId: inputs["channel_id"],
		Channel:   inputs["channel"],
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

	expected := "📦 *Github Workflow*\n\n📌 *Commit:* `def456`\n🔖 *Branch:* `develop`\n🛠️ *Workflow:* `Deploy`\n📝 *Message:* update feature\n👤 *Author:* developer\n🐳 *Image Tag:* v2.0\n🕗 *Commit Time:* 2023-01-02\n\nDeployment completed successfully"
	if finalMessage != expected {
		t.Errorf("Final message = %q, expected %q", finalMessage, expected)
	}
}

func TestSlackSendFlow(t *testing.T) {
	// Set up environment variables for Slack send
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Test message")
	os.Setenv("INPUT_API_KEY", "test_slack_key")
	os.Setenv("INPUT_CHANNEL_ID", "#test")
	os.Setenv("INPUT_CHANNEL", "slack")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "false")
	defer func() {
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_API_KEY")
		os.Unsetenv("INPUT_CHANNEL_ID")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
	}()

	// Simulate init logic
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
	if ParsedInputs.ApiKey != "test_slack_key" {
		t.Errorf("ParsedInputs.ApiKey = %s, expected test_slack_key", ParsedInputs.ApiKey)
	}
	if ParsedInputs.ChannelId != "#test" {
		t.Errorf("ParsedInputs.ChannelId = %s, expected #test", ParsedInputs.ChannelId)
	}
	if ParsedInputs.AddCommitInfo != false {
		t.Errorf("ParsedInputs.AddCommitInfo = %v, expected false", ParsedInputs.AddCommitInfo)
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
	os.Setenv("INPUT_API_KEY", "test_slack_key")
	os.Setenv("INPUT_CHANNEL_ID", "#test")
	os.Setenv("INPUT_MSGID", "123456")
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
		os.Unsetenv("INPUT_API_KEY")
		os.Unsetenv("INPUT_CHANNEL_ID")
		os.Unsetenv("INPUT_MSGID")
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
	if ParsedInputs.MsgID != "123456" {
		t.Errorf("ParsedInputs.MsgID = %s, expected 123456", ParsedInputs.MsgID)
	}
	if !ParsedInputs.AddCommitInfo {
		t.Errorf("ParsedInputs.AddCommitInfo = %v, expected true", ParsedInputs.AddCommitInfo)
	}

	// Simulate the update flow
	var commitMsg string
	if ParsedInputs.AddCommitInfo {
		commitMsg = templateCommitInfo()
	} else {
		commitMsg = "📦 *Github Workflow*\n\n"
	}
	finalMessage := commitMsg + ParsedInputs.Message

	expectedMessage := "📦 *Github Workflow*\n\n📌 *Commit:* `abc123`\n🔖 *Branch:* `main`\n🛠️ *Workflow:* `CI`\n📝 *Message:* fix\n👤 *Author:* user\n🐳 *Image Tag:* v1.0\n🕗 *Commit Time:* 2023-01-01\n\nUpdated message"
	if finalMessage != expectedMessage {
		t.Errorf("Final message = %q, expected %q", finalMessage, expectedMessage)
	}
}

func TestTelegramSendFlow(t *testing.T) {
	// Set up environment variables for Telegram send
	os.Setenv("INPUT_ACTION", "send")
	os.Setenv("INPUT_MESSAGE", "Telegram test")
	os.Setenv("INPUT_API_KEY", "test_telegram_key")
	os.Setenv("INPUT_CHANNEL_ID", "@testchannel")
	os.Setenv("INPUT_CHANNEL", "telegram")
	os.Setenv("INPUT_ADD_COMMIT_INFO", "false")
	defer func() {
		os.Unsetenv("INPUT_ACTION")
		os.Unsetenv("INPUT_MESSAGE")
		os.Unsetenv("INPUT_API_KEY")
		os.Unsetenv("INPUT_CHANNEL_ID")
		os.Unsetenv("INPUT_CHANNEL")
		os.Unsetenv("INPUT_ADD_COMMIT_INFO")
	}()

	// Simulate init logic
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
	if ParsedInputs.ApiKey != "test_telegram_key" {
		t.Errorf("ParsedInputs.ApiKey = %s, expected test_telegram_key", ParsedInputs.ApiKey)
	}
	if ParsedInputs.ChannelId != "@testchannel" {
		t.Errorf("ParsedInputs.ChannelId = %s, expected @testchannel", ParsedInputs.ChannelId)
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
