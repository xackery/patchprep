package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"
	ignore "github.com/sabhiram/go-gitignore"
)

var (
	// Version is what version is running
	Version string
)

const (
	contents = `# add .gitignore like patterns
*.ini
.git
.gitignore
patchprep.txt
patchprep.exe
patchprep
.vscode/
uifiles/
uiresources/
_test_data/
voice/
*.txt
!*_chr.txt
!eqhost.txt
*.bak
*.blend
userdata/
ts/
PlayerStudio/
Logs/
LayoutConverter/
Help/
EnvEmitterEffects/
Broon/
.DS_Store
`
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed patchprep:", err)
		os.Exit(1)
	}
}

func run() error {
	if Version == "" {
		Version = "0.0.0"
	}
	fmt.Printf("patchprep v%s\n", Version)

	fi, err := os.Stat("patchprep.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat patchprep.txt: %w", err)
		}
		err = os.WriteFile("patchprep.txt", []byte(contents), 0644)
		if err != nil {
			return fmt.Errorf("create patchprep.txt: %w", err)
		}
		fi, err = os.Stat("patchprep.txt")
		if err != nil {
			return fmt.Errorf("stat patchprep.txt: %w", err)
		}
	}
	if fi.IsDir() {
		return fmt.Errorf("patchprep.txt is a directory")
	}

	r, err := os.Open("patchprep.txt")
	if err != nil {
		return fmt.Errorf("open patchprep.txt: %w", err)
	}
	defer r.Close()

	path := "patch/"
	scanner := bufio.NewScanner(r)
	lines := []string{}
	lines = append(lines, "patch/")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		lines = append(lines, line)
	}

	fi, err = os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", path, err)
		}
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("mkdir %s: %w", path, err)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", path, err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	if path == "patch/" {
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}

	ignores := ignore.CompileIgnoreLines(lines...)

	err = filepath.WalkDir(".", func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ignores.MatchesPath(p) {
			return nil
		}
		// check if create date was within last 2 years
		fi, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("stat %s: %w", p, err)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		t, err := times.Stat(p)
		if err != nil {
			return fmt.Errorf("times.Stat %s: %w", p, err)
		}
		if !t.HasBirthTime() {
			return fmt.Errorf("no birth time for %s", p)
		}
		if time.Now().AddDate(-2, 0, 0).After(t.BirthTime()) {
			return nil
		}

		err = os.MkdirAll(path+filepath.Dir(p), 0755)
		if err != nil {
			return fmt.Errorf("mkdir %s: %w", path+filepath.Dir(p), err)
		}

		r, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("open %s: %w", p, err)
		}
		defer r.Close()

		w, err := os.Create(path + p)
		if err != nil {
			return fmt.Errorf("create %s: %w", path+p, err)
		}
		defer w.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("copy %s: %w", p, err)
		}

		// retain datestamps
		err = os.Chtimes(path+p, fi.ModTime(), fi.ModTime())
		if err != nil {
			return fmt.Errorf("chtimes %s: %w", path+p, err)
		}

		err = os.WriteFile(path+p, []byte{}, 0644)
		if err != nil {
			return fmt.Errorf("write %s: %w", path+p, err)
		}
		fmt.Println(path + p)

		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	return nil
}
