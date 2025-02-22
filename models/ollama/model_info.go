package ollama

type ModelInfo struct {
	Type            string `json:"general.type"`
	ParameterCount  int64  `json:"general.parameter_count"`
	ContextLength   int    `json:"model.context_length"`
	EmbeddingLength int    `json:"model.embedding_length"`
}
