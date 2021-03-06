package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

var (
	version = "-"
	branch  = "-"
	commit  = "-"
	tags    = "-"
	target  = "-"
	race    = "-"
)

var opts struct {
	expression string
	filename   string
	profile    string
	astStyle   string
	version    bool
}

func main() {
	flag.StringVar(&opts.expression, "e", "", "An expression to evaluate")
	flag.StringVar(&opts.filename, "f", "", "A path to a file to run")
	flag.StringVar(&opts.profile, "p", "", "The name of a profile to record")
	flag.StringVar(&opts.astStyle, "a", "", `Print the program's AST; available styles are: "sexpr", "tree", and "waterfall"`)
	flag.BoolVar(&opts.version, "version", false, "Display binary version information")
	flag.Parse()

	if opts.version || flag.Arg(0) == "version" {
		fmt.Println("Version:      ", version)
		fmt.Println("Branch:       ", branch)
		fmt.Println("Commit:       ", commit)
		fmt.Println("Tags:         ", tags)
		fmt.Println("Go version:   ", strings.TrimPrefix(runtime.Version(), "go"))
		fmt.Println("OS/Arch:      ", runtime.GOOS+"/"+runtime.GOARCH)
		fmt.Println("Target:       ", target)
		fmt.Println("Race detector:", race)

		return
	}

	if (opts.expression == "" && opts.filename == "") || (opts.expression != "" && opts.filename != "") {
		flag.Usage()

		os.Exit(2)
	}

	//nolint:gocritic // it's ok for defers to not run in profile code
	switch opts.profile {
	case "cpu":
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	case "mem":
		f, err := os.Create("mem.pprof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		runtime.MemProfileRate = 4096

		//nolint:errcheck // any errors here will be apparent so we don't check
		defer pprof.Lookup("heap").WriteTo(f, 0)
	case "trace":
		f, err := os.Create("trace.out")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := trace.Start(f); err != nil {
			log.Fatal(err)
		}
		defer trace.Stop()
	default:
		if opts.profile != "" {
			log.Fatalf("unknown profile type %q", opts.profile)
		}
	}

	var b []byte
	if opts.expression != "" {
		b = []byte(opts.expression)
	} else {
		var err error
		if b, err = os.ReadFile(opts.filename); err != nil {
			fmt.Println(err)

			os.Exit(1)
		}
	}

	l := lexer.New()
	p := parser.New(l)
	g := value.NewGlobalStore()
	program, err := p.Parse(g, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	if opts.astStyle != "" {
		switch opts.astStyle {
		case "sexpr":
			fmt.Println(program.ASTString(ast.StyleSexpr))
		case "tree":
			fmt.Println(program.ASTString(ast.StyleTree))
		case "waterfall":
			fmt.Println(program.ASTString(ast.StyleWaterfall))
		default:
			flag.Usage()

			os.Exit(2)
		}

		return
	}

	rand.Seed(time.Now().UnixNano())

	if _, err := interpreter.New(g, p).Execute(program); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
