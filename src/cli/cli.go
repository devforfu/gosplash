package cli

import (
    "../imutil"
    "flag"
    "fmt"
    "github.com/mitchellh/go-homedir"
    "log"
    "os"
    "path"
    "strconv"
    "strings"
)

type ParameterType uint32

const (
    String ParameterType = iota
    Integer
    Boolean
    Float
)

type Command struct {
    Parameters []*Parameter
}

func (c *Command) Parse() Params {
    params := make(Params)
    for _, p := range c.Parameters {
        switch p.Type {
        case String:
            p.Value = flag.String(p.Name, p.Default.(string), p.Usage)
        case Integer:
            p.Value = flag.Int(p.Name, p.Default.(int), p.Usage)
        case Boolean:
            p.Value = flag.Bool(p.Name, p.Default.(bool), p.Usage)
        case Float:
            p.Value = flag.Float64(p.Name, p.Default.(float64), p.Usage)
        default:
            panic(fmt.Sprintf("Unknown parameter type: %v", p.Type))
        }
        params[strings.TrimLeft(p.Name, "-")] = p
    }
    flag.Parse()
    for _, p := range params {
        if p.Value == nil {
            p.Value = p.Default
        }
    }
    return params
}

type Parameter struct {
    Name, Usage string
    Type ParameterType
    Value interface{}
    Default interface{}
}

func (p *Parameter) TypeName() string {
    switch p.Type {
    case String: return "String"
    case Integer: return "Integer"
    case Boolean: return "Boolean"
    case Float: return "Float"
    }
    return "Unknown"
}

func (p *Parameter) String() string {
    if p.Type == String  { return *p.Value.(*string) }
    if p.Type == Integer { return strconv.Itoa(*p.Value.(*int)) }
    panic(fmt.Sprintf("cannot convert into string a parameter of type %s", p.TypeName()))
}

func (p *Parameter) Integer() int {
    if p.Type == Integer { return *p.Value.(*int) }
    panic(fmt.Sprintf("cannot convert into integer a parameter of type %s", p.TypeName()))
}

func (p *Parameter) Boolean() bool {
    if p.Type == Boolean { return *p.Value.(*bool) }
    panic(fmt.Sprintf("cannot convert into boolean a parameter of type %s", p.TypeName()))
}

func (p *Parameter) Float() float64 {
    if p.Type == Integer { return float64(*p.Value.(*int)) }
    if p.Type == Float { return *p.Value.(*float64) }
    panic(fmt.Sprintf("cannot convert into float a parameter of type %s", p.TypeName()))
}

type Params map[string]*Parameter

// -----------------------
// Sub-commands definition
// -----------------------

func Parse(args []string) Params {
    commands := map[string]*Command {
        "download": {
            Parameters: []*Parameter{
                {Name: "-conf", Type: String, Default: "unsplash.key.json", Usage: "path to the file with Unsplash API keys"},
                {Name: "-n", Type: Integer, Default: 5, Usage: "number of images to download"},
                {Name: "-out", Type: String, Default: path.Join(mustHomeDir(), "Unsplash"), Usage: "path to the output folder"},
            },
        },
        "thumb": {
            Parameters: []*Parameter{
                {Name:"-dir", Type:String, Default:mustWorkDir(), Usage:"path to a folder with images"},
                {Name:"-p", Type:String, Default:imutil.ImageFormats, Usage:"image extensions pattern"},
            },
        },
        "canvas": {
            Parameters: []*Parameter{
                {Name:"-dir", Type:String, Default:mustWorkDir(), Usage:"path to a folder with images"},
            },
        },
    }
    return commands[args[1]].Parse()
}

func mustHomeDir() string {
    home, err := homedir.Dir()
    if err != nil { log.Fatal(err) }
    return home
}

func mustWorkDir() string {
    workdir, err := os.Getwd()
    if err != nil { log.Fatal(err) }
    return workdir
}