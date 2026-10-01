package engine

import (
	"context"
	"fmt"
	"ubot/internal/conversation"

	"github.com/cloudwego/eino/schema"
)

type Generator interface {
	Generate(ctx context.Context, messages []*schema.Message, tools []*schema.ToolInfo) (*schema.Message, error)
}

type Platform interface {
	LoadConversation(ctx context.Context, conversationID string) (*conversation.Conversation, error)
	GetAgentSnapshotByVersion(ctx context.Context, agentVersionID string) (*Snapshot, error)
}

type Snapshot struct {
	ID           string
	AgentID      string
	WorkspaceID  string
	SystemPrompt string
	MaxSteps     int
}

type Engine struct {
	platform Platform
	gen      Generator
}

func New(platform Platform, gen Generator) *Engine {
	return &Engine{
		platform: platform,
		gen:      gen,
	}
}

func (e *Engine) ResolveSnapshot(ctx context.Context, conversationID string) (*Snapshot, error) {
	if e.platform == nil {
		return nil, fmt.Errorf("platform is required")
	}

	conversation, err := e.platform.LoadConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("load conversation: %w", err)
	}
	if conversation.AgentVersionID == "" {
		return nil, fmt.Errorf("conversation has no pinned agent version")
	}

	snapshot, err := e.platform.GetAgentSnapshotByVersion(ctx, conversation.AgentVersionID)
	if err != nil {
		return nil, fmt.Errorf("get agent snapshot: %w", err)
	}
	if snapshot.ID != conversation.AgentVersionID {
		return nil, fmt.Errorf(
			"snapshot version mismatch: want %s, got %s",
			conversation.AgentVersionID,
			snapshot.ID,
		)
	}
	return snapshot, nil
}
