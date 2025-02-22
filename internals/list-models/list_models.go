package listmodels

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/padiazg/ollama-tools/internals/tools"
	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

const apiPath = "/api/tags"

// List
func List(cfg *settings.Settings, model_name string) {
	var (
		err  error
		tags *ollama.Tags
	)

	if model_name == "" {
		if tags, err = GetTags(cfg); err != nil {
			fmt.Printf("List getting tags: %v\n", err)
			return
		}

		fmt.Println("Available models:")
		fmt.Println("----------------------------------------------------")

		for _, tag := range tags.Models {
			ListModel(cfg, tag.Name)
		}
	} else {
		ListModel(cfg, model_name)
	}
}

func ListModel(cfg *settings.Settings, model_name string) {
	fmt.Printf("Model: %s\n", model_name)
	model, err := GetModelInfo(cfg, model_name)
	if err != nil {
		fmt.Printf("List getting model info: %+v\n", err)
		return
	}
	fmt.Printf("  Parameters: %s (%d)\n",
		tools.FormatParamCount(model.ModelInfo.ParameterCount),
		model.ModelInfo.ParameterCount)
	fmt.Printf("  Quantization: %s\n", model.Details.QuantizationLevel)
	fmt.Printf("  Context Length: %d tokens\n", model.ModelInfo.ContextLength)
	if model.ModelInfo.EmbeddingLength > 0 {
		fmt.Printf("  Embedding Length: %d\n", model.ModelInfo.EmbeddingLength)
	}

	mem := tools.EstimateMemory(model.ModelInfo.ParameterCount, model.ModelInfo.ContextLength, model.Details.QuantizationLevel)
	tools.PrintEstimatedMemoryPlain(mem)

	if model.ModelInfo.ContextLength > 8192 {
		fmt.Printf("\nNote: This model has a large context length (%d tokens).\n",
			model.ModelInfo.ContextLength)
		fmt.Printf("Reducing max_context in your Ollama request can significantly lower memory usage.\n")
	}
	fmt.Println("")
}

func ListTable(cfg *settings.Settings, model_name string) {
	var (
		err  error
		tags = &ollama.Tags{}
		t    table.Writer
	)

	if model_name == "" {
		if tags, err = GetTags(cfg); err != nil {
			fmt.Printf("List getting tags: %v\n", err)
			return
		}
	} else {
		tags.Models = append(tags.Models, ollama.TagModel{Name: model_name})
	}

	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		table.Row{"Model", "Parameters", "Parameters", "Quantization", "Quantization", "Context Length", "Embedding Length", "Base Model Size", "KV Cache", "GPU RAM", "System RAM"},
		table.RowConfig{AutoMerge: true},
	)
	t.AppendHeader(
		table.Row{"", "Billions", "Units", "level", "bits", "", "", "", "", "", ""},
	)
	// t.SetColumnConfigs([]table.ColumnConfig{
	// 	{Number: 1, AutoMerge: true},
	// 	{Number: 6, AutoMerge: true},
	// 	{Number: 7, AutoMerge: true},
	// 	{Number: 8, AutoMerge: true},
	// 	{Number: 9, AutoMerge: true},
	// 	{Number: 10, AutoMerge: true},
	// 	{Number: 11, AutoMerge: true},
	// })

	t.AppendSeparator()

	for _, tag := range tags.Models {
		model, err := GetModelInfo(cfg, tag.Name)
		if err != nil {
			fmt.Printf("List getting model info: %+v\n", err)
			return
		}

		mem := tools.EstimateMemory(model.ModelInfo.ParameterCount, model.ModelInfo.ContextLength, model.Details.QuantizationLevel)

		t.AppendRow([]interface{}{
			tag.Name,
			text.AlignRight.Apply(fmt.Sprintf("%s", tools.FormatParamCount(model.ModelInfo.ParameterCount)), 8),
			text.AlignRight.Apply(fmt.Sprintf("%d", model.ModelInfo.ParameterCount), 16),
			text.AlignLeft.Apply(model.Details.QuantizationLevel, 8),
			text.AlignRight.Apply(fmt.Sprintf("%d", tools.QuantizationBits(tools.NormalizeQuantizationLevel(model.Details.QuantizationLevel))), 6),
			text.AlignRight.Apply(fmt.Sprintf("%d", model.ModelInfo.ContextLength), 14),
			text.AlignRight.Apply(fmt.Sprintf("%d", model.ModelInfo.EmbeddingLength), 16),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.BaseModelSize), 15),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.KVCacheSize), 10),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.GPURAM), 12),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.SystemRAM), 12),
		})
	}
	t.Render()
}
