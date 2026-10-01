package agent

import "time"

type Version struct {
	ID           string
	AgentID      string
	WorkspaceID  string
	Version      int
	SystemPrompt string
	CreatedAt    time.Time
}