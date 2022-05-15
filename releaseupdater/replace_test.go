package releaseupdater

import (
	"reflect"
	"testing"
)

// note: only the ?P<name> grouping variant is support in go
const referenceRE = `(?P<key>\s*version:\s*)(?P<value>v.*?)(?P<other>$|\s+\.*)`
const withQuotesRE = `(?P<key>\s*version:\s*")(?P<value>(v.*?))(?P<other>$|"|\s+\.*?)`
const complexInput = `
version: v0.0.0-rc 

	version: v0.0.0
 version:v0.0.0
abcdef 
xyz
# Hallo version: x  
version:v1.0.0 #
`
const complexResult = `
version: v1.0 

	version: v1.0
 version:v1.0
abcdef 
xyz
# Hallo version: x  
version:v1.0 #
`

func TestReplaceVersionString(t *testing.T) {
	type args struct {
		content    string
		regex      string
		newVersion string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty", args: args{content: "", regex: ".*", newVersion: "v1.0"}, want: "v1.0"},
		{name: "invalid re", args: args{content: "", regex: "?<*", newVersion: "v1.0"}, want: ""},
		{name: "key value complete", args: args{content: "version: v0.0.0", regex: referenceRE, newVersion: "v1.0"}, want: "version: v1.0"},
		{name: "with quotes", args: args{content: "version: \"v0.0.0\"", regex: withQuotesRE, newVersion: "v1.0"}, want: "version: \"v1.0\""},
		{name: "complex multiline", args: args{
			regex: referenceRE, newVersion: "v1.0",
			content: complexInput},
			want: complexResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceVersionString(tt.args.content, tt.args.regex, tt.args.newVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
