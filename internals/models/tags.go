package models

import (
	"encoding/json"
	"fmt"

	"resty.dev/v3"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const apiPathTags = "/api/tags"

func GetTags(cfg *settings.Settings) (*ollama.Tags, error) {
	var (
		tags = &ollama.Tags{}
		err  error
		c    = resty.New()
	)
	defer c.Close()

	res, err := c.R().
		SetHeader("Accept", "application/josn").
		Get(cfg.OllamaUrl + apiPathTags)
	if err != nil {
		return nil, fmt.Errorf("requesting tags list: %+v", err)
	}

	defer res.Body.Close()

	if !res.IsSuccess() {
		return nil, fmt.Errorf("response status code: %d", res.StatusCode())
	}

	if err := json.NewDecoder(res.Body).Decode(tags); err != nil {
		return nil, fmt.Errorf("decoding response: %+v", err)
	}

	return tags, err
}
