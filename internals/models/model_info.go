package models

import (
	"encoding/json"
	"fmt"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
	"resty.dev/v3"
)

const (
	apiPathShow = "/api/show"
	ONE_GB      = 1_073_741_824 // 1024 * 1024 * 1024
)

func GetModelInfo(cfg *settings.Settings, model_name string) (*ollama.Model, error) {
	var (
		err   error
		model = &ollama.Model{}
		c     = resty.New()
	)
	defer c.Close()

	res, err := c.R().
		SetBody([]byte(`{"model": "` + model_name + `"}`)).
		Post(cfg.OllamaUrl + apiPathShow)
	if err != nil {
		return nil, fmt.Errorf("requesting model %s info list: %+v", model_name, err)
	}
	defer res.Body.Close()

	if !res.IsSuccess() {
		return nil, fmt.Errorf("response status code: %d", res.StatusCode())
	}

	if err := json.NewDecoder(res.Body).Decode(model); err != nil {
		return nil, fmt.Errorf("decoding response: %+v", err)
	}

	return model, nil
}
