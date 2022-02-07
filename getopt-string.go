package getopt

func (opts *GetOpt) StringValue(flag rune, longFlag string, required bool, help string) (*string, error) {
	return opts.StringValueV([]rune{flag}, []string{longFlag}, required, help)
}

func (opts *GetOpt) StringValueV(flags []rune, longFlags []string, required bool, help string) (*string, error) {
	var result string
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		required:  required,
		help:      help,
		argType:   "string",
	}
	def.argConv = func(arg string) error {
		result = arg
		def.count++
		return nil
	}
	def.argReset = func() {
		result = ""
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) StringDefault(flag rune, longFlag string, value string, help string) (*string, error) {
	return opts.StringDefaultV([]rune{flag}, []string{longFlag}, value, help)
}

func (opts *GetOpt) StringDefaultV(flags []rune, longFlags []string, value string, help string) (*string, error) {
	result := value
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "string",
	}
	def.argConv = func(arg string) error {
		result = arg
		def.count++
		return nil
	}
	def.argReset = func() {
		result = value
	}
	return &result, opts.safeAdd(def)
}

func (opts *GetOpt) StringList(flag rune, longFlag string, help string) (*[]string, error) {
	return opts.StringListV([]rune{flag}, []string{longFlag}, help)
}

func (opts *GetOpt) StringListV(flags []rune, longFlags []string, help string) (*[]string, error) {
	result := make([]string, 0)
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "string",
	}
	def.argConv = func(arg string) error {
		result = append(result, arg)
		def.count++
		return nil
	}
	def.argReset = func() {
		result = make([]string, 0)
	}
	return &result, opts.safeAdd(def)
}
