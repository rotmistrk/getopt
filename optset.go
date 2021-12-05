package getopt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type optDef struct {
	posixOpt rune
	longOpt  string
	help     string
	noArg    bool
	required bool
	multiple bool
	count    int
	argConv  func(string) error
	argReset func()
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

type OptSet struct {
	description  []string
	done         bool
	errorHandler func(err error, option Option) (bool, error)
	name         string
	options      map[string]optDef
	version      string
}

func NewOptSet() OptSet {
	return OptSet{
		errorHandler: func(err error, option Option) (bool, error) {
			fmt.Fprintln(os.Stderr, err, " while handling ", option)
			return true, err
		},
		description: make([]string, 0),
		options:     make(map[string]optDef),
	}
}

func (optSet *OptSet) AddDefaults(name string, version string, description []string) {
	optSet.SetName(name)
	optSet.SetVersion(version)
	optSet.SetDescription(description)
	optSet.FlagFunc('h', "--help", func() error { return optSet.Help() }, "Print help")
	optSet.FlagFunc('V', "--version", func() error { return optSet.Version() }, "Print version")
}

func (optset *OptSet) SetErrorHandler(handler func(err error, option Option) (bool, error)) {
	optset.errorHandler = handler
}

func (optSet *OptSet) ResetValues() {
	for _, v := range optSet.options {
		v.Reset()
	}
}

func (optSet *OptSet) Parse(args []string, posix bool) ([]string, error) {
	optstring := ""
	if posix {
		optstring += "+"
	}
	for _, v := range optSet.options {
		if v.posixOpt != 0 {
			optstring += string(v.posixOpt)
			if !v.noArg {
				optstring += ":"
			}
		}
	}
	content, err := Expand(args[1:], optstring)
	if err != nil {
		optSet.done, err = optSet.errorHandler(err, Option{"", &args[0]})
	}
	positional := make([]string, 0)
	for _, opt := range content {
		if opt.Opt != "" {
			if item, found := optSet.options[opt.Opt]; found == true {
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
			optSet.done, err = optSet.errorHandler(err, opt)
		}
	}
	return make([]string, 0), err
}

func (optSet OptSet) Done() bool {
	return optSet.done
}

func (optSet *OptSet) ArgFunc(flag rune, longFlag string, action func(string) error, help string) error {
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
	}
	def.argConv = action
	return optSet.safeAdd(def)
}

func (optSet *OptSet) FlagFunc(flag rune, longFlag string, action func() error, help string) error {
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		noArg:    true,
	}
	def.argConv = func(string) error {
		return action()
	}
	return optSet.safeAdd(def)
}

func (optSet *OptSet) Flag(flag rune, longFlag string, help string) (*bool, error) {
	var result bool
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		noArg:    true,
	}
	def.argConv = func(arg string) error {
		result = true
		def.count++
		return nil
	}
	def.argReset = func() {
		result = false
	}
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) StringValue(flag rune, longFlag string, required bool, help string) (*string, error) {
	var result string
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		required: required,
		help:     help,
	}
	def.argConv = func(arg string) error {
		result = arg
		def.count++
		return nil
	}
	def.argReset = func() {
		result = ""
	}
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) StringDefault(flag rune, longFlag string, value string, help string) (*string, error) {
	result := value
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
	}
	def.argConv = func(arg string) error {
		result = arg
		def.count++
		return nil
	}
	def.argReset = func() {
		result = value
	}
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) StringList(flag rune, longFlag string, help string) (*[]string, error) {
	result := make([]string, 0)
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
	}
	def.argConv = func(arg string) error {
		result = append(result, arg)
		def.count++
		return nil
	}
	def.argReset = func() {
		result = make([]string, 0)
	}
	return &result, optSet.safeAdd(def)
}

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
	return strconv.ParseInt(arg, 10, 64)
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
	return strconv.ParseUint(arg, 10, 64)
}

func (optSet *OptSet) IntValue(flag rune, longFlag string, required bool, help string) (*int64, error) {
	var result int64
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		required: required,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) IntDefault(flag rune, longFlag string, value int64, help string) (*int64, error) {
	result := value
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) IntList(flag rune, longFlag string, help string) (*[]int64, error) {
	result := make([]int64, 0)
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) UintValue(flag rune, longFlag string, required bool, help string) (*uint64, error) {
	var result uint64
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		required: required,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) UintDefault(flag rune, longFlag string, value uint64, help string) (*uint64, error) {
	result := value
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) UintList(flag rune, longFlag string, help string) (*[]uint64, error) {
	result := make([]uint64, 0)
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) FloatValue(flag rune, longFlag string, required bool, help string) (*float64, error) {
	var result float64
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		required: required,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) FloatDefault(flag rune, longFlag string, value float64, help string) (*float64, error) {
	result := value
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) FloatList(flag rune, longFlag string, help string) (*[]float64, error) {
	result := make([]float64, 0)
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) BoolValue(flag rune, longFlag string, required bool, help string) (*bool, error) {
	var result bool
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
		required: required,
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
	return &result, optSet.safeAdd(def)
}

func (optSet *OptSet) BoolDefault(flag rune, longFlag string, value bool, help string) (*bool, error) {
	result := value
	def := optDef{
		posixOpt: flag,
		longOpt:  longFlag,
		help:     help,
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
	return &result, optSet.safeAdd(def)
}

func (optSet OptSet) safeAdd(def optDef) error {
	if def.posixOpt != 0 {
		if err := optSet.safeAddKey(string("-"+string(def.posixOpt)), def); err != nil {
			return err
		}
	}
	if def.longOpt != "" {
		if err := optSet.safeAddKey(def.longOpt, def); err != nil {
			return err
		}
	}
	return nil
}

func (optSet *OptSet) safeAddKey(option string, opt optDef) error {
	if val, found := optSet.options[option]; found {
		return errors.New("Duplicate options key: " + option + ": " + val.help + " & " + opt.help)
	} else {
		optSet.options[option] = opt
	}
	return nil
}

func (optSet *OptSet) Help() error {
	optSet.done = true
	return nil
}

func (optSet *OptSet) Version() error {
	optSet.done = true
	return nil
}

func (optSet *OptSet) SetName(name string) {
	optSet.name = name
}

func (optSet *OptSet) SetVersion(version string) {
	optSet.version = version
}

func (optSet *OptSet) SetDescription(description []string) {
	optSet.description = description
}
