package listmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const (
	apiPathShow = "/api/show"
	ONE_GB      = 1_073_741_824 // 1024 * 1024 * 1024
)

func GetModelInfo(cfg *settings.Settings, modelName string) (*ollama.Model, error) {
	var (
		client = http.Client{Timeout: time.Duration(1) * time.Second}
		err    error
		req    *http.Request
		res    *http.Response
		model  = &ollama.Model{}
		body   []byte
		raw    string
	)

	if req, err = http.NewRequest(
		http.MethodPost,
		cfg.OllamaUrl+apiPathShow,
		bytes.NewReader([]byte(`{"model": "`+modelName+`"}`))); err != nil {
		return nil, fmt.Errorf("GetModel creating request: %v\n", err)
	}

	req.Header.Add("Accept", "application/josn")
	if res, err = client.Do(req); err != nil {
		return nil, fmt.Errorf("GetModel calling api: %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetModel status code: %d\n", res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("GetModel reading response: %+v", err)
	}

	if raw, err = normalizeFamilyFields(string(body)); err != nil {
		return nil, fmt.Errorf("GetModel normalizing response: %v", err)
	}

	if err := json.Unmarshal([]byte(raw), model); err != nil {
		return nil, fmt.Errorf("GetModel parsing response: %v", err)
	}

	// model.Details.ParseQuantizationLevel()

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
/* examples:
{
	"details": {
		"family": "llama",
	},
	"model_info": {
		"llama.context_length": 131072,
		"llama.embedding_length": 4096,
	}
}
{
	"details": {
		"family": "qwen2",
	},
	"model_info": {
		"qwen2.context_length": 131072,
		"qwen2.embedding_length": 5120
	}
}
*/
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
