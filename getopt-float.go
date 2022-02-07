package getopt

import "strconv"

func (opts *GetOpt) FloatValue(flag rune, longFlag string, required bool, help string) (*float64, error) {
	return opts.FloatValueV([]rune{flag}, []string{longFlag}, required, help)
}

func (opts *GetOpt) FloatValueV(flags []rune, longFlags []string, required bool, help string) (*float64, error) {
	var result float64
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		required:  required,
		argType:   "float",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = strconv.ParseFloat(arg, 64)
		def.count++
		return err
	}
	def.argReset = func() {
		result = 0
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) FloatDefault(flag rune, longFlag string, value float64, help string) (*float64, error) {
	return opts.FloatDefaultV([]rune{flag}, []string{longFlag}, value, help)
}
func (opts *GetOpt) FloatDefaultV(flags []rune, longFlags []string, value float64, help string) (*float64, error) {
	result := value
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "float",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = strconv.ParseFloat(arg, 64)
		def.count++
		return err
	}
	def.argReset = func() {
		result = value
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) FloatList(flag rune, longFlag string, help string) (*[]float64, error) {
	return opts.FloatListV([]rune{flag}, []string{longFlag}, help)
}

func (opts *GetOpt) FloatListV(flags []rune, longFlags []string, help string) (*[]float64, error) {
	result := make([]float64, 0)
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "float",
	}
	def.argConv = func(arg string) error {
		if value, err := strconv.ParseFloat(arg, 64); err == nil {
			result = append(result, value)
			def.count++
			return nil
		} else {
			return err
		}
	}
	def.argReset = func() {
		result = make([]float64, 0)
	}
	return &result, opts.safeAdd(def)
}
