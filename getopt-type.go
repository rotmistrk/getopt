package getopt

import (
	"errors"
	"fmt"
	"os"
)

type optDef struct {
	posixOpts []rune
	longOpts  []string
	help      string
	noArg     bool
	required  bool
	multiple  bool
	count     int
	argConv   func(string) error
	argReset  func()
	argType   string
}

func (optDef *optDef) Reset() {
	if optDef.argReset != nil {
		optDef.argReset()
	}
	optDef.count = 0
}

type ErrorHandler interface {
	Handle(err error, option Option) error
}

type GetOpt struct {
	description  []string
	done         bool
	errorHandler func(err error, option Option) (bool, error)
	name         string
	optionMap    map[string]optDef
	optionList   []optDef
	version      string
}

func New() *GetOpt {
	return &GetOpt{
		errorHandler: func(err error, option Option) (bool, error) {
			_, _ = fmt.Fprintln(os.Stderr, err, " while handling ", option)
			return true, err
		},
		description: make([]string, 0),
		optionMap:   make(map[string]optDef),
		optionList:  make([]optDef, 0),
	}
}

func (opts *GetOpt) AddDefaults(name string, version string, description []string) {
	opts.SetName(name)
	opts.SetVersion(version)
	opts.SetDescription(description)
	_ = opts.FlagFunc('h', "--help", func() error { return opts.Help() }, "Print help")
	_ = opts.FlagFunc('V', "--version", func() error { return opts.Version() }, "Print version")
}

func (opts *GetOpt) WithDefaults(programName string, version string, description ...string) *GetOpt {
	opts.AddDefaults(programName, version, description)
	return opts
}

func (opts *GetOpt) SetErrorHandler(handler func(err error, option Option) (bool, error)) {
	opts.errorHandler = handler
}

func (opts *GetOpt) WithErrorHandler(handler func(err error, option Option) (bool, error)) *GetOpt {
	opts.SetErrorHandler(handler)
	return opts
}

func (opts *GetOpt) ResetValues() {
	for _, v := range opts.optionMap {
		v.Reset()
	}
}

func (opts *GetOpt) separateFlagsFromLognopts(synonyms []string) ([]rune, []string) {
	flags := make([]rune, 0)
	longopts := make([]string, 0)
	for _, synonym := range synonyms {
		if len(synonym) == 0 {
			continue
		}
		runes := []rune(synonym)
		if len(runes) == 1 {
			flags = append(flags, runes[0])
		} else {
			longopts = append(longopts, "--"+synonym)
		}
	}
	return flags, longopts
}

func (opts *GetOpt) Parse(args []string, posix bool) ([]string, error) {
	optstring := ""
	if posix {
		optstring += "+"
	}
	for _, v := range opts.optionList {
		for _, posixOpt := range v.posixOpts {
			optstring += string(posixOpt)
			if !v.noArg {
				optstring += ":"
			}
		}
	}
	content, err := Tokenize(args, optstring)
	if err != nil {
		opts.done, err = opts.errorHandler(err, Option{"", &args[0]})
	}
	positional := make([]string, 0)
	for _, opt := range content {
		if opt.Opt != "" {
			if item, found := opts.optionMap[opt.Opt]; found == true {
				if item.noArg {
					err = item.argConv("")
				} else if opt.Arg == nil {
					err = errors.New("Argument required for " + opt.Opt + " (" + item.help + ")")
				} else {
					err = item.argConv(*opt.Arg)
				}
			} else {
				err = errors.New("Unknown option `" + opt.Opt + "`")
			}
		} else if opt.Arg != nil {
			positional = append(positional, *opt.Arg)
		} else {
			err = errors.New("Unexpedted empty option: no flag, no arg")
		}
		if err != nil {
			opts.done, err = opts.errorHandler(err, opt)
		}
	}
	for _, opt := range opts.optionList {
		if opt.required && opt.count == 0 {
			optflag := ""
			if len(opt.longOpts) > 0 {
				optflag = opt.longOpts[0]
			} else {
				optflag = "-" + string(opt.posixOpts[0])
			}
			opts.done, err = opts.errorHandler(errors.New("Missing required option"), Option{optflag, nil})
		}
	}
	return positional, err
}

func (opts GetOpt) Done() bool {
	return opts.done
}

func (opts *GetOpt) safeAdd(def optDef) error {
	for _, posixOpt := range def.posixOpts {
		if err := opts.safeAddKey(string("-"+string(posixOpt)), def); err != nil {
			return err
		}
	}
	for _, longOpt := range def.longOpts {
		if err := opts.safeAddKey(longOpt, def); err != nil {
			return err
		}
	}
	opts.optionList = append(opts.optionList, def)
	return nil
}

func (opts *GetOpt) safeAddKey(option string, opt optDef) error {
	if val, found := opts.optionMap[option]; found {
		return errors.New("Duplicate optionMap key: " + option + ": " + val.help + " & " + opt.help)
	} else {
		opts.optionMap[option] = opt
	}
	return nil
}

func (opts *GetOpt) Help() error {
	opts.done = true
	// Posix requires help to be printed to stdout,
	// that makes total sense as help is the requested result
	for _, desc := range opts.description {
		fmt.Println(desc)
	}
	for _, opt := range opts.optionList {
		arg := opt.argType
		nl := ""
		for _, f := range opt.longOpts {
			sep := "="
			if arg == "" {
				sep = ""
			}
			fmt.Printf("%s\t%s%s%s", nl, f, sep, arg)
			nl = "\n"
		}
		for _, f := range opt.posixOpts {
			fmt.Printf("%s\t-%c %s", nl, f, arg)
			nl = "\n"
		}
		fmt.Print("\t")
		nl = ""
		if opt.required {
			fmt.Print("required")
			nl = ", "
		}
		if opt.multiple {
			fmt.Printf("%smultiple", nl)
			nl = ", "
		}
		fmt.Printf("%s%s\n", nl, opt.help)
	}
	return nil
}

func (opts *GetOpt) Version() error {
	opts.done = true
	fmt.Println(opts.name, opts.version)
	return nil
}

func (opts *GetOpt) SetName(name string) {
	opts.name = name
}

func (opts *GetOpt) SetVersion(version string) {
	opts.version = version
}

func (opts *GetOpt) SetDescription(description []string) {
	opts.description = description
}
