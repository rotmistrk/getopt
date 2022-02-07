package getopt

import (
	"reflect"
	"testing"
)

func mkstr(str string) *string {
	return &str
}

func TestTokenize(t *testing.T) {
	type args struct {
		args    []string
		options string
	}
	tests := []struct {
		name    string
		args    args
		want    []Option
		wantErr bool
	}{
		{
			name: "gnuCompat",
			args: args{
				args: []string{
					"cmd",
					"-hV",
					"-l",
					"-fiinput",
					"-o",
					"output",
					"-xt",
					"-yqquote",
					"-zr",
					"--flag=value",
					"--other",
					"zeta",
					"kappa",
					"--",
					"-o",
					"theta",
				},
				options: "hVli:o:t::q::r::xyzf",
			},
			want: []Option{
				{Opt: "-h"},
				{Opt: "-V"},
				{Opt: "-l"},
				{Opt: "-f"},
				{Opt: "-i", Arg: mkstr("input")},
				{Opt: "-o", Arg: mkstr("output")},
				{Opt: "-x"},
				{Opt: "-t"},
				{Opt: "-y"},
				{Opt: "-q", Arg: mkstr("quote")},
				{Opt: "-z"},
				{Opt: "-r"},
				{Opt: "--flag", Arg: mkstr("value")},
				{Opt: "--other"},
				{Arg: mkstr("-o")},
				{Arg: mkstr("theta")},
				{Arg: mkstr("zeta")},
				{Arg: mkstr("kappa")},
			},
		},
		{
			name: "posixlyCorrect",
			args: args{
				args: []string{
					"cmd",
					"-hV",
					"-l",
					"-fiinput",
					"-o",
					"output",
					"-xt",
					"-yqquote",
					"-zr",
					"--flag=value",
					"--other",
					"zeta",
					"kappa",
					"--",
					"-o",
					"theta",
				},
				options: "+hVli:o:t::q::r::xyzf",
			},
			want: []Option{
				{Opt: "-h"},
				{Opt: "-V"},
				{Opt: "-l"},
				{Opt: "-f"},
				{Opt: "-i", Arg: mkstr("input")},
				{Opt: "-o", Arg: mkstr("output")},
				{Opt: "-x"},
				{Opt: "-t"},
				{Opt: "-y"},
				{Opt: "-q", Arg: mkstr("quote")},
				{Opt: "-z"},
				{Opt: "-r"},
				{Opt: "--flag", Arg: mkstr("value")},
				{Opt: "--other"},
				{Arg: mkstr("zeta")},
				{Arg: mkstr("kappa")},
				{Arg: mkstr("--")},
				{Arg: mkstr("-o")},
				{Arg: mkstr("theta")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.args.args, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() got/want =\n %v\n %v", got, tt.want)
			}
		})
	}
}

func makeOptMap(values ...rune) map[rune]int {
	result := make(map[rune]int)
	for i, val := range values {
		if i&1 == 0 {
			result[val] = 0
		} else {
			result[values[i-1]] = int(val)
		}
	}
	return result
}

func Test_newGetoptConfig(t *testing.T) {
	type args struct {
		options string
	}
	tests := []struct {
		name    string
		args    args
		want    getoptConfig
		wantErr bool
	}{
		{
			name:    "emptyString",
			args:    args{options: ""},
			want:    getoptConfig{optmap: make(map[rune]int)},
			wantErr: false,
		},
		{
			name:    "justFlags",
			args:    args{options: "hV"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0)},
			wantErr: false,
		},
		{
			name:    "withRequiredParameters",
			args:    args{options: "hVi:o:"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0, 'i', 1, 'o', 1)},
			wantErr: false,
		},
		{
			name:    "withOptionalParameters",
			args:    args{options: "hVi::o::"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0, 'i', 2, 'o', 2)},
			wantErr: false,
		},
		{
			name:    "withPosixlyCorrect",
			args:    args{options: "+hVo:"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0, 'o', 1), posixlyCorrect: true},
			wantErr: false,
		},
		{
			name:    "withPositionalAsFlags",
			args:    args{options: "-hVo:"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0, 'o', 1), positionalAsArgs: true},
			wantErr: false,
		},
		{
			name:    "withNoErrorReports",
			args:    args{options: "+:hVo:"},
			want:    getoptConfig{optmap: makeOptMap('h', 0, 'V', 0, 'o', 1), posixlyCorrect: true, dontPrintErrors: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newGetoptConfig(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("newGetoptConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGetoptConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
