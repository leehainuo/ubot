package conversation

import "time"

type Conversation struct {
	ID             string
	AgentID        string
	AgentVersionID string
	UserID         string
	CreatedAt      time.Time
}