package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
	"github.com/stretchr/testify/assert"
)

const (
	modelPhi4     string = "phi4:latest"
	modelLlama3_1 string = "llama3.1:latest"
	modelDeepSeek string = "deepseek-r1:7b"
	modelQwen     string = "qwen2.5-coder:7b"
)

type tags map[string]*ollama.TagModel
type modelsInfo map[string]testModelsInfo

func (t tags) getBody() string {
	if len(t) == 0 {
		return ""
	}

	tmp := ollama.Tags{}
	for _, model := range t {
		tmp.Models = append(tmp.Models, *model)
	}

	body, _ := json.Marshal(tmp)
	return string(body)
}

func (t tags) getModels() (models []string) {
	for _, tag := range t {
		models = append(models, tag.Name)
	}

	return
}

func (t tags) filter(list []string) tags {
	ret := tags{}
	for _, key := range list {
		if value, ok := t[key]; ok {
			ret[key] = value
		}
	}
	return ret
}

func (t tags) getConfig(wantError error) *settings.Settings {
	body := t.getBody()
	return &settings.Settings{
		OllamaUrl: "http://ollama:11434",
		Transport: &DryRunTransport{RoundTripFn: generalRoundTripper(body, wantError)},
	}
}

func (m modelsInfo) filter(list []string) modelsInfo {
	ret := modelsInfo{}
	for _, key := range list {
		if value, ok := m[key]; ok {
			ret[key] = value
		}
	}
	return ret
}

func (m modelsInfo) getModels(list []string) []*ModelItem {
	modelItemList := make([]*ModelItem, 0, len(models))

	for _, model := range list {
		m, ok := m[model]
		if !ok {
			modelItemList = append(modelItemList, &ModelItem{
				Error: fmt.Errorf("%s not found", model),
			})
			continue
		}

		modelItemList = append(modelItemList, &ModelItem{
			Name:  model,
			Model: m.model,
		})
	}

	return modelItemList
}

type testModelsInfo struct {
	name       string
	normalized string
	model      *ollama.Model
}

var (
	tagsList = tags{
		modelDeepSeek: &ollama.TagModel{
			Name:       modelDeepSeek,
			Model:      "deepseek-r1:7b",
			ModifiedAt: "2025-03-11T15:37:14.423620816-03:00",
			Size:       4683075271,
			Digest:     "0a8c266910232fd3291e71e5ba1e058cc5af9d411192cf88b6d30e92b6e73163",
			Details: ollama.TagModelDetails{
				ParentModel:       "",
				Format:            "gguf",
				Family:            "qwen2",
				Families:          []string{"qwen2"},
				ParameterSize:     "7.6B",
				QuantizationLevel: "Q4_K_M",
			},
		},
		modelQwen: &ollama.TagModel{
			Name:       modelQwen,
			Model:      "qwen2.5-coder:7b",
			ModifiedAt: "2025-02-26T13:02:11.408112713-03:00",
			Size:       4683087519,
			Digest:     "2b0496514337a3d5901f1d253d01726c890b721e891335a56d6e08cedf3e2cb0",
			Details: ollama.TagModelDetails{
				ParentModel:       "",
				Format:            "gguf",
				Family:            "qwen2",
				Families:          []string{"qwen2"},
				ParameterSize:     "7.6B",
				QuantizationLevel: "Q4_K_M",
			},
		},
		modelPhi4: &ollama.TagModel{
			Name:       modelPhi4,
			Model:      "phi4:14b",
			ModifiedAt: "2025-02-14T14:57:14.375271104-03:00",
			Size:       9053116391,
			Digest:     "ac896e5b8b34a1f4efa7b14d7520725140d5512484457fab45d2a4ea14c69dba",
			Details: ollama.TagModelDetails{
				ParentModel:       "",
				Format:            "gguf",
				Family:            "phi3",
				Families:          []string{"phi3"},
				ParameterSize:     "14.7B",
				QuantizationLevel: "Q4_K_M",
			},
		},
		modelLlama3_1: &ollama.TagModel{
			Name:       modelLlama3_1,
			Model:      "llama3.1:latest",
			ModifiedAt: "2025-01-08T18:46:33.340224609-03:00",
			Size:       4920753328,
			Digest:     "46e0c10c039e019119339687c3c1757cc81b9da49709a3b3924863ba87ca666e",
			Details: ollama.TagModelDetails{
				ParentModel:       "",
				Format:            "gguf",
				Family:            "llama",
				Families:          []string{"llama"},
				ParameterSize:     "8.0B",
				QuantizationLevel: "Q4_K_M",
			},
		},
	}

	modelsList = modelsInfo{
		modelPhi4: {
			name:       modelPhi4,
			normalized: `{"license":"Microsoft...SOFTWARE.","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|im_start|\u003e\"\nstop                           \"\u003c|im_end|\u003e\"\nstop                           \"\u003c|im_sep|\u003e\"","template":"{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 -}}\n\u003c|im_start|\u003e{{ .Role }}\u003c|im_sep|\u003e\n{{ .Content }}{{ if not $last }}\u003c|im_end|\u003e\n{{ end }}\n{{- if and (ne .Role \"assistant\") $last }}\u003c|im_end|\u003e\n\u003c|im_start|\u003eassistant\u003c|im_sep|\u003e\n{{ end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"phi3","families":["phi3"],"parameter_size":"14.7B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"phi3","general.basename":"phi","general.file_type":15,"general.languages":["en"],"general.license":"mit","general.license.link":"https://huggingface.co/microsoft/phi-4/resolve/main/LICENSE","general.organization":"Microsoft","general.parameter_count":14659507200,"general.quantization_version":2,"general.size_label":"15B","general.tags":["phi","nlp","math","code","chat","conversational","text-generation"],"general.type":"model","general.version":"4","phi3.attention.head_count":40,"phi3.attention.head_count_kv":10,"phi3.attention.layer_norm_rms_epsilon":0.00001,"phi3.attention.sliding_window":131072,"phi3.block_count":40,"model.context_length":16384,"model.embedding_length":5120,"phi3.feed_forward_length":17920,"phi3.rope.dimension_count":128,"phi3.rope.freq_base":250000,"phi3.rope.scaling.original_context_length":16384,"tokenizer.ggml.bos_token_id":100257,"tokenizer.ggml.eos_token_id":100257,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.padding_token_id":100257,"tokenizer.ggml.pre":"dbrx","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-01-14T17:21:17.785607967-03:00"}`,
			model: &ollama.Model{
				Details: ollama.ModelDetails{
					ParentModel:       "",
					Format:            "gguf",
					Family:            "phi3",
					Families:          []string{"phi3"},
					ParameterSize:     "14.7B",
					QuantizationLevel: "Q4_K_M",
				},
				ModelInfo: ollama.ModelInfo{
					Type:            "model",
					ParameterCount:  14659507200,
					ContextLength:   16384,
					EmbeddingLength: 5120,
				},
			},
		},
		modelLlama3_1: {
			name:       modelLlama3_1,
			normalized: `{"license":"LLAMA 3.1 ...","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|start_header_id|\u003e\"\nstop                           \"\u003c|end_header_id|\u003e\"\nstop                           \"\u003c|eot_id|\u003e\"","template":"{{- if or .System .Tools }}\u003c|start_header_id|\u003esystem\u003c|end_header_id|\u003e\n{{- if .System }}\n\n{{ .System }}\n{{- end }}\n{{- if .Tools }}\n\nCutting Knowledge Date: December 2023\n\nWhen you receive a tool call response, use the output to format an answer to the orginal user question.\n\nYou are a helpful assistant with tool calling capabilities.\n{{- end }}\u003c|eot_id|\u003e\n{{- end }}\n{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 }}\n{{- if eq .Role \"user\" }}\u003c|start_header_id|\u003euser\u003c|end_header_id|\u003e\n{{- if and $.Tools $last }}\n\nGiven the following functions, please respond with a JSON for a function call with its proper arguments that best answers the given prompt.\n\nRespond in the format {\"name\": function name, \"parameters\": dictionary of argument name and its value}. Do not use variables.\n\n{{ range $.Tools }}\n{{- . }}\n{{ end }}\nQuestion: {{ .Content }}\u003c|eot_id|\u003e\n{{- else }}\n\n{{ .Content }}\u003c|eot_id|\u003e\n{{- end }}{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- else if eq .Role \"assistant\" }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n{{- if .ToolCalls }}\n{{ range .ToolCalls }}\n{\"name\": \"{{ .Function.Name }}\", \"parameters\": {{ .Function.Arguments }}}{{ end }}\n{{- else }}\n\n{{ .Content }}\n{{- end }}{{ if not $last }}\u003c|eot_id|\u003e{{ end }}\n{{- else if eq .Role \"tool\" }}\u003c|start_header_id|\u003eipython\u003c|end_header_id|\u003e\n\n{{ .Content }}\u003c|eot_id|\u003e{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"llama","families":["llama"],"parameter_size":"8.0B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"llama","general.basename":"Meta-Llama-3.1","general.file_type":15,"general.finetune":"Instruct","general.languages":["en","de","fr","it","pt","hi","es","th"],"general.license":"llama3.1","general.parameter_count":8030261312,"general.quantization_version":2,"general.size_label":"8B","general.tags":["facebook","meta","pytorch","llama","llama-3","text-generation"],"general.type":"model","llama.attention.head_count":32,"llama.attention.head_count_kv":8,"llama.attention.layer_norm_rms_epsilon":0.00001,"llama.block_count":32,"model.context_length":131072,"model.embedding_length":4096,"llama.feed_forward_length":14336,"llama.rope.dimension_count":128,"llama.rope.freq_base":500000,"llama.vocab_size":128256,"tokenizer.ggml.bos_token_id":128000,"tokenizer.ggml.eos_token_id":128009,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.pre":"llama-bpe","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-02-03T19:22:17.054410969-03:00"}`,
			model: &ollama.Model{
				Details: ollama.ModelDetails{
					ParentModel:       "",
					Format:            "gguf",
					Family:            "llama",
					Families:          []string{"llama"},
					ParameterSize:     "8.0B",
					QuantizationLevel: "Q4_K_M",
				},
				ModelInfo: ollama.ModelInfo{
					Type:            "model",
					ParameterCount:  8030261312,
					ContextLength:   131072,
					EmbeddingLength: 4096,
				},
			},
		},
	}

	tagsRoundTripper = func(body string, wantError error) func(r *http.Request) (*http.Response, error) {
		return func(r *http.Request) (*http.Response, error) {
			if wantError != nil {
				return nil, wantError
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}
	}

	modelListRoundTripper = func(wantErr error) func(r *http.Request) (*http.Response, error) {
		return func(r *http.Request) (*http.Response, error) {
			type payload struct {
				Model string `json:"model"`
			}

			if wantErr != nil {
				return nil, wantErr
			}

			body := &payload{}

			if err := json.NewDecoder(r.Body).Decode(body); err != nil {
				return nil, fmt.Errorf("decoding response: %+v", err)
			}

			data, ok := modelsList[body.Model]
			if !ok {
				return &http.Response{StatusCode: http.StatusNotFound}, nil
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(data.normalized)),
			}, nil
		}
	}

	generalRoundTripper = func(body string, wantError error) func(r *http.Request) (*http.Response, error) {
		return func(r *http.Request) (*http.Response, error) {
			var (
				res *http.Response
				err error
			)

			switch r.URL.Path {
			case apiPathTags:
				return tagsRoundTripper(body, wantError)(r)
			case apiPathShow:
				return modelListRoundTripper(wantError)(r)
			}

			return res, err
		}
	}

	getConfigModelsList = func(wantError error) *settings.Settings {
		return &settings.Settings{
			OllamaUrl: "http://ollama:11434",
			Transport: &DryRunTransport{RoundTripFn: modelListRoundTripper(wantError)},
		}
	}
)

type DryRunTransport struct {
	http.RoundTripper
	RoundTripFn func(r *http.Request) (*http.Response, error)
}

func (dr *DryRunTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return dr.RoundTripFn(r)
}

func TestModelsInfoList(t *testing.T) {
	var (
		tests = []struct {
			name         string
			model_name   string
			models       []string
			wantErrorMsg string
		}{
			{
				name:   modelPhi4,
				models: []string{modelPhi4},
			},
			{
				name:   modelPhi4 + "+" + modelLlama3_1,
				models: []string{modelPhi4, modelLlama3_1},
			},
			{
				name:         "error",
				wantErrorMsg: "test-ModelsInfoList-error",
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					var (
						tags           tags
						cfg            *settings.Settings
						want           []*ModelItem
						wantError      error
						wantModelsList []string
					)

					if tt.model_name == "" {
						wantModelsList = tt.models
					} else {
						wantModelsList = []string{tt.model_name}
					}

					want = modelsList.getModels(wantModelsList)
					tags = tagsList.filter(wantModelsList)

					if tt.wantErrorMsg != "" {
						wantError = fmt.Errorf(tt.wantErrorMsg)
					}

					cfg = tags.getConfig(wantError)
					got, err := ModelsInfoList(cfg, tt.model_name)

					if wantError != nil {
						assert.ErrorContains(t, err, tt.wantErrorMsg)
					} else {
						if assert.NoError(t, err, "error not expected: %+v", err) {
							assert.EqualValues(t, got, want)
						}
					}
				})
			}
		})
	}
}

func Test_modelsInfoGenerator(t *testing.T) {
	type checkModelsInfoGeneratorFn func(*testing.T, nextFn)

	var (
		checkModels = func(models []string) checkModelsInfoGeneratorFn {
			return func(t *testing.T, next nextFn) {
				t.Helper()
				var (
					count    = 0
					expected = len(models)
					got      = []string{}
				)

				for {
					nextModel := next()
					if nextModel.model_name == "" {
						break
					}

					for _, m := range models {
						got = append(got, nextModel.model_name)
						if m == nextModel.model_name {
							count++
						}
					}
				}

				assert.Equalf(t, expected, count, "checkModels count=%+v, expected=%+v\ngot:%+v\nexpected:%+v", expected, count, got, models)
			}
		}

		tests = []struct {
			name         string
			model_name   string
			models       []string
			wantErrorMsg string
		}{
			{
				name:   modelDeepSeek,
				models: []string{modelDeepSeek},
			},
			{
				name:   modelDeepSeek + "+" + modelQwen,
				models: []string{modelDeepSeek, modelQwen},
			},
			{
				name:       "specific-model",
				model_name: modelPhi4,
			},
			{
				name:         "http-client-error",
				wantErrorMsg: "from http-client-mock",
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				tags      tags
				cfg       *settings.Settings
				want      []string
				wantError error
			)

			if tt.model_name == "" {
				tags = tagsList.filter(tt.models)
				want = tags.getModels()
			} else {
				want = []string{tt.model_name}
			}

			if tt.wantErrorMsg != "" {
				wantError = fmt.Errorf(tt.wantErrorMsg)
			}

			cfg = tags.getConfig(wantError)
			next, err := modelsInfoGenerator(cfg, tt.model_name)
			if wantError != nil {
				assert.ErrorContains(t, err, tt.wantErrorMsg)
			} else {
				if assert.NoError(t, err, "error not expected: %+v", err) {
					checkModels(want)(t, next)
				}
			}
		})
	}
}

func Test_modelsInfoList(t *testing.T) {
	var (
		generator = func(models []string) (nextFn, []*ModelItem) {
			var (
				modelItemList = modelsList.getModels(models)
				index         = 0
			)

			return func() nextData {
				if index == len(modelItemList) {
					return nextData{model_name: ""}
				}

				data := nextData{
					model_name: modelItemList[index].Name,
					cfg:        getConfigModelsList(nil),
				}
				index++
				return data
			}, modelItemList
		}

		tests = []struct {
			name   string
			models []string
		}{
			{
				name:   modelPhi4,
				models: []string{modelPhi4},
			},
			{
				name:   modelPhi4 + "+" + modelLlama3_1,
				models: []string{modelPhi4, modelLlama3_1},
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			next, want := generator(tt.models)
			got := modelsInfoList(next)
			assert.EqualValues(t, got, want)
		})
	}
}

func Test_modelsInfoReader(t *testing.T) {
	type checkModelsInfoReaderFn func(t *testing.T, pending <-chan nextData)

	var (
		checkModelName = func(model_name string) checkModelsInfoReaderFn {
			return func(t *testing.T, pending <-chan nextData) {
				t.Helper()
				fetched := make(chan nextData)

				go func() {
					for {
						data, ok := <-pending
						if !ok {
							close(fetched)
							break
						}
						fetched <- data
					}
				}()

				value, ok := <-fetched
				assert.Truef(t, ok, "channel expected to be open, currently closed")
				assert.Equal(t, model_name, value.model_name)
			}
		}

		checkEndReading = func() checkModelsInfoReaderFn {
			return func(t *testing.T, pending <-chan nextData) {
				t.Helper()

				go func() {
					for {
						_, ok := <-pending
						assert.Falsef(t, ok, "channel expected to be closed, currently open")
						break
					}
				}()
			}
		}

		// we need a generator to send the empty value and trigger the
		// end of the reader go routine and the close of the pending channel
		generator = func(models []string) nextFn {
			index := 0
			return func() nextData {
				if index == len(models) {
					return nextData{model_name: ""}
				}
				data := nextData{model_name: models[index]}
				index++
				return data
			}
		}

		tests = []struct {
			name       string
			model_name string
			check      checkModelsInfoReaderFn
			next       nextFn
		}{
			{
				name:       "success-one-model",
				model_name: "test",
				check:      checkModelName("test"),
				next:       generator([]string{"test"}),
			},
			{
				name:       "end-reading",
				model_name: "",
				check:      checkEndReading(),
				next: func() nextData {
					return nextData{model_name: ""}
				},
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			pending := modelsInfoReader(tt.next)
			tt.check(t, pending)
		})
	}
}

func Test_modelsInfoFetcher(t *testing.T) {
	type checkModelsInfoReaderFn func(t *testing.T, fetched <-chan pair)

	var (
		checkModel = func(model *ollama.Model) checkModelsInfoReaderFn {
			return func(t *testing.T, fetched <-chan pair) {
				data, ok := <-fetched
				assert.True(t, ok)
				assert.EqualValues(t, model, data.model)
			}
		}

		checkError = func(errorMsg string) checkModelsInfoReaderFn {
			return func(t *testing.T, fetched <-chan pair) {
				data, ok := <-fetched
				assert.True(t, ok)
				assert.ErrorContains(t, data.err, "test-GetModelInfo-error")
			}
		}

		tests = []struct {
			name  string
			data  nextData
			check checkModelsInfoReaderFn
		}{
			{
				name: "success",
				data: nextData{
					model_name: modelPhi4,
					cfg:        getConfigModelsList(nil),
				},
				check: checkModel(modelsList[modelPhi4].model),
			},
			{
				name: "error",
				data: nextData{
					model_name: modelPhi4,
					cfg:        getConfigModelsList(fmt.Errorf("test-GetModelInfo-error")),
				},
				check: checkError("test-GetModelInfo-error"),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				pending = make(chan nextData)
				fetched = modelsInfoFetcher(pending)
			)
			pending <- tt.data
			close(pending)
			tt.check(t, fetched)
		})
	}
}

func Test_modelsInfoFill(t *testing.T) {
	type checkModelsInfoFillFn func(*testing.T, []*ModelItem)

	var (
		fetcher = func(models []string) (chan pair, []*ModelItem) {
			var (
				fetched       = make(chan pair)
				modelItemList = modelsList.getModels(models)
			)

			go func() {
				for _, m := range modelItemList {
					fetched <- pair{name: m.Name, model: m.Model, err: m.Error}
				}
				close(fetched)
			}()

			return fetched, modelItemList
		}

		tests = []struct {
			name   string
			models []string
			check  checkModelsInfoFillFn
		}{
			{
				name:   "success-one-model",
				models: []string{modelPhi4},
			},
			{
				name:   "success-two-models",
				models: []string{modelPhi4, modelLlama3_1},
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fetched, want := fetcher(tt.models)
			got := modelsInfoFill(fetched)
			assert.EqualValues(t, got, want)
		})
	}
}
