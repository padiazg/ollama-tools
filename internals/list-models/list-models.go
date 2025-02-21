package listmodels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const apiPath = "/api/tags"

func List(cfg *settings.Settings) {
	var (
		client = http.Client{Timeout: time.Duration(1) * time.Second}
		err    error
		req    *http.Request
		res    *http.Response
		body   []byte
		tags   ollama.Tags
	)

	if req, err = http.NewRequest(http.MethodGet, cfg.OllamaUrl+apiPath, nil); err != nil {
		fmt.Printf("creating request: %v\n", err)
		return
	}

	req.Header.Add("Accept", "application/josn")
	if res, err = client.Do(req); err != nil {
		fmt.Printf("executing request: %v\n", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code: %d\n", res.StatusCode)
		return
	}

	if body, err = io.ReadAll(res.Body); err != nil {
		fmt.Printf("reading body: %v\n", err)
	}

	json.Unmarshal(body, &tags)

	for _, model := range tags.Models {
		fmt.Printf("name: %s, model: %s, parameter_size: %s (%.6f), quantization_level: %s (%d)\n",
			model.Name,
			model.Model,
			model.Details.ParameterSize,
			model.Details.ParameterSizeAsBillions(),
			model.Details.QuantizationLevel,
			model.Details.QuantizationLevelAsBitCount(),
		)

	}
}

func parseParameterSize(size string) float32 {
	var (
		re  = regexp.MustCompile(`(?m)([0-9]+(\.[0-9])?)([BM])`)
		res = re.FindAllStringSubmatch(size, -1)
	)

	if len(res) < 1 {
		fmt.Printf("parseParameterSize error parsing: %s\n", size)
		return 0.0
	}

	// extract size
	value, err := strconv.ParseFloat(res[0][1], 32)
	if err != nil {
		fmt.Printf("parseParameterSize error converting value: %s\n", res[0][1])
	}

	// scale if needed
	switch res[0][3] {
	case "M":
		value /= 1024.0
	}

	return float32(value)
}
