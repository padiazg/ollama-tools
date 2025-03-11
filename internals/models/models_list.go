package models

import (
	"fmt"

	"github.com/padiazg/ollama-tools/models/ollama"
	"github.com/padiazg/ollama-tools/models/settings"
)

type ModelItem struct {
	Name  string
	Model *ollama.Model
	Error error
}

type pair struct {
	name  string
	model *ollama.Model
	err   error
}

type nextData struct {
	model_name string
	cfg        *settings.Settings
}

type nextFn func() nextData

func ModelsInfoList(cfg *settings.Settings, model_name string) []*ModelItem {
	next := modelsInfoGenerator(cfg, model_name)

	return modelsInfoList(next)
}

func modelsInfoGenerator(cfg *settings.Settings, model_name string) func() nextData {
	var (
		tags  *ollama.Tags
		err   error
		index = 0
	)

	if model_name == "" {
		if tags, err = GetTags(cfg); err != nil {
			fmt.Printf("List getting tags: %v\n", err)
		}
	} else {
		tags = &ollama.Tags{
			Models: []ollama.TagModel{{Name: model_name}},
		}
	}

	return func() nextData {
		if index == len(tags.Models) {
			return nextData{model_name: ""}
		}
		data := nextData{
			model_name: tags.Models[index].Name,
			cfg:        cfg,
		}
		index++
		return data
	}
}

func modelsInfoList(next nextFn) []*ModelItem {
	pending := modelsInfoReader(next)
	fetched := modelsInfoFetcher(pending)

	return modelsInfoFill(fetched)
}

func modelsInfoReader(next nextFn) chan nextData {
	pending := make(chan nextData)
	go func() {
		for {
			nextModel := next()
			if nextModel.model_name == "" {
				close(pending)
				break
			}
			pending <- nextModel
		}
	}()

	return pending
}

func modelsInfoFetcher(pending <-chan nextData) chan pair {
	fetched := make(chan pair)

	go func() {
		for {
			var (
				model *ollama.Model
				err   error
			)

			data, ok := <-pending
			if !ok {
				close(fetched)
				break
			}

			model, err = GetModelInfo(data.cfg, data.model_name)
			fetched <- pair{
				name:  data.model_name,
				model: model,
				err:   err,
			}
		}
	}()

	return fetched
}

func modelsInfoFill(fetched <-chan pair) []*ModelItem {
	list := []*ModelItem{}

	for {
		data, ok := <-fetched
		if !ok {
			break
		}

		list = append(list, &ModelItem{
			Name:  data.name,
			Model: data.model,
			Error: data.err,
		})
	}

	return list
}
