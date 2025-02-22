package ollama

type Model struct {
	Details   ModelDetails `json:"details"`
	ModelInfo ModelInfo    `json:"model_info"`
}
