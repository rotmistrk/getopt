package getopt

import (
	"errors"
	"strings"
)

type Option struct {
	Opt string
	Arg *string
}

func (option Option) String() string {
	if option.Opt != "" {
		if option.Arg != nil {
			if len(option.Opt) > 2 {
				return option.Opt + "='" + *option.Arg + "'"
			} else {
				return option.Opt + " '" + *option.Arg + "'"
			}
		} else {
			return option.Opt
		}
	} else {
		return "'" + *option.Arg + "'"
	}
}

type getoptConfig struct {
	dontPrintErrors  bool
	positionalAsArgs bool
	posixlyCorrect   bool
	optmap           map[rune]int
}

func Expand(args []string, options string) ([]Option, error) {
	if cfg, err := newGetoptConfig(options); err != nil {
		return nil, err
	} else {
		result := make([]Option, 0)
		tail := make([]Option, 0)
		nextOpt := rune(0)
		nextOptIsOpt := false
		for pos, arg := range args[1:] {
			if nextOpt != 0 && nextOptIsOpt == false {
				sarg := arg
				result = append(result, Option{"-" + string(nextOpt), &sarg})
				nextOpt = 0
			} else {
				l := len(arg)
				if l > 0 && arg[0] == '-' { // -
					if nextOpt != 0 {
						result = append(result, Option{"-" + string(nextOpt), nil})
						nextOpt = 0
						nextOptIsOpt = false
					}
					if l > 1 && arg[1] == '-' { // --
						if l > 2 { // --flag
							if eq := strings.Index(arg, "="); eq < 0 {
								result = append(result, Option{arg, nil})
							} else {
								sarg := string(arg[eq+1:])
								result = append(result, Option{arg[0:eq], &sarg})
							}
						} else { // --
							for p := pos + 2; p < len(args); p++ {
								sarg := args[p]
								result = append(result, Option{"", &sarg})
							}
							break
						}
					} else { // -a
						argrunes := []rune(arg)[1:]
						for chpos, ch := range argrunes {
							if optType := cfg.optmap[ch]; optType == 0 {
								result = append(result, Option{"-" + string(ch), nil})
							} else {
								if chpos+1 == len(argrunes) {
									nextOpt = ch
									if optType > 1 {
										nextOptIsOpt = true
									}
								} else {
									sarg := string(argrunes[chpos+1:])
									result = append(result, Option{"-" + string(ch), &sarg})
									break
								}
							}
						}
					}
				} else if cfg.posixlyCorrect {
					for p := pos + 1; p < len(args); p++ {
						sarg := args[p]
						result = append(result, Option{"", &sarg})
					}
					break
				} else {
					sarg := arg
					tail = append(tail, Option{"", &sarg})
				}
			}
		}
		var err error
		if nextOpt != 0 {
			result = append(result, Option{"-" + string(nextOpt), nil})
			if !nextOptIsOpt {
				err = errors.New("Missing argument to required option -" + string(nextOpt))
			}
		}
		return append(result, tail...), err
	}
}

func newGetoptConfig(options string) (getoptConfig, error) {
	cfg := getoptConfig{
		optmap: make(map[rune]int),
	}
	last := rune(0)
	for i, c := range options {
		switch c {
		case ':':
			if last == 0 {
				cfg.dontPrintErrors = true
			} else {
				cfg.optmap[last] += 1
			}
		case '-':
			if i == 0 {
				cfg.positionalAsArgs = true
			}
		case '+':
			if i == 0 {
				cfg.posixlyCorrect = true
			}
		default:
			last = c
			cfg.optmap[last] = 0
		}
	}
	return cfg, nil
}
