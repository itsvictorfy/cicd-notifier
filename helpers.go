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
	if ParsedInputs.Message == "" {
		slog.Error("message is required")
		os.Exit(1)
	}
	if ParsedInputs.ApiKey == "" {
		slog.Error("ApiKey is required")
		os.Exit(1)
	}
	if ParsedInputs.Channel == "" {
		slog.Error("channel is required")
		os.Exit(1)
	}
	if ParsedInputs.ChannelId == "" {
		slog.Error("channelId is required")
		os.Exit(1)
	}
	if ParsedInputs.Action == "" {
		slog.Error("action is missing")
		os.Exit(1)
	}
	if ParsedInputs.Action == "update" && ParsedInputs.MsgID == "" {
		slog.Error("Missing MsgId for Action-Update")
		os.Exit(1)
	}
	fmt.Printf("%v\n", ParsedInputs)
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
	msg += "\n"
	return msg
}
