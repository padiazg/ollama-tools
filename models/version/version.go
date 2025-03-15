package version

import (
	"fmt"
	"regexp"
	"time"
)

var (
	version   = "v0.0.1"
	commit    = "unknown"
	buildDate = "unknown"
	builtBy   = "unknown"
	v         = &VersionInfo{}
)

type VersionInfo struct {
	Version   string
	Commit    string
	BuildDate string
	BuiltBy   string
	Major     int
	Minor     int
	Patch     int
	Extra     string
	TimeStamp *time.Time
}

func (v *VersionInfo) ParseVersion() {
	re := regexp.MustCompile(`(?m)v?(\d{1,3})\.(\d{1,3})\.?(\d{1,3})?-?(.*)?`)

	if re.MatchString(v.Version) {
		matches := re.FindStringSubmatch(v.Version)
		v.Major = parseInt(matches[1])
		v.Minor = parseInt(matches[2])
		v.Patch = parseInt(matches[3])
		v.Extra = matches[4]
	}
}

func (v *VersionInfo) ParseDate() error {
	ts, err := time.Parse(time.RFC3339, v.BuildDate)
	if err != nil {
		return err
	}

	v.TimeStamp = &ts

	return nil
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("%s %s %s", v.Version, v.Commit, v.BuildDate)
}

func init() {
	v.Version = version
	v.Commit = commit
	v.BuildDate = buildDate
	v.BuiltBy = builtBy
	_ = v.ParseDate()
	v.ParseVersion()
}
