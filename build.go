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
)

var (
	buildTags = ""
	testTags  = buildTags + " test"
)

var (
	debug         = flag.Bool("debug", false, "Enable symbol table/DWARF generation and disable optimisations/inlining")
	race          = flag.Bool("race", false, "Enable data race detection in the final binary")
	userBuildTags = flag.String("tags", "", "Additional build tags")
	userTestTags  = flag.String("test-tags", "", "Additional test build tags")

	buildEnabled    = flag.Bool("build", false, "Enable building the final binary")
	clearEnabled    = flag.Bool("clear", false, "If set then the console will be cleared before running the build")
	coverEnabled    = flag.Bool("cover", false, "Generates an HTML cover report that's opened in the browser if watch is disabled")
	generateEnabled = flag.Bool("generate", false, "Enable generating code before build")
	lintEnabled     = flag.Bool("lint", false, "Enable linting before build")
	testEnabled     = flag.Bool("test", false, "Enable tests before build")
	watchEnabled    = flag.Bool("watch", false, "Watches for changes and re-runs the build if changes are detected")
	watchExts       = flag.String("watch-exts", ".go .json .sql", "A space separated list of file extensions to watch")
	watchInterval   = flag.Duration("watch-interval", 2*time.Second, "The interval that watch mode checks for file changes")
)

func main() {
	flag.Parse()

	if tags := strings.TrimSpace(*userBuildTags); tags != "" {
		buildTags += " " + tags
	}

	if tags := strings.TrimSpace(*userTestTags); tags != "" {
		testTags += " " + tags
	}

	buildTags = strings.TrimSpace(buildTags)
	testTags = strings.TrimSpace(testTags)

	*watchExts = strings.TrimSpace(*watchExts)

	if !*buildEnabled && !*generateEnabled && !*lintEnabled && !*testEnabled && !*coverEnabled {
		*generateEnabled = true
		*lintEnabled = true
		*testEnabled = true
		*buildEnabled = true
	}

	pkg := flag.Arg(0)

	// Default to building all commands
	if pkg == "" {
		pkg = "./cmd/..."
	}

	// Always only build local commands
	if !strings.HasPrefix(pkg, "./") {
		fmt.Printf("Please build a local package by prefixing the package name with %q\n", "./")

		os.Exit(1)
	}

	run(pkg)

	if *watchEnabled {
		exts := make(map[string]struct{})
		for _, ext := range strings.Split(*watchExts, " ") {
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
				if d.IsDir() {
					return nil
				}

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

			time.Sleep(*watchInterval)
		}
	}
}

func run(pkg string) {
	if *clearEnabled {
		clear()
	}

	fmt.Printf("Running build @ %v:\n", time.Now().Format(time.UnixDate))

	if *generateEnabled {
		if err := generate(); err != nil {
			if !*watchEnabled {
				os.Exit(1)
			}

			return
		}
	}

	if *lintEnabled {
		if err := lint(); err != nil {
			if !*watchEnabled {
				os.Exit(1)
			}

			return
		}
	}

	if *testEnabled && !*coverEnabled {
		if err := test(); err != nil {
			if !*watchEnabled {
				os.Exit(1)
			}

			return
		}
	}

	if *coverEnabled {
		if err := cover(); err != nil {
			if !*watchEnabled {
				os.Exit(1)
			}

			return
		}
	}

	if *buildEnabled {
		if err := build(pkg); err != nil {
			if !*watchEnabled {
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

func generate() error {
	fmt.Print("    go generate... ")

	out, err := exec.Command("go", "generate", "./...").CombinedOutput()
	if err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")
	}

	return nil
}

func test() error {
	fmt.Print("    go test... ")

	out, err := exec.Command("go", "test", "-race", "-tags", testTags, "./...").CombinedOutput()
	if err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")
	}

	return nil
}

func cover() error {
	fmt.Print("    go test (cover)... ")

	out, err := exec.Command("go", "test", "-race", "-tags", testTags, "-coverprofile", "_cover.out", "./...").CombinedOutput()
	if err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Println("ok")
	}

	if !*watchEnabled {
		fmt.Print("    go tool cover... ")

		out, err = exec.Command("go", "tool", "cover", "-html", "_cover.out").CombinedOutput()
		if err != nil {
			fmt.Println("error")

			if len(out) > 0 {
				fmt.Println(string(out))
			}

			return err
		} else {
			fmt.Println("ok")
		}
	}

	return nil
}

func lint() error {
	fmt.Print("    golangci-lint... ")

	if lint, err := exec.LookPath("golangci-lint"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println("not found, please install")
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	} else {
		out, err := exec.Command(lint, "run", "--color", "always", "--build-tags", testTags).CombinedOutput()
		if err != nil {
			fmt.Println("error")

			if len(out) > 0 {
				fmt.Println(string(out))
			}

			return err
		} else {
			fmt.Println("ok")
		}
	}

	return nil
}

func build(pkg string) error {
	fmt.Printf("    go build %v... ", strings.TrimSuffix(pkg, "..."))

	cmd := "go"
	args := []string{"build", "-v", "-tags", buildTags}
	gcflags := []string{}
	ldflags := []string{
		fmt.Sprintf("-X 'main.version=%v'", closestTag()),
		fmt.Sprintf("-X 'main.branch=%v'", commitBranch()),
		fmt.Sprintf("-X 'main.commit=%v'", commitHash()),
		fmt.Sprintf("-X 'main.built=%v'", time.Now().UTC().Format(time.RFC3339)),
		fmt.Sprintf("-X 'main.tags=%v'", buildTags),
	}

	if *debug {
		// -N disables all optimisations
		// -l disables inlining
		// See: go tool compile --help
		gcflags = append(gcflags, "all=-N -l")

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

	if *race {
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

	if out, err := exec.Command(cmd, args...).CombinedOutput(); err != nil {
		fmt.Println("error")

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		return err
	} else {
		fmt.Print("ok ")
	}

	var info []string

	if *debug {
		info = append(info, "debug")
	} else {
		info = append(info, "release")
	}

	if *race {
		info = append(info, "race")
	}

	fmt.Printf("(%v)\n", strings.Join(info, "/"))

	return nil
}

func closestTag() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "describe", "--long", "--tags").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	parts := strings.Split(string(out), "-")
	tag := strings.Join(parts[:len(parts)-2], "-")
	commitsAhead, _ := strconv.Atoi(parts[len(parts)-2])
	commitHash := tagCommitHash(tag)

	version := strings.TrimPrefix(strings.TrimSpace(tag), "v")

	var additional []string
	if commitsAhead > 0 {
		noun := "commit"
		if commitsAhead > 1 {
			noun = "commits"
		}

		additional = append(additional, fmt.Sprintf("%v %v ahead of %v", commitsAhead, noun, commitHash))
	}
	if hasUncommittedChanges() {
		additional = append(additional, "built with uncommitted changes")
	}
	if len(additional) > 0 {
		version = fmt.Sprintf("%v (%v)", version, strings.Join(additional, "; "))
	}

	return version
}

func tagCommitHash(tag string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "show-ref", "-d", "--tags", tag).CombinedOutput()
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

func hasUncommittedChanges() bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return false
	}

	out, err := exec.Command(git, "status", "-su").CombinedOutput()
	if err != nil {
		return false
	}

	return len(out) > 0
}

func commitBranch() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "branch", "--show-current").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

func commitHash() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "rev-list", "-1", "HEAD").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}
