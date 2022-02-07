package getopt

import (
	"reflect"
	"testing"
)

func TestGetOpt_Int(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*int64, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   int64
		wantParseOk bool
	}
	tests := []test{
		{
			name: "optional no args",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: true,
		},
		{
			name: "optional longopt set",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog", "--int=10"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   10,
			wantParseOk: true,
		},
		{
			name: "optional longopt no arg should fail",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog", "--int"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "optional opt set",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog", "-f", "010"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   010,
			wantParseOk: true,
		},
		{
			name: "optional opt set no space",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog", "-f0x10"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0x10,
			wantParseOk: true,
		},
		{
			name: "optional flag no arg should fail",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", false, "help")
			},
			args:        []string{"prog", "-f"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "required flag fails if missing",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntValue('f', "--int", true, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "default flag value is used",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntDefault('f', "--int", 20, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   20,
			wantParseOk: true,
		},
		{
			name: "default flag value is replaced",
			init: func(getopt *GetOpt) (*int64, error) {
				return getopt.IntDefault('f', "--int", 20, "help")
			},
			args:        []string{"prog", "-f0b100"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   4,
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

func TestGetOpt_IntList(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*[]int64, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   []int64
		wantParseOk bool
	}
	tests := []test{
		{
			name: "empty list",
			init: func(getopt *GetOpt) (*[]int64, error) {
				return getopt.IntList('f', "--int", "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []int64{},
			wantParseOk: true,
		},
		{
			name: "single value",
			init: func(getopt *GetOpt) (*[]int64, error) {
				return getopt.IntList('f', "--int", "help")
			},
			args:        []string{"prog", "-f0tAB"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []int64{331},
			wantParseOk: true,
		},
		{
			name: "two values",
			init: func(getopt *GetOpt) (*[]int64, error) {
				return getopt.IntList('f', "--int", "help")
			},
			args:        []string{"prog", "-f0tAB", "--int=-0x11"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []int64{331, -17},
			wantParseOk: true,
		},
		{
			name: "three values",
			init: func(getopt *GetOpt) (*[]int64, error) {
				return getopt.IntList('f', "--int", "help")
			},
			args:        []string{"prog", "-f0d10", "--int=0x10", "-f", "0o10"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []int64{10, 16, 8},
			wantParseOk: true,
		},
		{
			name: "four values",
			init: func(getopt *GetOpt) (*[]int64, error) {
				return getopt.IntList('f', "--int", "help")
			},
			args:        []string{"prog", "-f10", "--int=0b10", "-f", "0x10", "--int=010"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []int64{10, 2, 16, 8},
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
					if !reflect.DeepEqual(test.wantValue, *result) {
						t.Errorf("Unexpected value: %v (expected: %v)", *result, test.wantValue)
					}
				}
			}
		})
	}
}
