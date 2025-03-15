/*
Copyright © 2025 Pato Diaz pato@patodiaz.io

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"testing"

	"github.com/padiazg/ollama-tools/models/settings"
	"github.com/spf13/viper"
)

func Test_bindEnvs(t *testing.T) {
	type SingleStruct struct {
		StringField string
		IntField    int
		BoolField   bool
	}

	type NestedStruct struct {
		StringField string
		IntField    int
		BoolField   bool
		Nested      SingleStruct
	}

	type NestedPointerStruct struct {
		StringField string
		IntField    int
		BoolField   bool
		Nested      *SingleStruct
	}

	type NestedEmbededStruct struct {
		SingleStruct
	}

	type UnexpectedStruct struct {
		ExportedField   string
		unexportedField string
	}

	type args struct {
		i     interface{}
		parts []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single struct",
			args: args{
				i:     SingleStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},

		{
			name: "single pointer struct",
			args: args{
				i:     &SingleStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},
		{
			name: "nested struct",
			args: args{
				i:     NestedStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
				"nested.stringfield",
				"nested.intfield",
				"nested.boolfield",
			},
		},
		{
			name: "nested pointer struct",
			args: args{
				i: &NestedPointerStruct{
					Nested: &SingleStruct{},
				},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
				"nested.stringfield",
				"nested.intfield",
				"nested.boolfield",
			},
		},
		{
			name: "nested embeded struct",
			args: args{
				i:     NestedEmbededStruct{},
				parts: []string{},
			},
			want: []string{
				"singlestruct_stringfield",
				"singlestruct_intfield",
				"singlestruct_boolfield",
			},
		},
		{
			name: "unexpected struct",
			args: args{
				i:     UnexpectedStruct{},
				parts: []string{},
			},
			want: []string{
				"exportedfield",
			},
		},
		{
			name: "unitialized pointer struct",
			args: args{
				i:     &NestedPointerStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},
		{
			name: "settigns",
			args: args{
				i:     &settings.Settings{},
				parts: []string{},
			},
			want: []string{
				"ollamaurl",
				// "webserver.port",
				// "webserver.tlsenabled",
				// "webserver.tlsminversion",
				// "webserver.tlsinsecureskipverify",
				// "webserver.tlscertificates",
				// "webserver.static.path",
				// "webserver.static.route",
				// "database.dialect",
				// "database.connectionstring",
				// "database.enablelogs",
				// "database.maxretries",
				// "database.retryinterval",
				// "log.level",
				// "log.format",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindEnvs(tt.args.i, tt.args.parts...)
			keys := viper.AllKeys()
			for _, key := range tt.want {
				found := false
				for _, k := range keys {
					if k == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("bindEnvs() missing key: %s", key)
				}
			}
		})
	}
}
