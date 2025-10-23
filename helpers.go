package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var outputs = make(map[string]string)

// addOutput adds a key-value pair to the outputs map
func addOutput(key, value string) {
	outputs[key] = value
}

// readInputs reads all INPUT_* environment variables and returns a map
// To add new params, just use inputs["new_param_name"] in your code
func readInputs() map[string]string {
	inputs := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "INPUT_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := strings.ToLower(strings.TrimPrefix(parts[0], "INPUT_"))
				value := parts[1]
				inputs[key] = value
			}
		}
	}
	return inputs
}

// setOutputs writes the global outputs map to GITHUB_OUTPUT or prints to console
// To add new outputs, use addOutput(key, value)
func setOutputs() {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		// Not running in GitHub Actions, print to console
		fmt.Println("Outputs:")
		for key, value := range outputs {
			fmt.Printf("%s=%s\n", key, value)
		}
		return
	}

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening output file: %v\n", err)
		return
	}
	defer file.Close()

	for key, value := range outputs {
		fmt.Fprintf(file, "%s=%s\n", key, value)
	}
}

func validateInputs() {
	if ParsedInputs.Action != "send" && ParsedInputs.Action != "update" {
		slog.Error("Wrong Operation", "action", ParsedInputs.Action)
		os.Exit(1)
	}

	// Validate channel and required inputs
	switch strings.ToLower(ParsedInputs.Channel) {
	case "slack":
		if err := hasAllRequiredInputs("slack"); err != nil {
			slog.Error("Missing required inputs for Slack", "error", err)
			os.Exit(1)
		}
	case "telegram":
		if err := hasAllRequiredInputs("telegram"); err != nil {
			slog.Error("Missing required inputs for Telegram", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("Unsupported channel specified", "channel", ParsedInputs.Channel)
		os.Exit(1)
	}
	fmt.Printf("%v\n", ParsedInputs)
}

func hasAllRequiredInputs(channel string) error {
	if ParsedInputs.Message == "" {
		return fmt.Errorf("message is required")
	}
	switch strings.ToLower(channel) {
	case "slack":
		if ParsedInputs.SlackApiKey == "" {
			return fmt.Errorf("slack_api_key is required for Slack")
		}
		if ParsedInputs.SlackChannel == "" {
			return fmt.Errorf("slack_channel is required for Slack")
		}
		if ParsedInputs.Action == "update" && ParsedInputs.SlackMsgID == "" {
			return fmt.Errorf("slack_msgid is required for Slack update")
		}
	case "telegram":
		if ParsedInputs.TelegramApiKey == "" {
			return fmt.Errorf("telegram_api_key is required for Telegram")
		}
		if ParsedInputs.TelegramChannel == "" {
			return fmt.Errorf("telegram_channel is required for Telegram")
		}
		if ParsedInputs.Action == "update" && ParsedInputs.TelegramMsgID == "" {
			return fmt.Errorf("telegram_msgid is required for Telegram update")
		}
	default:
		return fmt.Errorf("unsupported channel: %s", channel)
	}
	return nil
}

func templateCommitInfo() string {
	msg := "📦 *Github Workflow*\n\n"

	if ParsedInputs.CommitSha != "" {
		msg += fmt.Sprintf("📌 *Commit:* `%s`\n", ParsedInputs.CommitSha)
	}
	if ParsedInputs.Branch != "" {
		msg += fmt.Sprintf("🔖 *Branch:* `%s`\n", ParsedInputs.Branch)
	}
	if ParsedInputs.WorkflowName != "" {
		msg += fmt.Sprintf("🛠️ *Workflow:* `%s`\n", ParsedInputs.WorkflowName)
	}
	if ParsedInputs.CommitMsg != "" {
		msg += fmt.Sprintf("📝 *Message:* %s\n", ParsedInputs.CommitMsg)
	}
	if ParsedInputs.Author != "" {
		msg += fmt.Sprintf("👤 *Author:* %s\n", ParsedInputs.Author)
	}
	if ParsedInputs.ImageTag != "" {
		msg += fmt.Sprintf("🐳 *Image Tag:* %s\n", ParsedInputs.ImageTag)
	}
	if ParsedInputs.CommitTime != "" {
		msg += fmt.Sprintf("🕗 *Commit Time:* %s\n", ParsedInputs.CommitTime)
	}
	return msg
}
