package ollama

import (
	"fmt"
	"regexp"
)

type ShowModel struct {
	raw string
}

func (m ShowModel) getFamily() string {
	var (
		re  = regexp.MustCompile(`(?U)"family":\s?"(.*)",`)
		res = re.FindAllStringSubmatch(m.raw, -1)
	)

	if len(res) < 1 {
		fmt.Printf("getFamily can't get family\n")
		return ""
	}

	return res[0][1]
}

func (m ShowModel) normalizeFamilyFields() {
	var (
		family = m.getFamily()
	)

}
