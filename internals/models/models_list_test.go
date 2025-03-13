package models

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/padiazg/ollama-tools/models/settings"
	"github.com/stretchr/testify/assert"
)

type testTags struct {
	body   string
	models []string
}

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
	type checkModelsInfoReader func(t *testing.T, pending <-chan nextData)

	var (
		checkModelName = func(model_name string) checkModelsInfoReader {
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

		checkEndReading = func() checkModelsInfoReader {
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
			check      checkModelsInfoReader
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
	//	type args struct {
	//		pending <-chan nextData
	//	}
	//
	//	tests := []struct {
	//		name string
	//		args args
	//		want chan pair
	//	}{
	//
	//		// TODO: Add test cases.
	//	}
	//
	//	for _, tt := range tests {
	//		t.Run(tt.name, func(t *testing.T) {
	//			if got := modelsInfoFetcher(tt.args.pending); !reflect.DeepEqual(got, tt.want) {
	//				t.Errorf("modelsInfoFetcher() = %v, want %v", got, tt.want)
	//			}
	//		})
	//	}
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
