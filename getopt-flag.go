package getopt

func (opts *GetOpt) Flag(flag rune, longFlag string, help string) (*bool, error) {
	return opts.FlagV([]rune{flag}, []string{longFlag}, help)
}

func (opts *GetOpt) FlagV(flags []rune, longFlags []string, help string) (*bool, error) {
	var result bool
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		noArg:     true,
	}
	def.argConv = func(arg string) error {
		result = true
		def.count++
		return nil
	}
	def.argReset = func() {
		result = false
	}
	return &result, opts.safeAdd(def)
}
