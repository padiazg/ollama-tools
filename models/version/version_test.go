package version

import (
	"testing"
	"time"
)

func TestVersion_ParseVersion(t *testing.T) {
	type Expect struct {
		Major int
		Minor int
		Patch int
		Extra string
	}
	tests := []struct {
		name    string
		version string
		want    Expect
	}{
		{
			name:    "v0.0.1",
			version: "v0.0.1",
			want: Expect{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Extra: "",
			},
		},
		{
			name:    "v0.0.1-rc1",
			version: "v0.0.1-rc1",
			want: Expect{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Extra: "rc1",
			},
		},
		{
			name:    "v0.0.1-rc1-dirty",
			version: "v0.0.1-rc1-dirty",
			want: Expect{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Extra: "rc1-dirty",
			},
		},
		{
			name:    "0.0.1-rc1-SNAPSHOT-a4dca029",
			version: "0.0.1-rc1-SNAPSHOT-a4dca029",
			want: Expect{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Extra: "rc1-SNAPSHOT-a4dca029",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VersionInfo{Version: tt.version}
			v.ParseVersion()

			if v.Major != tt.want.Major ||
				v.Minor != tt.want.Minor ||
				v.Patch != tt.want.Patch ||
				v.Extra != tt.want.Extra {
				t.Errorf("Version.ParseVersion() = %v, want %v", Expect{v.Major, v.Minor, v.Patch, v.Extra}, tt.want)
			}
		})
	}
}

func TestVersion_ParseDate(t *testing.T) {
	tests := []struct {
		name      string
		buildDate string
		want      time.Time
		wantErr   bool
	}{
		{
			name:      "2019-01-01T00:00:00Z",
			buildDate: "2019-01-01T00:00:00Z",
			want:      time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "2023-10-16T00:00:00-00:00",
			buildDate: "2023-10-16T00:00:00-00:00",
			want:      time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "2023-10-16T00:00:00-03:00",
			buildDate: "2023-10-16T00:00:00-03:00",
			want:      time.Date(2023, 10, 16, 0, 0, 0, 0, time.FixedZone("UTC-3", -3*60*60)),
			wantErr:   false,
		},
		{
			name:      "2019-01-01 00:00:00",
			buildDate: "2019-01-01 00:00:00", // wrong format, dot missint between date and time
			want:      time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				v   = &VersionInfo{BuildDate: tt.buildDate}
				err = v.ParseDate()
			)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Version.ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if v.TimeStamp == nil {
				t.Errorf("Version.ParseDate() = %v, want %v", v.TimeStamp, tt.want)
				return
			}

			if !v.TimeStamp.Equal(tt.want) {
				t.Errorf("Version.ParseDate() = %v, want %v", v.TimeStamp, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	type fields struct {
		Version   string
		Commit    string
		BuildDate string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "v0.0.1",
			fields: fields{
				Version:   "v0.0.1",
				Commit:    "none",
				BuildDate: "none",
			},
			want: "v0.0.1 none none",
		},
		{
			name: "v0.0.1",
			fields: fields{
				Version:   "v0.0.1",
				Commit:    "dd00d1766495cb704a6d2c1c594800ced58e88b3",
				BuildDate: "2023-09-12.15:37:54",
			},
			want: "v0.0.1 dd00d1766495cb704a6d2c1c594800ced58e88b3 2023-09-12.15:37:54",
		},
		{
			name: "v0.0.1-rc1",
			fields: fields{
				Version:   "v0.0.1-rc1",
				Commit:    "dd00d1766495cb704a6d2c1c594800ced58e88b3",
				BuildDate: "2023-09-12.15:37:54",
			},
			want: "v0.0.1-rc1 dd00d1766495cb704a6d2c1c594800ced58e88b3 2023-09-12.15:37:54",
		},
		{
			name: "v0.0.1-rc1-dirty",
			fields: fields{
				Version:   "v0.0.1-rc1-dirty",
				Commit:    "dd00d1766495cb704a6d2c1c594800ced58e88b3",
				BuildDate: "2023-09-12.15:37:54",
			},
			want: "v0.0.1-rc1-dirty dd00d1766495cb704a6d2c1c594800ced58e88b3 2023-09-12.15:37:54",
		},
		{
			name: "1.7.0-rc1-SNAPSHOT-a4dca029",
			fields: fields{
				Version:   "1.7.0-rc1-SNAPSHOT-a4dca029",
				Commit:    "dd00d1766495cb704a6d2c1c594800ced58e88b3",
				BuildDate: "2023-09-12.15:37:54",
			},
			want: "1.7.0-rc1-SNAPSHOT-a4dca029 dd00d1766495cb704a6d2c1c594800ced58e88b3 2023-09-12.15:37:54",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VersionInfo{
				Version:   tt.fields.Version,
				Commit:    tt.fields.Commit,
				BuildDate: tt.fields.BuildDate,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
