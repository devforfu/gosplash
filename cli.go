package main

import (
    "./src/imutil"
    "./src/unsplash"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/mitchellh/go-homedir"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strconv"
)

func main() {
    //cmd := Command{
    //    {Type:String, Name:"x", Usage:"first param", Default:"some"},
    //    {Type:Integer, Name:"y", Usage:"second param", Default:1},
    //}
    //params := cmd.Parse()
    //fmt.Printf("%s\n", params["x"].String())

    //fmt.Printf("%d\n", params["y"].Integer())
    switch cmd, params := Parse(); cmd {
    case "download":
        DownloadImages(params)
    case "thumb":
        MakeThumbnails(params)
    case "concat":
        ConcatIntoSingle(params)
    }
}

func Parse() (string, Params) {
    parser := os.Args[1]
    args := make([]string, 0)
    args = append(args, os.Args[0])
    args = append(args, os.Args[2:]...)
    os.Args = args

    var cmd Command
    switch parser {
    case "download":
        cmd = Command{
            {Type: String, Name: "conf", Default: "unsplash.key.json", Usage: "path to the file with Unsplash API keys"},
            {Type: Integer, Name: "n", Default: 5, Usage: "number of images to download"},
            {Type: String, Name: "out", Default: path.Join(mustHomeDir(), "Unsplash"), Usage: "path to the output folder"},
        }
    default:
        panic("unexpected command: " + parser)
    }

    return parser, cmd.Parse()
}

func DownloadImages(params Params) {
    log.Println("Downloading random images...")
    jsonConfig := MustLoadJSON(params["conf"].String())
    client := unsplash.Client{
        AccessKey:jsonConfig["accessKey"],
        SecretKey:jsonConfig["secretKey"]}
    err := client.DownloadRandomPhotos(
        params["out"].String(),
        params["n"].Integer())
    if err != nil { log.Fatal(err) }
    log.Printf("Images downloaded into folder: %s\n", params["out"].String())
}

func MakeThumbnails(params Params) {
    log.Println("Converting images into thumbnails...")
    maker := imutil.NewThumbnailMaker(imutil.PNG, 128, 128)
    thumbs, err := maker.ConvertFolder(
        params["dir"].String(), params["p"].String())
    if err != nil { log.Fatal(err) }
    log.Println("Created thumbnails:")
    n := len(strconv.Itoa(len(thumbs)))
    for i, t := range thumbs { log.Printf("\t%0*d: %s", n, i, t) }
}

func ConcatIntoSingle(params Params) {
    log.Fatal("Not implemented!")
}

func MustLoadJSON(filename string) map[string]string {
    data, err := ioutil.ReadFile(filename)
    if err != nil { log.Fatal(err) }
    config := make(map[string]string)
    err = json.Unmarshal(data, &config)
    if err != nil { log.Fatal(err) }
    return config
}

type ParameterType uint32

const (
    String ParameterType = iota
    Integer
    Boolean
    Float
)

type Parameter struct {
    Type ParameterType
    Name, Usage string
    Default interface{}
    Value interface{}
}

type Command []*Parameter

type Params map[string]*Parameter

func (c Command) Parse() Params {
    params := make(Params)
    for _, p := range c {
        switch p.Type {
        case String:
            p.Value = flag.String(p.Name, p.Default.(string), p.Usage)
        case Integer:
            p.Value = flag.Int(p.Name, p.Default.(int), p.Usage)
        }
        params[p.Name] = p
    }
    flag.Parse()
    return params
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