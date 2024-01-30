package v1

import (
	"context"

	"connectrpc.com/connect"
	studioerrors "github.com/missingstudio/studio/backend/internal/errors"
	"github.com/missingstudio/studio/backend/internal/providers"
	"github.com/missingstudio/studio/backend/internal/providers/base"
	llmv1 "github.com/missingstudio/studio/protos/pkg/llm"
)

func (s *V1Handler) ChatCompletions(
	ctx context.Context,
	req *connect.Request[llmv1.CompletionRequest],
) (*connect.Response[llmv1.CompletionResponse], error) {
	provider, err := providers.GetProvider(ctx)
	if err != nil {
		return nil, studioerrors.NewNotFound("provider not found")
	}

	completionProvider, ok := provider.(base.ChatCompilationInterface)
	if !ok {
		return nil, studioerrors.NewInternalError("not able to get chat compilation provider")
	}

	data, err := completionProvider.ChatCompilation(ctx, req.Msg)
	if err != nil {
		return nil, studioerrors.New(err)
	}

	return connect.NewResponse(data), nil
}