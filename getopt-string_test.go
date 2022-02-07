package getopt

import (
	"reflect"
	"testing"
)

func TestGetOpt_String(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*string, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   string
		wantParseOk bool
	}
	tests := []test{
		{
			name: "optional no args",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "",
			wantParseOk: true,
		},
		{
			name: "optional longopt set",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "--str=abc"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "abc",
			wantParseOk: true,
		},
		{
			name: "optional longopt no arg should fail",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "--str"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "",
			wantParseOk: false,
		},
		{
			name: "optional longopt separated arg should fail",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "--str", "val"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "",
			wantParseOk: false,
		},
		{
			name: "optional opt set",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "-f", "abc"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "abc",
			wantParseOk: true,
		},
		{
			name: "optional opt set no space",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "-face"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "ace",
			wantParseOk: true,
		},
		{
			name: "optional flag no arg should fail",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", false, "help")
			},
			args:        []string{"prog", "-f"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "",
			wantParseOk: false,
		},
		{
			name: "required flag fails if missing",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringValue('f', "--str", true, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "",
			wantParseOk: false,
		},
		{
			name: "default flag value is used",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringDefault('f', "--str", "hello", "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "hello",
			wantParseOk: true,
		},
		{
			name: "default flag value is replaced",
			init: func(getopt *GetOpt) (*string, error) {
				return getopt.StringDefault('f', "--str", "hello", "help")
			},
			args:        []string{"prog", "-fworld"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   "world",
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

func TestGetOpt_StringList(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*[]string, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   []string
		wantParseOk bool
	}
	tests := []test{
		{
			name: "empty list",
			init: func(getopt *GetOpt) (*[]string, error) {
				return getopt.StringList('f', "--str", "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []string{},
			wantParseOk: true,
		},
		{
			name: "single value",
			init: func(getopt *GetOpt) (*[]string, error) {
				return getopt.StringList('f', "--str", "help")
			},
			args:        []string{"prog", "-feet"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []string{"eet"},
			wantParseOk: true,
		},
		{
			name: "two values",
			init: func(getopt *GetOpt) (*[]string, error) {
				return getopt.StringList('f', "--str", "help")
			},
			args:        []string{"prog", "-foot", "--str=abc"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []string{"oot", "abc"},
			wantParseOk: true,
		},
		{
			name: "three values",
			init: func(getopt *GetOpt) (*[]string, error) {
				return getopt.StringList('f', "--str", "help")
			},
			args:        []string{"prog", "-face", "--str=abc", "-f", "def"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []string{"ace", "abc", "def"},
			wantParseOk: true,
		},
		{
			name: "four values",
			init: func(getopt *GetOpt) (*[]string, error) {
				return getopt.StringList('f', "--str", "help")
			},
			args:        []string{"prog", "-fun", "--str=abc", "-f", "def", "--str=", "-f", ""},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []string{"un", "abc", "def", "", ""},
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
