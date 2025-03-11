package settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type Settings struct {
	OllamaUrl string `json:"ollamaurl"`
	Transport http.RoundTripper
}

func (s *Settings) Show() {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Settings:\n%s\n", string(b))
}

func (s *Settings) ShowKeyValuePairs() {
	fmt.Printf("Key/Value pairs:\n")
	// print all the keys
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		fmt.Printf("  %s: %v\n", key, val)
	}
}

func (s *Settings) Save(name string) error {
	if err := viper.WriteConfigAs(name); err != nil {
		return fmt.Errorf("writing config file: %+v", err)
	}

	return nil
}
