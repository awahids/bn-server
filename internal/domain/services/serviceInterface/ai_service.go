package serviceinterface

import "context"

type AIService interface {
	GetCoachResponse(ctx context.Context, systemPrompt string, userMessage string) (string, error)
}
