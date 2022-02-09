package getopt

import (
	"reflect"
	"testing"
	"time"
)

type testMarshal struct {
	Flag      bool               `flag:"b,boolean" help:"boolean value"`
	Str       string             `flag:"s,str" help:"string value"`
	StrList   []string           `flag:"S,str-list" help:"string list"`
	StrMap    map[string]string  `flag:"M,str-map" help:"string map (-Mkey:val -M k:v --str-map=ky:vl)"`
	Int       int64              `flag:"i,int-val" help:"integer value"`
	IntList   []int64            `flag:"I,int-list" help:"integer list"`
	IntMap    map[string]int     `flag:"int-map" help:"integer map"`
	Uint      uint64             `flag:"u,uint-val" help:"unsigned int value"`
	UintList  []uint64           `flag:"U,uint-list" help:"unsigned int list"`
	UintMap   map[string]int     `flag:"uint-map" help:"unsigned integer map"`
	Float     float64            `flag:"f,float-val" help:"float value"`
	FloatList []float64          `flag:"F,float-list" help:"float list"`
	FloatMap  map[string]float64 `flag:"float-map" help:"float map"`
	Wait      time.Duration      `flag:"d,duration-val" help:"duration value"`
	WaitList  []time.Duration    `flag:"D,duration-list" help:"duration list"`
	Exec      func(string) error `flag:"x,exec" help:"execute cmd"`
	Show      func() error       `flag:"show" help:"show stats"`
}

func TestGetOpt_Marshal(t *testing.T) {
	type args struct {
		target *testMarshal
		argv   []string
		posix  bool
	}
	var topic string
	var showCalled bool
	tests := []struct {
		name      string
		getOpt    *GetOpt
		args      args
		expected  testMarshal
		want      []string
		wantErr   bool
		wantTopic string
		wantShown bool
	}{
		// TODO: Add test cases.
		{
			"empty",
			New().WithDefaults("prog", "v0"),
			args{
				&testMarshal{},
				[]string{
					"prog",
				},
				true,
			},
			testMarshal{},
			[]string{},
			false,
			"",
			false,
		},
		{
			"happy path",
			New().WithDefaults("prog", "v0"),
			args{
				&testMarshal{
					Exec: func(t string) error {
						topic = t
						return nil
					},
					Show: func() error {
						showCalled = true
						return nil
					},
				},
				[]string{
					"prog",
					"-b",
					"-sstr",
					"-Sone", "-S", "two", "--str-list=three",
					"-Mleft:right", "-M", "top:bottom", "--str-map=far:close",
					"-i-1",
					"-I-2", "-I", "3", "--int-list=4",
					"--int-map=x:1", "--int-map=y:-2",
					"-u0x10",
					"-U20", "-U", "030", "--uint-list=0b1000",
					"--uint-map=x:10", "--uint-map=y:0x10",
					"-f-1.1",
					"-F1.1", "-F", "2.2", "--float-list=3.3",
					"--float-map=z:1.23",
					"-d1m",
					"-D1s", "-D2m30s", "--duration-list=3h",
					"-xls -alt",
					"--show",
					"une", "doux", "trois",
				},
				true,
			},
			testMarshal{
				Flag:      true,
				Str:       "str",
				StrList:   []string{"one", "two", "three"},
				StrMap:    map[string]string{"left": "right", "top": "bottom", "far": "close"},
				Int:       -1,
				IntList:   []int64{-2, 3, 4},
				IntMap:    map[string]int{"x": 1, "y": -2},
				Uint:      16,
				UintList:  []uint64{20, 24, 8},
				UintMap:   map[string]int{"x": 10, "y": 16},
				Float:     -1.1,
				FloatList: []float64{1.1, 2.2, 3.3},
				FloatMap:  map[string]float64{"z": 1.23},
				Wait:      1 * time.Minute,
				WaitList:  []time.Duration{time.Second, 150 * time.Second, 3 * time.Hour},
			},
			[]string{"une", "doux", "trois"},
			false,
			"ls -alt",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.getOpt
			got, err := opts.Marshal(tt.args.target, tt.args.argv, tt.args.posix)
			tt.args.target.Show = nil
			tt.args.target.Exec = nil
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
			if topic != tt.wantTopic {
				t.Errorf("Exec call is wrong - %v vs %v", topic, tt.wantTopic)
			}
			if showCalled != tt.wantShown {
				t.Errorf("Show call is wrong - %v vs %v", showCalled, tt.wantShown)
			}
		})
	}
}
