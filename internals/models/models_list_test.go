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

type testTags struct {
	body   string
	models []string
}

type testModelsInfo struct {
	name       string
	normalized string
	model      *ollama.Model
}

const (
	modelPhi4     string = "phi4:latest"
	modelLlama3_1 string = "llama3.1:latest"
)

var (
	tagsList = map[string]testTags{
		"one-model": {
			body:   `{"models":[{"name":"deepseek-r1:7b","model":"deepseek-r1:7b","modified_at":"2025-03-11T15:37:14.423620816-03:00","size":4683075271,"digest":"0a8c266910232fd3291e71e5ba1e058cc5af9d411192cf88b6d30e92b6e73163","details":{"parent_model":"","format":"gguf","family":"qwen2","families":["qwen2"],"parameter_size":"7.6B","quantization_level":"Q4_K_M"}}]}`,
			models: []string{"deepseek-r1:7b"},
		},
		"two-models": {
			body:   `{"models":[{"name":"deepseek-r1:7b","model":"deepseek-r1:7b","modified_at":"2025-03-11T15:37:14.423620816-03:00","size":4683075271,"digest":"0a8c266910232fd3291e71e5ba1e058cc5af9d411192cf88b6d30e92b6e73163","details":{"parent_model":"","format":"gguf","family":"qwen2","families":["qwen2"],"parameter_size":"7.6B","quantization_level":"Q4_K_M"}},{"name":"qwen2.5-coder:7b","model":"qwen2.5-coder:7b","modified_at":"2025-02-26T13:02:11.408112713-03:00","size":4683087519,"digest":"2b0496514337a3d5901f1d253d01726c890b721e891335a56d6e08cedf3e2cb0","details":{"parent_model":"","format":"gguf","family":"qwen2","families":["qwen2"],"parameter_size":"7.6B","quantization_level":"Q4_K_M"}}]}`,
			models: []string{"deepseek-r1:7b", "qwen2.5-coder:7b"},
		},
	}

	modelsList = map[string]testModelsInfo{
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
		// "nomic-embed-text": {
		// 	name:               "nomic-embed-text:latest",
		// 	raw:                `{"license":"Apache...the License.\n","modelfile":"# Modelfile ...","parameters":"num_ctx                        8192","template":"{{ .Prompt }}","details":{"parent_model":"","format":"gguf","family":"nomic-bert","families":["nomic-bert"],"parameter_size":"137M","quantization_level":"F16"},"model_info":{"general.architecture":"nomic-bert","general.file_type":1,"general.parameter_count":136727040,"nomic-bert.attention.causal":false,"nomic-bert.attention.head_count":12,"nomic-bert.attention.layer_norm_epsilon":1e-12,"nomic-bert.block_count":12,"nomic-bert.context_length":2048,"nomic-bert.embedding_length":768,"nomic-bert.feed_forward_length":3072,"nomic-bert.pooling_type":1,"nomic-bert.rope.freq_base":1000,"tokenizer.ggml.bos_token_id":101,"tokenizer.ggml.cls_token_id":101,"tokenizer.ggml.eos_token_id":102,"tokenizer.ggml.mask_token_id":103,"tokenizer.ggml.model":"bert","tokenizer.ggml.padding_token_id":0,"tokenizer.ggml.scores":null,"tokenizer.ggml.seperator_token_id":102,"tokenizer.ggml.token_type":null,"tokenizer.ggml.token_type_count":2,"tokenizer.ggml.tokens":null,"tokenizer.ggml.unknown_token_id":100},"modified_at":"2025-02-03T19:22:18.145435125-03:00"}`,
		// 	normalized:         `{"license":"Apache...the License.\n","modelfile":"# Modelfile ...","parameters":"num_ctx                        8192","template":"{{ .Prompt }}","details":{"parent_model":"","format":"gguf","family":"nomic-bert","families":["nomic-bert"],"parameter_size":"137M","quantization_level":"F16"},"model_info":{"general.architecture":"nomic-bert","general.file_type":1,"general.parameter_count":136727040,"nomic-bert.attention.causal":false,"nomic-bert.attention.head_count":12,"nomic-bert.attention.layer_norm_epsilon":1e-12,"nomic-bert.block_count":12,"model.context_length":2048,"model.embedding_length":768,"nomic-bert.feed_forward_length":3072,"nomic-bert.pooling_type":1,"nomic-bert.rope.freq_base":1000,"tokenizer.ggml.bos_token_id":101,"tokenizer.ggml.cls_token_id":101,"tokenizer.ggml.eos_token_id":102,"tokenizer.ggml.mask_token_id":103,"tokenizer.ggml.model":"bert","tokenizer.ggml.padding_token_id":0,"tokenizer.ggml.scores":null,"tokenizer.ggml.seperator_token_id":102,"tokenizer.ggml.token_type":null,"tokenizer.ggml.token_type_count":2,"tokenizer.ggml.tokens":null,"tokenizer.ggml.unknown_token_id":100},"modified_at":"2025-02-03T19:22:18.145435125-03:00"}`,
		// 	family:             "nomic-bert",
		// 	context_length:     2048,
		// 	embedding_length:   768,
		// 	parameter_count:    136727040,
		// 	parameter_size:     "137M",
		// 	quantization_level: "F16",
		// },
	}
)

type DryRunTransport struct {
	http.RoundTripper
	RoundTripFn func(r *http.Request) (*http.Response, error)
}

func (dr *DryRunTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return dr.RoundTripFn(r)
}

// func TestModelsInfoList(t *testing.T) {
// 	type args struct {
// 		cfg        *settings.Settings
// 		model_name string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []*ModelItem
// 	}{
// 		{
// 			name: "",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ModelsInfoList(tt.args.cfg, tt.args.model_name); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ModelsInfoList() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
			name       string
			model_name string
			cfg        *settings.Settings
			check      func(*testing.T, nextFn)
			wantError  bool
		}{
			{
				name: "one-model",
				cfg: &settings.Settings{
					OllamaUrl: "http://ollama:11434",
					Transport: &DryRunTransport{
						RoundTripFn: func(r *http.Request) (*http.Response, error) {
							var (
								res *http.Response
								err error
							)

							res = &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(strings.NewReader(tagsList["one-model"].body)),
							}

							return res, err
						},
					},
				},
				check: checkModels(tagsList["one-model"].models),
			},
			{
				name: "two-models",
				cfg: &settings.Settings{
					OllamaUrl: "http://ollama:11434",
					Transport: &DryRunTransport{
						RoundTripFn: func(r *http.Request) (*http.Response, error) {
							var (
								res *http.Response
								err error
							)

							res = &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(strings.NewReader(tagsList["two-models"].body)),
							}

							return res, err
						},
					},
				},
				check: checkModels(tagsList["two-models"].models),
			},
			{
				name: "specific-models",
				cfg: &settings.Settings{
					OllamaUrl: "http://ollama:11434",
				},
				model_name: "fake-model",
				check:      checkModels([]string{"fake-model"}),
			},
			{
				name: "http-client-error",
				cfg: &settings.Settings{
					OllamaUrl: "http://ollama:11434",
					Transport: &DryRunTransport{
						RoundTripFn: func(r *http.Request) (*http.Response, error) {
							return nil, fmt.Errorf("from http-client-mock")
						},
					},
				},
				wantError: true,
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			next, err := modelsInfoGenerator(tt.cfg, tt.model_name)
			if tt.wantError {
				assert.Error(t, err, "error expected, not found")
			} else {
				if assert.NoError(t, err, "error not expected: %+v", err) {
					tt.check(t, next)
				}
			}
		})
	}
}

// func Test_modelsInfoList(t *testing.T) {
// 	type args struct {
// 		next nextFn
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []*ModelItem
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := modelsInfoList(tt.args.next); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("modelsInfoList() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
		modelListRoundtripper = func(wantErr error) func(r *http.Request) (*http.Response, error) {
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

		checkModels = func(model *ModelItem) checkModelsInfoReaderFn {
			return func(t *testing.T, fetched <-chan pair) {
				data, ok := <-fetched
				assert.True(t, ok)
				assert.EqualValues(t, model, data.model)
				// assert.NotNil(t, data.model)
				// list := modelsInfoFill(fetched)
				// assert.ElementsMatch(t, models, list)
			}
		}

		checkError = func(errorMsg string) checkModelsInfoReaderFn {
			return func(t *testing.T, fetched <-chan pair) {
				list := modelsInfoFill(fetched)
				assert.Len(t, list, 1)
				assert.ErrorContains(t, list[0].Error, "test-GetModelInfo-error")
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
					cfg: &settings.Settings{
						OllamaUrl: "http://localhost:11434",
						Transport: &DryRunTransport{
							RoundTripFn: modelListRoundtripper(nil),
						},
					},
				},
				// check: checkModels([]*ModelItem{{
				// 	Name:  modelPhi4,
				// 	Model: modelsList[modelPhi4].model,
				// }}),
				check: checkModels(&ModelItem{
					Name:  modelPhi4,
					Model: modelsList[modelPhi4].model,
				}),
			},
			{
				name: "error",
				data: nextData{
					model_name: modelPhi4,
					cfg: &settings.Settings{
						OllamaUrl: "http://localhost:11434",
						Transport: &DryRunTransport{
							RoundTripFn: modelListRoundtripper(fmt.Errorf("test-GetModelInfo-error")),
						},
					},
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

// func Test_modelsInfoFill(t *testing.T) {
// 	type args struct {
// 		fetched <-chan pair
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []*ModelItem
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := modelsInfoFill(tt.args.fetched); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("modelsInfoFill() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
