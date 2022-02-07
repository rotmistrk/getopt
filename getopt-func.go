package getopt

func (opts *GetOpt) ArgFunc(flag rune, longFlag string, action func(string) error, help string) error {
	return opts.ArgFuncV([]rune{flag}, []string{longFlag}, action, help)
}

func (opts *GetOpt) ArgFuncV(flags []rune, longFlags []string, action func(string) error, help string) error {
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		argType:   "value",
	}
	def.argConv = action
	return opts.safeAdd(def)
}

func (opts *GetOpt) FlagFunc(flag rune, longFlag string, action func() error, help string) error {
	return opts.FlagFuncV([]rune{flag}, []string{longFlag}, action, help)
}

func (opts *GetOpt) FlagFuncV(flags []rune, longFlags []string, action func() error, help string) error {
	def := optDef{
		posixOpts: flags,
		longOpts:  longFlags,
		help:      help,
		noArg:     true,
	}
	def.argConv = func(string) error {
		return action()
	}
	return opts.safeAdd(def)
}
