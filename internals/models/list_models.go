package models

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/padiazg/ollama-tools/internals/tools"
	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

// List
func List(cfg *settings.Settings, model_name string, table bool) {
	models, err := ModelsInfoList(cfg, model_name)
	if err != nil {
		fmt.Printf("listing models: %+v", err)
	}

	if table {
		listModelsTable(models)
	} else {
		listModelsDetail(models)
	}
}

func listModelsDetail(models []*ModelItem) {
	fmt.Println("Available models:")
	fmt.Println("----------------------------------------------------")
	for _, model := range models {
		printModel(model)
	}
}

func printModel(model *ModelItem) {
	modelInfo := model.Model.ModelInfo
	details := model.Model.Details

	fmt.Printf("Model: %s\n", model.Name)
	fmt.Printf("  Parameters: %s (%d)\n",
		tools.FormatParamCount(modelInfo.ParameterCount),
		modelInfo.ParameterCount)
	fmt.Printf("  Quantization: %s\n", details.QuantizationLevel)
	fmt.Printf("  Context Length: %d tokens\n", modelInfo.ContextLength)
	if modelInfo.EmbeddingLength > 0 {
		fmt.Printf("  Embedding Length: %d\n", modelInfo.EmbeddingLength)
	}

	mem := tools.EstimateMemory(modelInfo.ParameterCount, modelInfo.ContextLength, details.QuantizationLevel)
	tools.PrintEstimatedMemoryPlain(mem)

	if modelInfo.ContextLength > 8192 {
		fmt.Printf("\nNote: This model has a large context length (%d tokens).\n",
			modelInfo.ContextLength)
		fmt.Printf("Reducing max_context in your Ollama request can significantly lower memory usage.\n")
	}
	fmt.Println("")
}

func listModelsTable(models []*ModelItem) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		table.Row{"Model", "Parameters", "Parameters", "Quantization", "Quantization", "Context Length", "Embedding Length", "Base Model Size", "KV Cache", "GPU RAM", "System RAM"},
		table.RowConfig{AutoMerge: true},
	)
	t.AppendHeader(
		table.Row{"", "Billions", "Units", "level", "bits", "", "", "", "", "", ""},
	)

	t.AppendSeparator()

	for _, model := range models {
		modelInfo := model.Model.ModelInfo
		details := model.Model.Details
		mem := tools.EstimateMemory(modelInfo.ParameterCount, modelInfo.ContextLength, details.QuantizationLevel)

		t.AppendRow([]interface{}{
			model.Name,
			text.AlignRight.Apply(tools.FormatParamCount(modelInfo.ParameterCount), 8),
			text.AlignRight.Apply(fmt.Sprintf("%d", modelInfo.ParameterCount), 16),
			text.AlignLeft.Apply(details.QuantizationLevel, 8),
			text.AlignRight.Apply(fmt.Sprintf("%d", tools.QuantizationBits(tools.NormalizeQuantizationLevel(details.QuantizationLevel))), 6),
			text.AlignRight.Apply(fmt.Sprintf("%d", modelInfo.ContextLength), 14),
			text.AlignRight.Apply(fmt.Sprintf("%d", modelInfo.EmbeddingLength), 16),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.BaseModelSize), 15),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.KVCacheSize), 10),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.GPURAM), 12),
			text.AlignRight.Apply(fmt.Sprintf("%.2f Gb", mem.SystemRAM), 12),
		})
	}

	t.Render()
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
			text.AlignRight.Apply(tools.FormatParamCount(model.ModelInfo.ParameterCount), 8),
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
