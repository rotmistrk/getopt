package getopt

import (
	"reflect"
	"testing"
)

func TestGetOpt_Float(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*float64, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   float64
		wantParseOk bool
	}
	tests := []test{
		{
			name: "optional no args",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: true,
		},
		{
			name: "optional longopt set",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog", "--float=1.61"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   1.61,
			wantParseOk: true,
		},
		{
			name: "optional longopt no arg should fail",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog", "--float"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "optional opt set",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog", "-f", "1.61"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   1.61,
			wantParseOk: true,
		},
		{
			name: "optional opt set no space",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog", "-f1.61"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   1.61,
			wantParseOk: true,
		},
		{
			name: "optional flag no arg should fail",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", false, "help")
			},
			args:        []string{"prog", "-f"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "required flag fails if missing",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatValue('f', "--float", true, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   0,
			wantParseOk: false,
		},
		{
			name: "default flag value is used",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatDefault('f', "--float", 3.14, "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   3.14,
			wantParseOk: true,
		},
		{
			name: "default flag value is replaced",
			init: func(getopt *GetOpt) (*float64, error) {
				return getopt.FloatDefault('f', "--float", 3.14, "help")
			},
			args:        []string{"prog", "-f1.61"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   1.61,
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

func TestGetOpt_FloatList(t *testing.T) {
	type test struct {
		name        string
		init        func(getopt *GetOpt) (*[]float64, error)
		args        []string
		posix       bool
		wantSetupOk bool
		wantValue   []float64
		wantParseOk bool
	}
	tests := []test{
		{
			name: "empty list",
			init: func(getopt *GetOpt) (*[]float64, error) {
				return getopt.FloatList('f', "--float", "help")
			},
			args:        []string{"prog"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []float64{},
			wantParseOk: true,
		},
		{
			name: "single value",
			init: func(getopt *GetOpt) (*[]float64, error) {
				return getopt.FloatList('f', "--float", "help")
			},
			args:        []string{"prog", "-f1.1"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []float64{1.1},
			wantParseOk: true,
		},
		{
			name: "two values",
			init: func(getopt *GetOpt) (*[]float64, error) {
				return getopt.FloatList('f', "--float", "help")
			},
			args:        []string{"prog", "-f1.1", "--float=2.2"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []float64{1.1, 2.2},
			wantParseOk: true,
		},
		{
			name: "three values",
			init: func(getopt *GetOpt) (*[]float64, error) {
				return getopt.FloatList('f', "--float", "help")
			},
			args:        []string{"prog", "-f1.1", "--float=2.2", "-f", "3.3"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []float64{1.1, 2.2, 3.3},
			wantParseOk: true,
		},
		{
			name: "four values",
			init: func(getopt *GetOpt) (*[]float64, error) {
				return getopt.FloatList('f', "--float", "help")
			},
			args:        []string{"prog", "-f1.1", "--float=2.2", "-f", "3.3", "--float=4.4"},
			posix:       true,
			wantSetupOk: true,
			wantValue:   []float64{1.1, 2.2, 3.3, 4.4},
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
