# go-getopt

One more getopt implementation for golang.

- tokenizer supports standard getopt specifications string, including modes of:
  - posixly correct
  - parse positionals as arguments
  - don't trport errors to stderr
- high-level option configuration 
  - follows POSIX standard to print help or version to STDOUT
  - gives more flexibility on what to pring in help

## USAGE

### As tokenizer

````
  if opts, err := getopt.Expand(os.Args(), "hVi:o:f::p") ; err == nil {
    for opt := range opts {
      switch opt.Opt {
        case "-h": help()
        case "-i": parseInpot(opt.Arg())
        ...
      }
    }
  }
````

where opt is []interface{Opt() string, Arg() *string}

Opt contains option in form "-l" or "-x" (even if merged form of "-lx" 
was used in argument list).  Arg may be nil for flags.  For positional
arguments opt is empty.  Long options are reported as positional parameters
(same as regular getopt would behave).

### Rich form

````
  opts := getopt.NewOptSet()
  opts.AddDefaults("myProg", "v1.0", []string{"the description"})
  opts.SetErrorHandler(...)
  input, _ := opts.StringValue('i', "--input", true, "Input file name")
  output, _ := opts.StringDefault('o', "--output", "/dev/stdout", "Output file name")
  dbs, _ := opts.StringList('d', "--db", "List of databases to query")
  if (args, err := opts.Parse(os.Args(), false); err == nil {
      if opts.Done() {
          os.Exit(0) // There was help or version used - EX_OK
      }
      for _, arg : range args {
          ....
      }
  } else {
      os.Exit(64) // EX_USAGE
  }
  ...
````

Supported following configurators: 
 - ArgFunc(flag rune, longopt string, action func(string) error, help string) error 
 - FlagFunc(opt rune, longopt string, action func() error, help string) error
 - Flag(opt rune, longopt, help []string) (*bool, error)
 - <type>Value(opt rune, longopt string, required bool, help []string) (*<type>, error)
 - <type>Default(opt rune, longopt string, defaultValue <type>, help []string) (*<type>, error)
 - <type>List(opt rune, longopt string, help []string) (*[]<type>, error)

Where <type> is one of:
  - String
  - Int
  - Uint
  - Float
  - Bool

For unsigned integers, input format supports 
 - decimals, 
 - octals /0[0-7]+/, 
 - hexadecimals /0x[0-9a-fA-F]+/, 
 - binary /0b[01]+/, 
 - base32 /0t[0-9A-Va-v]*/,
 - base64 /0s[0-9a-zA-Z/+]={0,2}/.

Iht, Uint, and Float values are parsed as 64bit.

# EOF
