package main

import (
	"os"
	"testing"
)

func Test_checkRequiredEnvVars(t *testing.T) {
	type args struct {
		requiredEnvVars []string
		actualEnvVars   []struct {
			key   string
			value string
		}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success test",
			args: args{
				requiredEnvVars: []string{
					"TEST_ENV_VAR1",
					"TEST_ENV_VAR2",
					"TEST_ENV_VAR3",
				},
				actualEnvVars: []struct {
					key   string
					value string
				}{
					{
						key:   "TEST_ENV_VAR1",
						value: "test1",
					},
					{
						key:   "TEST_ENV_VAR2",
						value: "test2",
					},
					{
						key:   "TEST_ENV_VAR3",
						value: "test3",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing environment variable",
			args: args{
				requiredEnvVars: []string{
					"TEST_ENV_VAR1",
					"TEST_ENV_VAR2",
					"TEST_ENV_VAR3",
				},
				actualEnvVars: []struct {
					key   string
					value string
				}{
					{
						key:   "TEST_ENV_VAR1",
						value: "test1",
					},
					{
						key:   "TEST_ENV_VAR2",
						value: "test2",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Clearenv()

			var keys []string

			for _, envVar := range tt.args.actualEnvVars {
				_ = os.Setenv(envVar.key, envVar.value)
				keys = append(keys, envVar.key)
			}

			if err := checkRequiredEnvVars(tt.args.requiredEnvVars); (err != nil) != tt.wantErr {
				t.Errorf("checkRequiredEnvVars() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
