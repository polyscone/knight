//go:build ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	tags     = ""
	testTags = "test"
)

var opts struct {
	goos              string
	goarch            string
	tags              string
	testTags          string
	unoptimised       bool
	debug             bool
	race              bool
	build             bool
	clear             bool
	cover             bool
	generate          bool
	lint              bool
	lintTimeout       time.Duration
	test              bool
	verbose           bool
	watch             bool
	watchExts         string
	watchSkipPatterns string
	watchInterval     time.Duration
}

func main() {
	flag.StringVar(&opts.goos, "goos", "", "Sets the GOOS environment variable for the build")
	flag.StringVar(&opts.goarch, "goarch", "", "Sets the GOARCH environment variable for the build")
	flag.StringVar(&opts.tags, "tags", "", "Additional build tags")
	flag.StringVar(&opts.testTags, "test-tags", "", "Additional test build tags")
	flag.BoolVar(&opts.unoptimised, "unoptimised", false, "Disable optimisations/inlining")
	flag.BoolVar(&opts.debug, "debug", false, "Enable symbol table/DWARF generation")
	flag.BoolVar(&opts.race, "race", false, "Enable data race detection in the final binary")
	flag.BoolVar(&opts.build, "build", false, "Enable building the final binary")
	flag.BoolVar(&opts.clear, "clear", false, "If set then the console will be cleared before running the build")
	flag.BoolVar(&opts.cover, "cover", false, "Generates an HTML cover report that's opened in the browser if watch is disabled")
	flag.BoolVar(&opts.generate, "generate", false, "Enable generating code before build")
	flag.BoolVar(&opts.lint, "lint", false, "Enable linting before build")
	flag.DurationVar(&opts.lintTimeout, "lint-timeout", 1*time.Minute, "The timeout duration to pass to the linter")
	flag.BoolVar(&opts.test, "test", false, "Enable tests before build")
	flag.BoolVar(&opts.verbose, "verbose", false, "Print the commands that are being run along with all command output")
	flag.BoolVar(&opts.watch, "watch", false, "Watches for changes and re-runs the build if changes are detected")
	flag.StringVar(&opts.watchExts, "watch-exts", ".go .h .c .sql .json", "A space separated list of file extensions to watch")
	flag.StringVar(&opts.watchSkipPatterns, "watch-skip-patterns", ".git/ .hg/ .svn/ node_modules/ build.go", "A space separated list of patterns to skip in watch mode")
	flag.DurationVar(&opts.watchInterval, "watch-interval", 2*time.Second, "The interval that watch mode checks for file changes")
	flag.Parse()

	if s := strings.TrimSpace(opts.tags); s != "" {
		tags += " " + s
	}
	tags = strings.TrimSpace(tags)

	if s := strings.TrimSpace(opts.testTags); s != "" {
		testTags += " " + s
	}
	testTags = strings.TrimSpace(tags + " " + testTags)

	opts.watchExts = strings.TrimSpace(opts.watchExts)
	opts.watchSkipPatterns = strings.TrimSpace(opts.watchSkipPatterns)

	if !opts.build && !opts.generate && !opts.lint && !opts.test && !opts.cover {
		opts.generate = true
		opts.lint = true
		opts.test = true
		opts.build = true
	}

	pkg := flag.Arg(0)

	// Default to building all commands
	if pkg == "" {
		if fi, err := os.Stat("cmd"); errors.Is(err, fs.ErrNotExist) || !fi.IsDir() {
			pkg = "./..."
		} else {
			pkg = "./cmd/..."
		}
	}

	// Always only build local commands
	if !strings.HasPrefix(pkg, "./") {
		fmt.Printf("Please build a local package by prefixing the package name with %q\n", "./")

		os.Exit(1)
	}

	if opts.goos == "" {
		opts.goos = runtime.GOOS
	}

	if opts.goarch == "" {
		opts.goarch = runtime.GOARCH
	}

	// Always immediately run the build pipeline at least once, even if in watch mode
	run(pkg)

	if opts.watch {
		skipPatterns := strings.Fields(opts.watchSkipPatterns)
		skip := func(path string) bool {
			for _, pattern := range skipPatterns {
				path = filepath.ToSlash(path)

				if path == pattern {
					return true
				}

				if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern) {
					return true
				}
			}

			return false
		}

		exts := make(map[string]struct{})
		for _, ext := range strings.Fields(opts.watchExts) {
			if ext == "" {
				continue
			}

			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}

			exts[ext] = struct{}{}
		}

		files := make(map[string]time.Time)
		for {
			var changed bool

			filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
				// Completely skip directories that are in the skip patterns
				if d.IsDir() && skip(path) {
					return filepath.SkipDir
				}

				// Individually skip directories/files that haven't been entirely skipped by the previous check
				if d.IsDir() || skip(path) {
					return nil
				}

				// Skip any files that don't match watch extensions
				if _, ok := exts[filepath.Ext(path)]; !ok {
					return nil
				}

				info, err := d.Info()
				if err != nil {
					return err
				}

				if modified, ok := files[path]; !changed && ok {
					changed = modified.Before(info.ModTime())
				}

				files[path] = info.ModTime()

				return nil
			})

			if changed {
				run(pkg)
			}

			time.Sleep(opts.watchInterval)
		}
	}
}

func run(pkg string) {
	if opts.clear {
		clear()
	}

	fmt.Printf("Running build @ %v:\n", time.Now().Format(time.UnixDate))

	if opts.generate {
		if err := generate(); err != nil {
			if !opts.watch {
				os.Exit(1)
			}

			return
		}
	}

	if opts.lint {
		if err := lint(); err != nil {
			if !opts.watch {
				os.Exit(1)
			}

			return
		}
	}

	// If cover is enabled then we skip this because cover includes a call to test with cover profile flags anyway
	if opts.test && !opts.cover {
		if err := test(); err != nil {
			if !opts.watch {
				os.Exit(1)
			}

			return
		}
	}

	if opts.cover {
		if err := cover(); err != nil {
			if !opts.watch {
				os.Exit(1)
			}

			return
		}
	}

	if opts.build {
		if err := build(pkg); err != nil {
			if !opts.watch {
				os.Exit(1)
			}

			return
		}
	}
}

func clear() {
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func command(program string, args ...string) (string, []string, string) {
	messageValues := make([]interface{}, len(args))
	for i, arg := range args {
		messageValues[i] = arg
	}

	verbs := make([]string, len(args))
	for i, arg := range args {
		if strings.IndexFunc(arg, unicode.IsSpace) >= 0 {
			verbs[i] = "%q"
		} else {
			verbs[i] = "%v"
		}
	}
	message := fmt.Sprintf("$ %v "+strings.Join(verbs, " "), append([]interface{}{program}, messageValues...)...)

	return program, args, message
}

func generate() error {
	program, args, message := command("go", "generate", "./...")
	if opts.verbose {
		fmt.Println("    " + message)
	}

	fmt.Print("    go generate... ")

	if out, err := exec.Command(program, args...).CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")

		if len(out) > 0 && opts.verbose {
			fmt.Println(string(out))
		}
	}

	return nil
}

func test() error {
	program, args, message := command("go", "test", "-race", "-tags", testTags, "./...")
	if opts.verbose {
		fmt.Println("    " + message)
	}

	fmt.Print("    go test... ")

	if out, err := exec.Command(program, args...).CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")

		if len(out) > 0 && opts.verbose {
			fmt.Println(string(out))
		}
	}

	return nil
}

func cover() error {
	program, args, message := command("go", "test", "-race", "-tags", testTags, "-coverprofile", "_cover.out", "./...")
	if opts.verbose {
		fmt.Println("    " + message)
	}

	fmt.Print("    go test (cover)... ")

	if out, err := exec.Command(program, args...).CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")

		if len(out) > 0 && opts.verbose {
			fmt.Println(string(out))
		}
	}

	if !opts.watch {
		program, args, message := command("go", "tool", "cover", "-html", "_cover.out")
		if opts.verbose {
			fmt.Println("    " + message)
		}

		fmt.Print("    go tool cover... ")

		if out, err := exec.Command(program, args...).CombinedOutput(); err != nil {
			fmt.Println("error")

			if len(out) > 0 {
				fmt.Println(string(out))
			}

			return err
		} else {
			fmt.Println("ok")

			if len(out) > 0 && opts.verbose {
				fmt.Println(string(out))
			}
		}
	}

	return nil
}

func lint() error {
	linter := "golangci-lint"
	program, args, message := command(linter, "run", "--allow-parallel-runners", "--timeout", opts.lintTimeout.String(), "--color", "always", "--build-tags", testTags)
	if opts.verbose {
		fmt.Println("    " + message)
	}

	fmt.Print("    golangci-lint... ")

	if _, err := exec.LookPath(linter); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println("not found, please install")
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	} else if out, err := exec.Command(program, args...).CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")

		if len(out) > 0 && opts.verbose {
			fmt.Println(string(out))
		}
	}

	return nil
}

func build(pkg string) error {
	tagsMessage := tags
	if tagsMessage == "" {
		tagsMessage = "-"
	}

	args := []string{"build", "-v", "-x", "-o", ".", "-tags", tags}
	gcflags := []string{}
	ldflags := []string{
		fmt.Sprintf("-X 'main.branch=%v'", commitBranch("")),
		fmt.Sprintf("-X 'main.version=%v'", closestTag("")),
		fmt.Sprintf("-X 'main.commit=%v'", commitHash("")),
		fmt.Sprintf("-X 'main.tags=%v'", tagsMessage),
	}

	if opts.unoptimised {
		// -N disables all optimisations
		// -l disables inlining
		// See: go tool compile --help
		gcflags = append(gcflags, "all=-N -l")
	}

	if opts.debug {
		if opts.goos == "windows" {
			// This is required on Windows to view disassembly in things like pprof
			args = append(args, "-buildmode", "exe")
		}

		ldflags = append(ldflags, "-X 'main.target=Debug'")
	} else {
		args = append(args, "-trimpath")

		ldflags = append(ldflags, "-X 'main.target=Release'")

		// -s disables the symbol table
		// -w disables DWARF generation
		// See: go tool link --help
		ldflags = append(ldflags, "-s")
		ldflags = append(ldflags, "-w")
	}

	if opts.race {
		args = append(args, "-race")

		ldflags = append(ldflags, "-X 'main.race=Enabled'")
	} else {
		ldflags = append(ldflags, "-X 'main.race=Disabled'")
	}

	if len(gcflags) > 0 {
		args = append(args, "-gcflags", strings.Join(gcflags, " "))
	}

	if len(ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(ldflags, " "))
	}

	args = append(args, pkg)

	var env []string
	if opts.goos != "" {
		env = append(env, "GOOS="+opts.goos)
	}
	if opts.goarch != "" {
		env = append(env, "GOARCH="+opts.goarch)
	}

	program, args, message := command("go", args...)
	if opts.verbose {
		fmt.Println("    " + message)
	}

	fmt.Printf("    go build %v... ", strings.TrimSuffix(pkg, "..."))

	cmd := exec.Command(program, args...)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Print("ok ")

		var info []string

		if opts.debug {
			info = append(info, "debug")
		} else {
			info = append(info, "release")
		}

		if opts.race {
			info = append(info, "race")
		}

		fmt.Printf("(%v)\n", strings.Join(info, "/"))

		if len(out) > 0 && opts.verbose {
			fmt.Println(string(out))
		}
	}

	return nil
}

func closestTag(dir string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	cmd := exec.Command(git, "describe", "--long", "--tags")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	parts := strings.Split(string(out), "-")
	tag := strings.Join(parts[:len(parts)-2], "-")
	commitsAhead, _ := strconv.Atoi(parts[len(parts)-2])
	commitHash := tagCommitHash(dir, tag)
	version := strings.TrimPrefix(strings.TrimSpace(tag), "v")

	var additional []string
	if commitsAhead > 0 {
		noun := "commit"
		if commitsAhead > 1 {
			noun = "commits"
		}

		additional = append(additional, fmt.Sprintf("%v %v ahead of %v", commitsAhead, noun, commitHash))
	}
	if hasUncommittedChanges(dir) {
		additional = append(additional, "built with uncommitted changes")
	}
	if len(additional) > 0 {
		version = fmt.Sprintf("%v (%v)", version, strings.Join(additional, "; "))
	}

	return version
}

func tagCommitHash(dir, tag string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	cmd := exec.Command(git, "show-ref", "-d", "--tags", tag)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, "^{}") {
			parts := strings.Split(line, " ")

			return parts[0][:7]
		}
	}

	return "unknown"
}

func hasUncommittedChanges(dir string) bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return false
	}

	cmd := exec.Command(git, "status", "-su")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return len(out) > 0
}

func commitBranch(dir string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	cmd := exec.Command(git, "branch", "--show-current")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

func commitHash(dir string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	cmd := exec.Command(git, "rev-list", "-1", "HEAD")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}
