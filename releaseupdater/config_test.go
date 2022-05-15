package releaseupdater

import (
	"reflect"
	"testing"
)

func Test_getConfig(t *testing.T) {
	type args struct {
		configFilename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test a",
			args: args{configFilename: "testdata/testa.yml"},
			want: &Config{Actions: []Action{{From: "demo_app1", Update: []Update{{Repo: "demo_deploy1", File: "versions.properties", Regex: `(?P<key>\s*version:\s*)(?P<value>v.*?)(?P<other>$|\s+\.*)`}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig(tt.args.configFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigFrom(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "empty",
			args: args{content: []byte(``)},
			want: &Config{},
		},
		{
			name: "x",
			args: args{content: []byte(`actions: []`)},
			want: &Config{Actions: []Action{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConfigFrom(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfigFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfigFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
