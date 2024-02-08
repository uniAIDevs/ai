package v1

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect"
	"github.com/missingstudio/studio/backend/internal/providers"
	"github.com/missingstudio/studio/backend/internal/providers/base"
	"github.com/missingstudio/studio/common/errors"
	llmv1 "github.com/missingstudio/studio/protos/pkg/llm"
)

var (
	ErrChatCompletionStreamNotSupported = errors.NewBadRequest("streaming is not supported with this method, please use StreamChatCompletions")
	ErrChatCompletionNotSupported       = errors.NewInternalError("provider don't have chat Completion capabilities")
)

func (s *V1Handler) ChatCompletions(
	ctx context.Context,
	req *connect.Request[llmv1.ChatCompletionRequest],
) (*connect.Response[llmv1.ChatCompletionResponse], error) {
	provider, err := providers.GetProvider(ctx, req.Header())
	if err != nil {
		return nil, errors.New(err)
	}

	if err := provider.Validate(); err != nil {
		return nil, errors.New(err)
	}

	chatCompletionProvider, ok := provider.(base.ChatCompletionInterface)
	if !ok {
		return nil, ErrChatCompletionNotSupported
	}

	payload, err := json.Marshal(req.Msg)
	if err != nil {
		return nil, errors.New(err)
	}

	resp, err := chatCompletionProvider.ChatCompletion(ctx, payload)
	if err != nil {
		return nil, errors.New(err)
	}

	data := &llmv1.ChatCompletionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, errors.New(err)
	}

	ingesterdata := make(map[string]interface{})
	ingesterdata["provider"] = provider.GetName()
	ingesterdata["model"] = data.Model
	ingesterdata["usage"] = map[string]interface{}{
		"total_tokens":      data.Usage.TotalTokens,
		"prompt_tokens":     data.Usage.PromptTokens,
		"completion_tokens": data.Usage.CompletionTokens,
	}

	go s.ingester.Ingest(ingesterdata, "logs")
	return connect.NewResponse(data), nil
}
