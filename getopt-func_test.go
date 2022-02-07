package getopt

import (
	"reflect"
	"testing"
)

func TestGetOpt_Func(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) error
		args        []string
		posix       bool
		wantSetupOk bool
		valueList   []string
		wantParseOk bool
	}
	var args []string
	tests := []test{
		{
			name: "not set if no args",
			init: func(getopt *GetOpt) error {
				return getopt.FlagFunc('b', "flag", func() error { args = append(args, ""); return nil }, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{},
			wantParseOk: true,
		},
		{
			name: "set for longflag",
			init: func(getopt *GetOpt) error {
				return getopt.FlagFunc('b', "--flag", func() error { args = append(args, ""); return nil }, "help")
			},
			args:        []string{"prog", "--flag"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{""},
			wantParseOk: true,
		},
		{
			name: "set for short flag",
			init: func(getopt *GetOpt) error {
				return getopt.FlagFunc('b', "--flag", func() error { args = append(args, ""); return nil }, "help")
			},
			args:        []string{"prog", "-b"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{""},
			wantParseOk: true,
		},
		{
			name: "set for few flags",
			init: func(getopt *GetOpt) error {
				return getopt.FlagFunc('b', "--flag", func() error { args = append(args, ""); return nil }, "help")
			},
			args:        []string{"prog", "-bbb", "--flag", "--flag"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{"", "", "", "", ""},
			wantParseOk: true,
		},
		{
			name: "argfunc not set if no args",
			init: func(getopt *GetOpt) error {
				return getopt.ArgFunc('b', "flag", func(arg string) error { args = append(args, arg); return nil }, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{},
			wantParseOk: true,
		},
		{
			name: "argfunc set for longflag",
			init: func(getopt *GetOpt) error {
				return getopt.ArgFunc('b', "--flag", func(arg string) error { args = append(args, arg); return nil }, "help")
			},
			args:        []string{"prog", "--flag=hello"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{"hello"},
			wantParseOk: true,
		},
		{
			name: "argfunc set for short flag",
			init: func(getopt *GetOpt) error {
				return getopt.ArgFunc('b', "--flag", func(arg string) error { args = append(args, arg); return nil }, "help")
			},
			args:        []string{"prog", "-boot"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{"oot"},
			wantParseOk: true,
		},
		{
			name: "argfunc set for few flags",
			init: func(getopt *GetOpt) error {
				return getopt.ArgFunc('b', "--flag", func(arg string) error { args = append(args, arg); return nil }, "help")
			},
			args:        []string{"prog", "-bbb", "--flag=abc", "--flag=def"},
			posix:       true,
			wantSetupOk: true,
			valueList:   []string{"bb", "abc", "def"},
			wantParseOk: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			getopt := New()
			args = make([]string, 0)
			if err := test.init(getopt); err != nil {
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
					if !reflect.DeepEqual(test.valueList, args) {
						t.Errorf("Unexpected value: %v (expected: %v)", args, test.valueList)
					}
				}
			}
		})
	}
}
