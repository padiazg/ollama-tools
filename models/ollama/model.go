package ollama

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Model struct {
	Details   ModelDetails `json:"details"`
	ModelInfo ModelInfo    `json:"model_info"`
}

// UnmarshalJSON will try to normalize field names before
// unmarshaling
func (m *Model) UnmarshalJSON(raw []byte) error {
	var (
		data []byte
		err  error
	)

	if data, err = replaceFamilyFields(raw); err != nil {
		return fmt.Errorf("normalizig fields: %+v", err)
	}

	type Alias Model
	temp := (*Alias)(m)

	return json.Unmarshal([]byte(data), &temp)
}

// replaceFamilyFields will raplace the family name with a plain `model`
// at the beggining of some fields
func replaceFamilyFields(raw []byte) ([]byte, error) {
	familySearch := regexp.
		MustCompile(`(?U)"family":\s?"(.*)",`).
		FindAllStringSubmatch(string(raw), -1)

	if len(familySearch) < 1 {
		return nil, fmt.Errorf("no family found")
	}

	family := familySearch[0][1]
	data := regexp.
		MustCompile(fmt.Sprintf(`(?m)(?U)%s\.(context_length|embedding_length)`, family)).
		ReplaceAllString(string(raw), "model.$1")

	return []byte(data), nil
}
