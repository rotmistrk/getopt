# getopt

## SYMOPSIS

### Tokenizer

```go
getopt.Tokenize(os.Args(), "hVi:o:f::p") ([]getopt.Option, error)
```

### Parser

```go
opts := getopt.New().WithDefaults("MyProg", "v1.0", "Usage: [opts] [file...]", ...)
nums := opts.IntList('n', "--num", "number of iterations")
...
positional, err := opts.Parse(os.Args, true)
```

### Marshaller

```go
type myStruct struct {
Names []string  `flag:"n,names" help:"name (may be repeated"`
Show func (arg string) error `flas:"S,show" help:"show status for arg"`
...
}

ms := myStruct{}
opts := getopt.New().WithDefaults("MyProg", "v1.0", "Usage: [opts] [file...]", ...)
positional, err := opts.Marshal(&ms, os.Args, true)
```

## DESCRIPTION

One more getopt implementation for golang.

- tokenizer supports standard getopt specifications string, including modes of:
    - posixly correct
    - parse positionals as arguments
    - suppress error reporting to stderr
    - allows flag concatenations (like ls -alt or tail -n100)
- high-level option configuration
    - follows POSIX standard to print help or version to STDOUT
    - gives more flexibility on what to pring in help
    - requires two dashes and = for long options
- high level struct marshalling
    - list support
    - callback function support

## USAGE

### As tokenizer

````
  if opts, err := getopt.Tokenize(os.Args(), "hVi:o:f::p") ; err == nil {
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
was used in argument list). Arg may be nil for flags. For positional arguments opt is empty.

Long options are reported as positional parameters
(same as regular getopt would behave).

### Rich form

````
  opts := getopt.New().WithDefaults("MyProg","v1.0","Usage: [opts] [file...]", ...)
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

Highlights

- AddDefaults adds -h, -V, --help, and --version flags
    - Help is auto-generated, uses description provided as header
    - Help and Version are reported to stdout
- check opts.Done() to exit on errors or Help/Version request

Supported following configurators:

- ArgFunc(flag rune, longopt string, action func(string) error, help string) error
- FlagFunc(opt rune, longopt string, action func() error, help string) error
- Flag(opt rune, longopt, help []string) (*bool, error)
- &lt;type>Value(opt rune, longopt string, required bool, help []string) (*&lt;type>, error)
- &lt;type>Default(opt rune, longopt string, defaultValue &lt;type>, help []string) (*&lt;type>, error)
- &lt;type>List(opt rune, longopt string, help []string) (*[]&lt;type>, error)

Where <type> is one of:

- String // *string or *[]string
- Int // *int64 or *[]int64
- Uint // *uint64 or *[]uint64
- Float // *float64 or *[]float64
- Bool // *bool

Each function also has variant

- &lt;type>&lt;variant>V(flags[]rune, longopts[]string, ...)
  to allow synonyms; for example
- StringListV(flags[]rune, longopts[]string, help string) (*[]string, error)

For unsigned integers, input format supports

- decimals,
- octals /0[0-7]+/,
- hexadecimals /0x[0-9a-fA-F]+/,
- binary /0b[01]+/,
- base32 /0t[0-9A-Va-v]*/,
- base64 /0s[0-9a-zA-Z/+]={0,2}/.

### Marshaller

For a struct passes by pointer, for public fields that have annotation "flag", sets up parser with all short and long
flahs. One-letter become short, rest become long. Bool are flags only (no args supported).

Supports same types as rich form does, plus:

- int
- uint
- float32
- time.Time // in RFC3339
- time.Duration
- func () error // flag callback, has to be not nil
- func (val string) error // flag with arg callback, has to be not nil

In addition to scalar and vector (repeatable) types,

- map[string]&lt;type> are supported (in form like --map-item=key:value -Mkey1:value1)

Example:

```go
type testMarshal struct {
    Flag      bool               `flag:"b,boolean" help:"boolean value"`
    Str       string             `flag:"s,str" help:"string value"`
    StrList   []string           `flag:"S,str-list" help:"string list"`
    StrMap    map[string]string  `flag:"M,str-map" help:"string map (-Mkey:val -M k:v --str-map=ky:vl)"`
    Int       int64              `flag:"i,int-val" help:"integer value"`
    IntList   []int64            `flag:"I,int-list" help:"integer list"`
    IntMap    map[string]int     `flag:"int-map" help:"integer map"`
    Uint      uint64             `flag:"u,uint-val" help:"unsigned int value"`
    UintList  []uint64           `flag:"U,uint-list" help:"unsigned int list"`
    UintMap   map[string]int     `flag:"uint-map" help:"unsigned integer map"`
    Float     float64            `flag:"f,float-val" help:"float value"`
    FloatList []float64          `flag:"F,float-list" help:"float list"`
    FloatMap  map[string]float64 `flag:"float-map" help:"float map"`
    Wait      time.Duration      `flag:"d,duration-val" help:"duration value"`
    WaitList  []time.Duration    `flag:"D,duration-list" help:"duration list"`
    Exec      func (string) error `flag:"x,exec" help:"execute cmd"`
    Show      func() error       `flag:"show" help:"show stats"`
}
```

All arguments are treated as optional.  
Fields initialized prior to call are not changed if flag did not appear in command line.

# EOF
