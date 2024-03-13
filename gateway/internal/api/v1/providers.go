package v1

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect"
	"github.com/jeremywohl/flatten"
	"github.com/missingstudio/ai/gateway/core/provider"
	llmv1 "github.com/missingstudio/ai/protos/pkg/llm/v1"
	"github.com/missingstudio/common/errors"

	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *V1Handler) ListProviders(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[llmv1.ListProvidersResponse], error) {
	providers := s.iProviderService.GetProviders()

	data := []*llmv1.Provider{}
	for _, provider := range providers {
		providerInfo := provider.Info()
		data = append(data, &llmv1.Provider{
			Title:       providerInfo.Title,
			Name:        providerInfo.Name,
			Description: providerInfo.Description,
		})
	}

	return connect.NewResponse(&llmv1.ListProvidersResponse{
		Providers: data,
	}), nil
}

func (s *V1Handler) GetProvider(ctx context.Context, req *connect.Request[llmv1.GetProviderRequest]) (*connect.Response[llmv1.GetProviderResponse], error) {
	provider, err := s.iProviderService.GetProvider(provider.Provider{Name: req.Msg.Name})
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	conn, err := s.providerService.GetByName(ctx, req.Msg.Name)
	if err != nil {
		return nil, errors.New(err)
	}

	stConfigs, err := structpb.NewStruct(conn.Config)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	info := provider.Info()
	p := &llmv1.Provider{
		Title:       info.Title,
		Name:        info.Name,
		Description: info.Description,
		Config:      stConfigs,
	}

	return connect.NewResponse(&llmv1.GetProviderResponse{
		Provider: p,
	}), nil
}

func (s *V1Handler) CreateProvider(ctx context.Context, req *connect.Request[llmv1.CreateProviderRequest]) (*connect.Response[llmv1.CreateProviderResponse], error) {
	provider := provider.Provider{Name: req.Msg.Name, Config: req.Msg.Config.AsMap()}

	p, err := s.iProviderService.GetProvider(provider)
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	providerSchema := gojsonschema.NewBytesLoader(p.Schema())
	connectionSchema := gojsonschema.NewGoLoader(provider.Config)

	result, err := gojsonschema.Validate(providerSchema, connectionSchema)
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	if !result.Valid() {
		var err error
		for _, desc := range result.Errors() {
			if desc.Type() == "required" {
				// ignore required validation checks in patch call
				continue
			}
			err = errors.NewBadRequest(desc.String())
		}

		if err != nil {
			return nil, errors.NewNotFound(err.Error())
		}
	}

	conn, err := s.providerService.Upsert(ctx, provider)
	if err != nil {
		return nil, errors.New(err)
	}

	stConfigs, err := structpb.NewStruct(conn.Config)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	info := p.Info()
	return connect.NewResponse(&llmv1.CreateProviderResponse{
		Name:   info.Name,
		Config: stConfigs,
	}), nil
}

func (s *V1Handler) UpsertProvider(ctx context.Context, req *connect.Request[llmv1.UpdateProviderRequest]) (*connect.Response[llmv1.UpdateProviderResponse], error) {
	provider := provider.Provider{Name: req.Msg.Name, Config: req.Msg.Config.AsMap()}

	p, err := s.iProviderService.GetProvider(provider)
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	providerSchema := gojsonschema.NewBytesLoader(p.Schema())
	connectionSchema := gojsonschema.NewGoLoader(provider.Config)

	result, err := gojsonschema.Validate(providerSchema, connectionSchema)
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	if !result.Valid() {
		var err error
		for _, desc := range result.Errors() {
			if desc.Type() == "required" {
				// ignore required validation checks in patch call
				continue
			}
			err = errors.NewBadRequest(desc.String())
		}

		if err != nil {
			return nil, errors.NewNotFound(err.Error())
		}
	}

	source, err := s.providerService.GetByName(ctx, req.Msg.Name)
	if err != nil {
		return nil, errors.New(err)
	}

	requiredMap, err := flatten.Flatten(provider.Config, "", flatten.DotStyle)
	if err != nil {
		return nil, errors.New(err)
	}
	if err := source.MergeConfig(requiredMap); err != nil {
		return nil, errors.New(err)
	}

	conn, err := s.providerService.Upsert(ctx, source)
	if err != nil {
		return nil, errors.New(err)
	}

	stConfigs, err := structpb.NewStruct(conn.Config)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	info := p.Info()
	return connect.NewResponse(&llmv1.UpdateProviderResponse{
		Name:   info.Name,
		Config: stConfigs,
	}), nil
}

func (s *V1Handler) GetProviderConfig(ctx context.Context, req *connect.Request[llmv1.GetProviderConfigRequest]) (*connect.Response[llmv1.GetProviderConfigResponse], error) {
	p, err := s.iProviderService.GetProvider(provider.Provider{Name: req.Msg.Name})
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}

	configs := map[string]any{}
	if err := json.Unmarshal(p.Schema(), &configs); err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	stConfigs, err := structpb.NewStruct(configs)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	return connect.NewResponse(&llmv1.GetProviderConfigResponse{
		Config: stConfigs,
	}), nil
}
