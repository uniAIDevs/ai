package provider

import (
	"encoding/json"

	"github.com/Jeffail/gabs/v2"
	"github.com/google/uuid"
)

type ProviderState string

const (
	ProviderStateActive   ProviderState = "active"
	ProviderStateInactive ProviderState = "inactive"
)

func (s ProviderState) String() string {
	return string(s)
}

type Provider struct {
	ID     uuid.UUID      `json:"id"`
	Name   string         `json:"name"`
	State  ProviderState  `json:"state"`
	Config map[string]any `json:"config"`
}

func (c *Provider) MergeConfig(kv map[string]any) (err error) {
	container := gabs.New()
	for k, v := range kv {
		_, err = container.SetP(v, k)
		if err != nil {
			return err
		}
	}
	if err := container.MergeFn(gabs.Wrap(c.Config), func(destination, source any) any {
		return destination
	}); err != nil {
		return err
	}

	// get back connection config
	containerBytes, err := container.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(containerBytes, &c.Config)
}

func (c *Provider) GetConfig(keys []string) map[string]any {
	fetched := map[string]any{}
	container := gabs.Wrap(c.Config)
	for _, key := range keys {
		fetched[key] = container.Path(key).Data()
	}
	return fetched
}
