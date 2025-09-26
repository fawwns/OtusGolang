package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    Environment
		wantErr bool
	}{
		{
			name: "valid envdir",
			dir:  "testdata/env",
			want: Environment{
				"HELLO": EnvValue{Value: "hello", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
			},
			wantErr: false,
		},
		{
			name:    "non-existing dir",
			dir:     "testdata/not_exists",
			want:    Environment{},
			wantErr: true,
		},
		{
			name:    "empty dir",
			dir:     "testdata/empty",
			want:    Environment{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.dir)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
