package ollama

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
