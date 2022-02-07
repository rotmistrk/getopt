package getopt

import (
	"testing"
)

func TestGetOpt_Flag(t *testing.T) {
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
			name: "not set if no args",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.Flag('b', "bool", "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   false,
			wantParseOk: true,
		},
		{
			name: "set for longflag",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.Flag('b', "--bool", "help")
			},
			args:        []string{"prog", "--bool"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   true,
			wantParseOk: true,
		},
		{
			name: "set for short flag",
			init: func(getopt *GetOpt) (*bool, error) {
				return getopt.Flag('b', "--bool", "help")
			},
			args:        []string{"prog", "-b"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   true,
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
