package version

import "strconv"

func parseInt(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return i
}

func CurrentVersion() *VersionInfo {
	return v
}
