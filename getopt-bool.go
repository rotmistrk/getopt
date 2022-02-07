package getopt

import "strconv"

func (opts *GetOpt) BoolValue(flag rune, longFlag string, required bool, help string) (*bool, error) {
	return opts.BoolValueV([]rune{flag}, []string{longFlag}, required, help)
}

func (opts *GetOpt) BoolValueV(flags []rune, longFlags []string, required bool, help string) (*bool, error) {
	var result bool
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		required:  required,
		argType:   "bool",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = strconv.ParseBool(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = false
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) BoolDefault(flag rune, longFlag string, value bool, help string) (*bool, error) {
	return opts.BoolDefaultV([]rune{flag}, []string{longFlag}, value, help)
}

func (opts *GetOpt) BoolDefaultV(flags []rune, longFlags []string, value bool, help string) (*bool, error) {
	result := value
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "bool",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = strconv.ParseBool(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = value
	}
	return &result, opts.safeAdd(def)
}
