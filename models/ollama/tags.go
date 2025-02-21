package ollama

import (
	"fmt"
	"regexp"
	"strconv"
)

type TagModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type TagModel struct {
	Name       string          `json:"name"`
	Model      string          `json:"model"`
	ModifiedAt string          `json:"modified_at"`
	Size       int             `json:"size"`
	Digest     string          `json:"digest"`
	Details    TagModelDetails `json:"details"`
}

type Tags struct {
	Models []TagModel `json:"models"`
}

func (d TagModelDetails) ParameterSizeAsBillions() float64 {
	var (
		re            = regexp.MustCompile(`(?m)([0-9]+(\.[0-9])?)([BM])`)
		res           = re.FindAllStringSubmatch(d.ParameterSize, -1)
		parameterSize = 0.0
	)

	if len(res) < 1 {
		fmt.Printf("parseParameterSize error parsing: %s\n", d.ParameterSize)
		return 0.0
	}

	// extract size
	value, err := strconv.ParseFloat(res[0][1], 32)
	if err != nil {
		fmt.Printf("parseParameterSize error converting value: %s\n", res[0][1])
	}

	// scale if needed
	switch res[0][3] {
	case "B":
		parameterSize = value
	case "M":
		parameterSize = value / 1024.0
	}

	return parameterSize
}

func (d TagModelDetails) QuantizationLevelAsBitCount() int {
	var (
		quantizationLevel = 0
		err               error
		prefix            = d.QuantizationLevel[0]
	)

	switch prefix {
	case 'Q':
		if quantizationLevel, err = strconv.Atoi(string(d.QuantizationLevel[1])); err != nil {
			return 0
		}
	case 'F':
		if quantizationLevel, err = strconv.Atoi(string(d.QuantizationLevel[1:3])); err != nil {
			return 0
		}
	}

	return quantizationLevel
}

// func (d TagModelDetails) EstimatedRequiredMemory() float64 {
// 	var (
// 		memory            float64
// 		parameterSize     = d.ParameterSizeAsBillions()
// 		quantizationLevel = d.QuantizationLevelAsBitCount()
// 	)

// 	return memory
// }
