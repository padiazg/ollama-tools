## Family fields
At `model_info` there are fields that starts with the family name as we can see in the following examples
**llama3.1**
```json
{
	"details": {
		"family": "llama",
        "parameter_size": "8.0B",
        "quantization_level": "Q4_K_M"
	},
	"model_info": {
        "general.basename": "Meta-Llama-3.1",
		"llama.context_length": 131072,
		"llama.embedding_length": 4096,
	}
}
```
**deepseek-r1:14b**
```json
{
	"details": {
		"family": "qwen2",
        "parameter_size": "14.8B",
        "quantization_level": "Q4_K_M"
	},
	"model_info": {
        "general.basename": "DeepSeek-R1-Distill-Qwen",
		"qwen2.context_length": 131072,
		"qwen2.embedding_length": 5120
	}
}
```
We need to recover the `context_length` and the `embedding_legth` values out of these fields, so I opted for the normalizing approach. 

~~We replace the family name for a common name before umarshaling the string into the models's struct. See `getFamily` for how do we get the family name and `normalizeFamilyFields` for how we replace the family names with a common value.~~

The current approach is to implement a JsonMarshaler `func (m *Model) UnmarshalJSON(raw []byte) error` at `models/ollama/model.go`

