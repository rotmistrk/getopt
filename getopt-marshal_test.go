package getopt

import (
	"reflect"
	"testing"
	"time"
)

type testMarshal struct {
	Flag      bool            `flag:"b,boolean" help:"boolean value"`
	Str       string          `flag:"s,str" help:"string value"`
	StrList   []string        `flag:"S,str-list" help:"string list"`
	Int       int64           `flag:"i,int-val" help:"integer value"`
	IntList   []int64         `flag:"I,int-list" help:"integer list"`
	Uint      uint64          `flag:"u,uint-val" help:"unsigned int value"`
	UintList  []uint64        `flag:"U,uint-list" help:"unsigned int list"`
	Float     float64         `flag:"f,float-val" help:"float value"`
	FloatList []float64       `flag:"F,float-list" help:"float list"`
	Wait      time.Duration   `flag:"d,duration-val" help:"duration value"`
	WaitList  []time.Duration `flag:"D,duration-list" help:"duration list"`
}

func TestGetOpt_Marshal(t *testing.T) {
	type args struct {
		target *testMarshal
		argv   []string
		posix  bool
	}
	tm := testMarshal{}
	tests := []struct {
		name     string
		getOpt   *GetOpt
		args     args
		expected testMarshal
		want     []string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"empty",
			New().WithDefaults("prog", "v0"),
			args{
				&tm,
				[]string{
					"prog",
				},
				true,
			},
			testMarshal{},
			[]string{},
			false,
		},
		{
			"happy path",
			New().WithDefaults("prog", "v0"),
			args{
				&tm,
				[]string{
					"prog",
					"-b",
					"-sstr",
					"-Sone", "-S", "two", "--str-list=three",
					"-i-1",
					"-I-2", "-I", "3", "--int-list=4",
					"-u0x10",
					"-U20", "-U", "030", "--uint-list=0b1000",
					"-f-1.1",
					"-F1.1", "-F", "2.2", "--float-list=3.3",
					"-d1m",
					"-D1s", "-D2m30s", "--duration-list=3h",
					"une", "doux", "trois",
				},
				true,
			},
			testMarshal{
				Flag:      true,
				Str:       "str",
				StrList:   []string{"one", "two", "three"},
				Int:       -1,
				IntList:   []int64{-2, 3, 4},
				Uint:      16,
				UintList:  []uint64{20, 24, 8},
				Float:     -1.1,
				FloatList: []float64{1.1, 2.2, 3.3},
				Wait:      1 * time.Minute,
				WaitList:  []time.Duration{time.Second, 150 * time.Second, 3 * time.Hour},
			},
			[]string{"une", "doux", "trois"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.getOpt
			got, err := opts.Marshal(tt.args.target, tt.args.argv, tt.args.posix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.expected, *tt.args.target) {
				t.Errorf("Marshal() got = %v, want %v", tt.expected, *tt.args.target)
			}
		})
	}
}
