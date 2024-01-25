package providers

import (
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/missingstudio/studio/backend/internal/providers/base"
	"github.com/missingstudio/studio/backend/internal/providers/openai"
)

type ProviderFactory interface {
	Create(token string) base.LLMProvider
}

var providerFactories = make(map[string]ProviderFactory)

func init() {
	providerFactories["openai"] = openai.OpenAIProviderFactory{}
}

func GetProvider(headers http.Header) (base.LLMProvider, error) {
	provider := headers.Get("x-ms-provider")
	providerFactory, ok := providerFactories[provider]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("provider not found"))
	}

	authHeader := headers.Get("Authorization")
	accessToken := strings.Replace(authHeader, "Bearer ", "", 1)

	return providerFactory.Create(accessToken), nil
}