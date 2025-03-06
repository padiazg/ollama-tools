package listmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/padiazg/ollama-tools/internals/client"
	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const (
	apiPathShow = "/api/show"
	ONE_GB      = 1_073_741_824 // 1024 * 1024 * 1024
)

// HTTPClient to help create mocks at the tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func GetModelInfo(cfg *settings.Settings, modelName string) (*ollama.Model, error) {
	var (
		err        error
		req        *http.Request
		model      = &ollama.Model{}
		restClient = client.New(&client.ResClientConfig{
			OnDecode: func(rc io.ReadCloser, i interface{}) error {
				var (
					body []byte
					raw  string
					err  error
				)

				body, err = io.ReadAll(rc)
				if err != nil {
					return fmt.Errorf("GetModelInfo/OnDecode reading response: %+v", err)
				}

				if raw, err = normalizeFamilyFields(string(body)); err != nil {
					return fmt.Errorf("GetModelInfo/OnDecode normalizing response: %v", err)
				}

				if err := json.Unmarshal([]byte(raw), i); err != nil {
					return fmt.Errorf("GetModelInfo/OnDecode parsing response: %v", err)
				}

				return nil
			},
		})
	)

	if req, err = http.NewRequest(
		http.MethodPost,
		cfg.OllamaUrl+apiPathShow,
		bytes.NewReader([]byte(`{"model": "`+modelName+`"}`))); err != nil {
		return nil, fmt.Errorf("GetModelInfo creating request: %v\n", err)
	}

	if err = restClient.Request(req, model); err != nil {
		return nil, fmt.Errorf("GetTags calling the api: %v\n", err)
	}

	return model, nil

}

// getFamily recovers value for {"details": {"family": "..."}} from the unprocessed response
func getFamily(raw string) (string, error) {
	var (
		re  = regexp.MustCompile(`(?U)"family":\s?"(.*)",`)
		res = re.FindAllStringSubmatch(string(raw), -1)
	)
	if len(res) < 1 {
		return "", fmt.Errorf("no family found")
	}

	return res[0][1], nil
}

// normalizeFamilyFields normalize fields that starts with the family name at {"model_info": {...}} so
// those field/values ara unmarshaled to prefedined fields (ContextLength, EmbeddingLength) in the struct
func normalizeFamilyFields(raw string) (string, error) {
	var (
		family string
		err    error
	)

	if family, err = getFamily(raw); err != nil {
		return "", fmt.Errorf("normalizeFamilyFields %+v", err)
	}

	re := regexp.MustCompile(fmt.Sprintf(`(?m)(?U)%s\.(context_length|embedding_length)`, family))
	raw = re.ReplaceAllString(raw, "model.$1")

	return raw, nil
}
