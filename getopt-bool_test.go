package getopt

import (
	"testing"
)

func TestGetOpt_Bool(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*bool, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   bool
		wantParseOk bool
	}
	tests := []test{
		{
			name: "optional no args",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "bool", false, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: true,
		},
		{
			name: "optional flag true",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "--bool", false, "help")
			},
			args:        []string{"prog", "--bool=true"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   true,
			wantParseOk: true,
		},
		{
			name: "optional flag false",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "--bool", false, "help")
			},
			args:        []string{"prog", "--bool=false"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: true,
		},
		{
			name: "optional flag no arg should fail",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "bool", false, "help")
			},
			args:        []string{"prog", "--bool"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: false,
		},
		{
			name: "optional flag true",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "bool", false, "help")
			},
			args:        []string{"prog", "-b", "true"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   true,
			wantParseOk: true,
		},
		{
			name: "optional flag false",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "--bool", false, "help")
			},
			args:        []string{"prog", "-b", "false"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: true,
		},
		{
			name: "optional flag no arg should fail",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "bool", false, "help")
			},
			args:        []string{"prog", "-b"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: false,
		},
		{
			name: "required flag fails if missing",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolValue('b', "bool", true, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: false,
		},
		{
			name: "default flag value is used",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolDefault('b', "bool", true, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   true,
			wantParseOk: true,
		},
		{
			name: "default flag value is replaced",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.BoolDefault('b', "bool", true, "help")
			},
			args:        []string{"prog", "-bfalse"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			getopt := New()
			if result, err := test.init(getopt); err != nil {
				if test.wantSetupOk {
					t.Errorf("Unexpected error %n on setup", err)
				}
			} else {
				if !test.wantSetupOk {
					t.Errorf("Expected setup to fail")
				}
				if _, err := getopt.Parse(test.args, test.posix); err != nil {
					if test.wantParseOk {
						t.Errorf("Unexpected error %n on parse", err)
					}
				} else {
					if !test.wantParseOk {
						t.Errorf("Expeted parse to fail")
					}
					if test.wantValue != *result {
						t.Errorf("Unexpected value: %v (expected: %v)", *result, test.wantValue)
					}
				}
			}
		})
	}
}
