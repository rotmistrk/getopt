package getopt

import "strconv"

func parseInt(arg string) (int64, error) {
	if len(arg) > 0 && arg[0] == '0' {
		if len(arg) > 1 {
			switch arg[1] {
			case 'x', 'X':
				return strconv.ParseInt(arg[2:], 16, 64)
			case 'b', 'B':
				return strconv.ParseInt(arg[2:], 2, 64)
			case 'd', 'D':
				return strconv.ParseInt(arg[2:], 10, 64)
			case 'o', 'O':
				return strconv.ParseInt(arg[2:], 8, 64)
			case 't', 'T':
				return strconv.ParseInt(arg[2:], 32, 64)
			case 's', 'S':
				return strconv.ParseInt(arg[2:], 64, 64)
			}
		} else {
			return strconv.ParseInt(arg[1:], 8, 64)
		}
	}
	return strconv.ParseInt(arg, 0, 64)
}

func parseUint(arg string) (uint64, error) {
	if len(arg) > 0 && arg[0] == '0' {
		if len(arg) > 1 {
			switch arg[1] {
			case 'x', 'X':
				return strconv.ParseUint(arg[2:], 16, 64)
			case 'b', 'B':
				return strconv.ParseUint(arg[2:], 2, 64)
			case 'd', 'D':
				return strconv.ParseUint(arg[2:], 10, 64)
			case 'o', 'O':
				return strconv.ParseUint(arg[2:], 8, 64)
			case 't', 'T':
				return strconv.ParseUint(arg[2:], 32, 64)
			case 's', 'S':
				return strconv.ParseUint(arg[2:], 64, 64)
			}
		} else {
			return strconv.ParseUint(arg[1:], 8, 64)
		}
	}
	return strconv.ParseUint(arg, 0, 64)
}

func (opts *GetOpt) IntValue(flag rune, longFlag string, required bool, help string) (*int64, error) {
	return opts.IntValueV([]rune{flag}, []string{longFlag}, required, help)
}

func (opts *GetOpt) IntValueV(flags []rune, longFlags []string, required bool, help string) (*int64, error) {
	var result int64
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		required:  required,
		argType:   "int",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = parseInt(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = 0
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) IntDefault(flag rune, longFlag string, value int64, help string) (*int64, error) {
	return opts.IntDefaultV([]rune{flag}, []string{longFlag}, value, help)
}

func (opts *GetOpt) IntDefaultV(flags []rune, longFlags []string, value int64, help string) (*int64, error) {
	result := value
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "int",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = parseInt(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = value
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) IntList(flag rune, longFlag string, help string) (*[]int64, error) {
	return opts.IntListV([]rune{flag}, []string{longFlag}, help)
}

func (opts *GetOpt) IntListV(flags []rune, longFlags []string, help string) (*[]int64, error) {
	result := make([]int64, 0)
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "int",
	}
	def.argConv = func(arg string) error {
		if value, err := parseInt(arg); err == nil {
			result = append(result, value)
			def.count++
			return nil
		} else {
			return err
		}
	}
	def.argReset = func() {
		result = make([]int64, 0)
	}
	return &result, opts.safeAdd(def)
}
