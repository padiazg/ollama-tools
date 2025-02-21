package ollama

import (
	"math"
	"testing"
)

func TestModelDetails_ParameterSizeAsBillions(t *testing.T) {
	var (
		tests = []struct {
			name          string
			parameterSize string
			want          float64
		}{
			{name: "8.0B", parameterSize: "8.0B", want: 8.0},
			{name: "14.8B", parameterSize: "14.8B", want: 14.8},
			{name: "7.6B", parameterSize: "7.6B", want: 7.6},
			{name: "137M", parameterSize: "137M", want: 0.13379},
		}
		round64 = func(value float64, precision uint) float64 {
			ratio := math.Pow(10, float64(precision))
			return math.Round(value*ratio) / ratio
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := TagModelDetails{ParameterSize: tt.parameterSize}
			got := round64(d.ParameterSizeAsBillions(), 5)
			if got != tt.want {
				t.Errorf("ModelDetails.ParameterSizeAsBillions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelDetails_QuantizationLevelAsBitCount(t *testing.T) {
	tests := []struct {
		name              string
		QuantizationLevel string
		want              int
	}{
		{name: "Q2_K", QuantizationLevel: "Q2_K", want: 2},
		{name: "Q3_K", QuantizationLevel: "Q3_K", want: 3},
		{name: "Q3_K_S", QuantizationLevel: "Q3_K_S", want: 3},
		{name: "Q4_K:", QuantizationLevel: "Q4_K:", want: 4},
		{name: "Q4_K_S", QuantizationLevel: "Q4_K_S", want: 4},
		{name: "Q5_K", QuantizationLevel: "Q5_K", want: 5},
		{name: "Q8_0", QuantizationLevel: "Q8_0", want: 8},
		{name: "F16", QuantizationLevel: "F16", want: 16},
		{name: "F32", QuantizationLevel: "F32", want: 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := TagModelDetails{QuantizationLevel: tt.QuantizationLevel}
			if got := d.QuantizationLevelAsBitCount(); got != tt.want {
				t.Errorf("ModelDetails.QuantizationLevelAsBitCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
