package listmodels

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

type testModels struct {
	name               string
	raw                string
	normalized         string
	family             string
	parameter_size     string
	quantization_level string
	parameter_count    int
	context_length     int
	embedding_length   int
}

var (
	models = []testModels{
		{
			name:               "phi4:latest",
			raw:                `{"license":"Microsoft...SOFTWARE.","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|im_start|\u003e\"\nstop                           \"\u003c|im_end|\u003e\"\nstop                           \"\u003c|im_sep|\u003e\"","template":"{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 -}}\n\u003c|im_start|\u003e{{ .Role }}\u003c|im_sep|\u003e\n{{ .Content }}{{ if not $last }}\u003c|im_end|\u003e\n{{ end }}\n{{- if and (ne .Role \"assistant\") $last }}\u003c|im_end|\u003e\n\u003c|im_start|\u003eassistant\u003c|im_sep|\u003e\n{{ end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"phi3","families":["phi3"],"parameter_size":"14.7B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"phi3","general.basename":"phi","general.file_type":15,"general.languages":["en"],"general.license":"mit","general.license.link":"https://huggingface.co/microsoft/phi-4/resolve/main/LICENSE","general.organization":"Microsoft","general.parameter_count":14659507200,"general.quantization_version":2,"general.size_label":"15B","general.tags":["phi","nlp","math","code","chat","conversational","text-generation"],"general.type":"model","general.version":"4","phi3.attention.head_count":40,"phi3.attention.head_count_kv":10,"phi3.attention.layer_norm_rms_epsilon":0.00001,"phi3.attention.sliding_window":131072,"phi3.block_count":40,"phi3.context_length":16384,"phi3.embedding_length":5120,"phi3.feed_forward_length":17920,"phi3.rope.dimension_count":128,"phi3.rope.freq_base":250000,"phi3.rope.scaling.original_context_length":16384,"tokenizer.ggml.bos_token_id":100257,"tokenizer.ggml.eos_token_id":100257,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.padding_token_id":100257,"tokenizer.ggml.pre":"dbrx","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-01-14T17:21:17.785607967-03:00"}`,
			normalized:         `{"license":"Microsoft...SOFTWARE.","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|im_start|\u003e\"\nstop                           \"\u003c|im_end|\u003e\"\nstop                           \"\u003c|im_sep|\u003e\"","template":"{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 -}}\n\u003c|im_start|\u003e{{ .Role }}\u003c|im_sep|\u003e\n{{ .Content }}{{ if not $last }}\u003c|im_end|\u003e\n{{ end }}\n{{- if and (ne .Role \"assistant\") $last }}\u003c|im_end|\u003e\n\u003c|im_start|\u003eassistant\u003c|im_sep|\u003e\n{{ end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"phi3","families":["phi3"],"parameter_size":"14.7B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"phi3","general.basename":"phi","general.file_type":15,"general.languages":["en"],"general.license":"mit","general.license.link":"https://huggingface.co/microsoft/phi-4/resolve/main/LICENSE","general.organization":"Microsoft","general.parameter_count":14659507200,"general.quantization_version":2,"general.size_label":"15B","general.tags":["phi","nlp","math","code","chat","conversational","text-generation"],"general.type":"model","general.version":"4","phi3.attention.head_count":40,"phi3.attention.head_count_kv":10,"phi3.attention.layer_norm_rms_epsilon":0.00001,"phi3.attention.sliding_window":131072,"phi3.block_count":40,"model.context_length":16384,"model.embedding_length":5120,"phi3.feed_forward_length":17920,"phi3.rope.dimension_count":128,"phi3.rope.freq_base":250000,"phi3.rope.scaling.original_context_length":16384,"tokenizer.ggml.bos_token_id":100257,"tokenizer.ggml.eos_token_id":100257,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.padding_token_id":100257,"tokenizer.ggml.pre":"dbrx","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-01-14T17:21:17.785607967-03:00"}`,
			family:             "phi3",
			context_length:     16384,
			embedding_length:   5120,
			parameter_count:    14659507200,
			parameter_size:     "14.7B",
			quantization_level: "Q4_K_M",
		},
		{
			name:               "llama3.1:latest",
			raw:                `{"license":"LLAMA 3.1 ...","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|start_header_id|\u003e\"\nstop                           \"\u003c|end_header_id|\u003e\"\nstop                           \"\u003c|eot_id|\u003e\"","template":"{{- if or .System .Tools }}\u003c|start_header_id|\u003esystem\u003c|end_header_id|\u003e\n{{- if .System }}\n\n{{ .System }}\n{{- end }}\n{{- if .Tools }}\n\nCutting Knowledge Date: December 2023\n\nWhen you receive a tool call response, use the output to format an answer to the orginal user question.\n\nYou are a helpful assistant with tool calling capabilities.\n{{- end }}\u003c|eot_id|\u003e\n{{- end }}\n{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 }}\n{{- if eq .Role \"user\" }}\u003c|start_header_id|\u003euser\u003c|end_header_id|\u003e\n{{- if and $.Tools $last }}\n\nGiven the following functions, please respond with a JSON for a function call with its proper arguments that best answers the given prompt.\n\nRespond in the format {\"name\": function name, \"parameters\": dictionary of argument name and its value}. Do not use variables.\n\n{{ range $.Tools }}\n{{- . }}\n{{ end }}\nQuestion: {{ .Content }}\u003c|eot_id|\u003e\n{{- else }}\n\n{{ .Content }}\u003c|eot_id|\u003e\n{{- end }}{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- else if eq .Role \"assistant\" }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n{{- if .ToolCalls }}\n{{ range .ToolCalls }}\n{\"name\": \"{{ .Function.Name }}\", \"parameters\": {{ .Function.Arguments }}}{{ end }}\n{{- else }}\n\n{{ .Content }}\n{{- end }}{{ if not $last }}\u003c|eot_id|\u003e{{ end }}\n{{- else if eq .Role \"tool\" }}\u003c|start_header_id|\u003eipython\u003c|end_header_id|\u003e\n\n{{ .Content }}\u003c|eot_id|\u003e{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"llama","families":["llama"],"parameter_size":"8.0B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"llama","general.basename":"Meta-Llama-3.1","general.file_type":15,"general.finetune":"Instruct","general.languages":["en","de","fr","it","pt","hi","es","th"],"general.license":"llama3.1","general.parameter_count":8030261312,"general.quantization_version":2,"general.size_label":"8B","general.tags":["facebook","meta","pytorch","llama","llama-3","text-generation"],"general.type":"model","llama.attention.head_count":32,"llama.attention.head_count_kv":8,"llama.attention.layer_norm_rms_epsilon":0.00001,"llama.block_count":32,"llama.context_length":131072,"llama.embedding_length":4096,"llama.feed_forward_length":14336,"llama.rope.dimension_count":128,"llama.rope.freq_base":500000,"llama.vocab_size":128256,"tokenizer.ggml.bos_token_id":128000,"tokenizer.ggml.eos_token_id":128009,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.pre":"llama-bpe","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-02-03T19:22:17.054410969-03:00"}`,
			normalized:         `{"license":"LLAMA 3.1 ...","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|start_header_id|\u003e\"\nstop                           \"\u003c|end_header_id|\u003e\"\nstop                           \"\u003c|eot_id|\u003e\"","template":"{{- if or .System .Tools }}\u003c|start_header_id|\u003esystem\u003c|end_header_id|\u003e\n{{- if .System }}\n\n{{ .System }}\n{{- end }}\n{{- if .Tools }}\n\nCutting Knowledge Date: December 2023\n\nWhen you receive a tool call response, use the output to format an answer to the orginal user question.\n\nYou are a helpful assistant with tool calling capabilities.\n{{- end }}\u003c|eot_id|\u003e\n{{- end }}\n{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 }}\n{{- if eq .Role \"user\" }}\u003c|start_header_id|\u003euser\u003c|end_header_id|\u003e\n{{- if and $.Tools $last }}\n\nGiven the following functions, please respond with a JSON for a function call with its proper arguments that best answers the given prompt.\n\nRespond in the format {\"name\": function name, \"parameters\": dictionary of argument name and its value}. Do not use variables.\n\n{{ range $.Tools }}\n{{- . }}\n{{ end }}\nQuestion: {{ .Content }}\u003c|eot_id|\u003e\n{{- else }}\n\n{{ .Content }}\u003c|eot_id|\u003e\n{{- end }}{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- else if eq .Role \"assistant\" }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n{{- if .ToolCalls }}\n{{ range .ToolCalls }}\n{\"name\": \"{{ .Function.Name }}\", \"parameters\": {{ .Function.Arguments }}}{{ end }}\n{{- else }}\n\n{{ .Content }}\n{{- end }}{{ if not $last }}\u003c|eot_id|\u003e{{ end }}\n{{- else if eq .Role \"tool\" }}\u003c|start_header_id|\u003eipython\u003c|end_header_id|\u003e\n\n{{ .Content }}\u003c|eot_id|\u003e{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","family":"llama","families":["llama"],"parameter_size":"8.0B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"llama","general.basename":"Meta-Llama-3.1","general.file_type":15,"general.finetune":"Instruct","general.languages":["en","de","fr","it","pt","hi","es","th"],"general.license":"llama3.1","general.parameter_count":8030261312,"general.quantization_version":2,"general.size_label":"8B","general.tags":["facebook","meta","pytorch","llama","llama-3","text-generation"],"general.type":"model","llama.attention.head_count":32,"llama.attention.head_count_kv":8,"llama.attention.layer_norm_rms_epsilon":0.00001,"llama.block_count":32,"model.context_length":131072,"model.embedding_length":4096,"llama.feed_forward_length":14336,"llama.rope.dimension_count":128,"llama.rope.freq_base":500000,"llama.vocab_size":128256,"tokenizer.ggml.bos_token_id":128000,"tokenizer.ggml.eos_token_id":128009,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.pre":"llama-bpe","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-02-03T19:22:17.054410969-03:00"}`,
			family:             "llama",
			context_length:     131072,
			embedding_length:   4096,
			parameter_count:    8030261312,
			parameter_size:     "8.0B",
			quantization_level: "Q4_K_M",
		},
		{
			name:               "nomic-embed-text:latest",
			raw:                `{"license":"Apache...the License.\n","modelfile":"# Modelfile ...","parameters":"num_ctx                        8192","template":"{{ .Prompt }}","details":{"parent_model":"","format":"gguf","family":"nomic-bert","families":["nomic-bert"],"parameter_size":"137M","quantization_level":"F16"},"model_info":{"general.architecture":"nomic-bert","general.file_type":1,"general.parameter_count":136727040,"nomic-bert.attention.causal":false,"nomic-bert.attention.head_count":12,"nomic-bert.attention.layer_norm_epsilon":1e-12,"nomic-bert.block_count":12,"nomic-bert.context_length":2048,"nomic-bert.embedding_length":768,"nomic-bert.feed_forward_length":3072,"nomic-bert.pooling_type":1,"nomic-bert.rope.freq_base":1000,"tokenizer.ggml.bos_token_id":101,"tokenizer.ggml.cls_token_id":101,"tokenizer.ggml.eos_token_id":102,"tokenizer.ggml.mask_token_id":103,"tokenizer.ggml.model":"bert","tokenizer.ggml.padding_token_id":0,"tokenizer.ggml.scores":null,"tokenizer.ggml.seperator_token_id":102,"tokenizer.ggml.token_type":null,"tokenizer.ggml.token_type_count":2,"tokenizer.ggml.tokens":null,"tokenizer.ggml.unknown_token_id":100},"modified_at":"2025-02-03T19:22:18.145435125-03:00"}`,
			normalized:         `{"license":"Apache...the License.\n","modelfile":"# Modelfile ...","parameters":"num_ctx                        8192","template":"{{ .Prompt }}","details":{"parent_model":"","format":"gguf","family":"nomic-bert","families":["nomic-bert"],"parameter_size":"137M","quantization_level":"F16"},"model_info":{"general.architecture":"nomic-bert","general.file_type":1,"general.parameter_count":136727040,"nomic-bert.attention.causal":false,"nomic-bert.attention.head_count":12,"nomic-bert.attention.layer_norm_epsilon":1e-12,"nomic-bert.block_count":12,"model.context_length":2048,"model.embedding_length":768,"nomic-bert.feed_forward_length":3072,"nomic-bert.pooling_type":1,"nomic-bert.rope.freq_base":1000,"tokenizer.ggml.bos_token_id":101,"tokenizer.ggml.cls_token_id":101,"tokenizer.ggml.eos_token_id":102,"tokenizer.ggml.mask_token_id":103,"tokenizer.ggml.model":"bert","tokenizer.ggml.padding_token_id":0,"tokenizer.ggml.scores":null,"tokenizer.ggml.seperator_token_id":102,"tokenizer.ggml.token_type":null,"tokenizer.ggml.token_type_count":2,"tokenizer.ggml.tokens":null,"tokenizer.ggml.unknown_token_id":100},"modified_at":"2025-02-03T19:22:18.145435125-03:00"}`,
			family:             "nomic-bert",
			context_length:     2048,
			embedding_length:   768,
			parameter_count:    136727040,
			parameter_size:     "137M",
			quantization_level: "F16",
		},
		{
			name:   "no-family",
			raw:    `{"license":"LLAMA 3.1 ...","modelfile":"# Modelfile ...","parameters":"stop                           \"\u003c|start_header_id|\u003e\"\nstop                           \"\u003c|end_header_id|\u003e\"\nstop                           \"\u003c|eot_id|\u003e\"","template":"{{- if or .System .Tools }}\u003c|start_header_id|\u003esystem\u003c|end_header_id|\u003e\n{{- if .System }}\n\n{{ .System }}\n{{- end }}\n{{- if .Tools }}\n\nCutting Knowledge Date: December 2023\n\nWhen you receive a tool call response, use the output to format an answer to the orginal user question.\n\nYou are a helpful assistant with tool calling capabilities.\n{{- end }}\u003c|eot_id|\u003e\n{{- end }}\n{{- range $i, $_ := .Messages }}\n{{- $last := eq (len (slice $.Messages $i)) 1 }}\n{{- if eq .Role \"user\" }}\u003c|start_header_id|\u003euser\u003c|end_header_id|\u003e\n{{- if and $.Tools $last }}\n\nGiven the following functions, please respond with a JSON for a function call with its proper arguments that best answers the given prompt.\n\nRespond in the format {\"name\": function name, \"parameters\": dictionary of argument name and its value}. Do not use variables.\n\n{{ range $.Tools }}\n{{- . }}\n{{ end }}\nQuestion: {{ .Content }}\u003c|eot_id|\u003e\n{{- else }}\n\n{{ .Content }}\u003c|eot_id|\u003e\n{{- end }}{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- else if eq .Role \"assistant\" }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n{{- if .ToolCalls }}\n{{ range .ToolCalls }}\n{\"name\": \"{{ .Function.Name }}\", \"parameters\": {{ .Function.Arguments }}}{{ end }}\n{{- else }}\n\n{{ .Content }}\n{{- end }}{{ if not $last }}\u003c|eot_id|\u003e{{ end }}\n{{- else if eq .Role \"tool\" }}\u003c|start_header_id|\u003eipython\u003c|end_header_id|\u003e\n\n{{ .Content }}\u003c|eot_id|\u003e{{ if $last }}\u003c|start_header_id|\u003eassistant\u003c|end_header_id|\u003e\n\n{{ end }}\n{{- end }}\n{{- end }}","details":{"parent_model":"","format":"gguf","_family_":"llama","families":["llama"],"parameter_size":"8.0B","quantization_level":"Q4_K_M"},"model_info":{"general.architecture":"llama","general.basename":"Meta-Llama-3.1","general.file_type":15,"general.finetune":"Instruct","general.languages":["en","de","fr","it","pt","hi","es","th"],"general.license":"llama3.1","general.parameter_count":8030261312,"general.quantization_version":2,"general.size_label":"8B","general.tags":["facebook","meta","pytorch","llama","llama-3","text-generation"],"general.type":"model","llama.attention.head_count":32,"llama.attention.head_count_kv":8,"llama.attention.layer_norm_rms_epsilon":0.00001,"llama.block_count":32,"llama.context_length":131072,"llama.embedding_length":4096,"llama.feed_forward_length":14336,"llama.rope.dimension_count":128,"llama.rope.freq_base":500000,"llama.vocab_size":128256,"tokenizer.ggml.bos_token_id":128000,"tokenizer.ggml.eos_token_id":128009,"tokenizer.ggml.merges":null,"tokenizer.ggml.model":"gpt2","tokenizer.ggml.pre":"llama-bpe","tokenizer.ggml.token_type":null,"tokenizer.ggml.tokens":null},"modified_at":"2025-02-03T19:22:17.054410969-03:00"}`,
			family: "",
		},
	}
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type CheckClientInfoFn func(t *testing.T, m *ollama.Model)

var CheckClientInfo = func(fns ...CheckClientInfoFn) []CheckClientInfoFn { return fns }

func TestGetModelInfo(t *testing.T) {
	type args struct {
		cfg       *settings.Settings
		modelName string
	}
	tests := []struct {
		name    string
		args    args
		want    *ollama.Model
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModelInfo(tt.args.cfg, tt.args.modelName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModelInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModelInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFamily(t *testing.T) {
	var (
		tests = []struct {
			name    string
			raw     string
			want    string
			wantErr bool
		}{
			{
				name:    models[0].name,
				raw:     models[0].raw,
				want:    models[0].family,
				wantErr: false,
			},
			{
				name:    models[1].name,
				raw:     models[1].raw,
				want:    models[1].family,
				wantErr: false,
			},
			{
				name:    models[2].name,
				raw:     models[2].raw,
				want:    models[2].family,
				wantErr: false,
			},
			{
				name:    models[3].name,
				raw:     models[3].raw,
				wantErr: true,
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFamily(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFamily() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && got != tt.want {
				t.Errorf("getFamily() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeFamilyFields(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{
			name:    models[0].name,
			raw:     models[0].raw,
			want:    models[0].normalized,
			wantErr: false,
		},
		{
			name:    models[1].name,
			raw:     models[1].raw,
			want:    models[1].normalized,
			wantErr: false,
		},
		{
			name:    models[2].name,
			raw:     models[2].raw,
			want:    models[2].normalized,
			wantErr: false,
		},
		{
			name:    models[3].name,
			raw:     models[3].raw,
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeFamilyFields(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeFamilyFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeFamilyFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
