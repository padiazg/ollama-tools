package listmodels

import (
	"fmt"
	"net/http"

	"github.com/padiazg/ollama-tools/internals/client"
	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const apiPathTags = "/api/tags"

// GetTags calls tags api and returns a list of models
// func GetTags(cfg *settings.Settings) (*ollama.Tags, error) {
// 	var (
// 		client = http.Client{Timeout: time.Duration(1) * time.Second}
// 		err    error
// 		req    *http.Request
// 		res    *http.Response
// 		tags   = &ollama.Tags{}
// 	)

// 	if req, err = http.NewRequest(http.MethodGet, cfg.OllamaUrl+apiPathTags, nil); err != nil {
// 		return nil, fmt.Errorf("GetTags creating request: %v\n", err)
// 	}

// 	req.Header.Add("Accept", "application/josn")
// 	if res, err = client.Do(req); err != nil {
// 		return nil, fmt.Errorf("GetTags calling api: %v\n", err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("GetTags status code: %d\n", res.StatusCode)
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("GetTags reading stream: %+v", err)
// 	}

// 	if err = json.Unmarshal(body, tags); err != nil {
// 		return nil, fmt.Errorf("GetTags unmarshaling: %+v", err)
// 	}

// 	return tags, nil
// }

func GetTags(cfg *settings.Settings) (*ollama.Tags, error) {
	var (
		tags       = &ollama.Tags{}
		err        error
		req        *http.Request
		restClient = client.New(nil)
	)

	if req, err = http.NewRequest(http.MethodGet, cfg.OllamaUrl+apiPathTags, nil); err != nil {
		return nil, fmt.Errorf("GetTags creating request: %v\n", err)
	}

	if err = restClient.Request(req, tags); err != nil {
		return nil, fmt.Errorf("GetTags calling the api: %v\n", err)
	}

	return tags, err
}
