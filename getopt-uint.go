package getopt

func (opts *GetOpt) UintValue(flag rune, longFlag string, required bool, help string) (*uint64, error) {
	return opts.UintValueV([]rune{flag}, []string{longFlag}, required, help)
}

func (opts *GetOpt) UintValueV(flags []rune, longFlags []string, required bool, help string) (*uint64, error) {
	var result uint64
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		required:  required,
		argType:   "uint",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = parseUint(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = 0
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) UintDefault(flag rune, longFlag string, value uint64, help string) (*uint64, error) {
	return opts.UintDefaultV([]rune{flag}, []string{longFlag}, value, help)
}

func (opts *GetOpt) UintDefaultV(flags []rune, longFlags []string, value uint64, help string) (*uint64, error) {
	result := value
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "uint",
	}
	def.argConv = func(arg string) error {
		var err error
		result, err = parseUint(arg)
		def.count++
		return err
	}
	def.argReset = func() {
		result = value
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) UintList(flag rune, longFlag string, help string) (*[]uint64, error) {
	return opts.UintListV([]rune{flag}, []string{longFlag}, help)
}

func (opts *GetOpt) UintListV(flags []rune, longFlags []string, help string) (*[]uint64, error) {
	result := make([]uint64, 0)
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "uint",
	}
	def.argConv = func(arg string) error {
		if value, err := parseUint(arg); err == nil {
			result = append(result, value)
			def.count++
			return nil
		} else {
			return err
		}
	}
	def.argReset = func() {
		result = make([]uint64, 0)
	}
	return &result, opts.safeAdd(def)
}
