package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var version = "1.0.0"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type config struct {
	dryRun  bool
	verbose bool
	dir     string
	target  string
	ignore  arrayFlags
}

// wrapper struct to hold information about source file
type sourceFile struct {
	fileInfo os.FileInfo
	fileName string
	path     string
}

func (s *sourceFile) DestSymlink(cfg *config) string {
	return path.Join(cfg.target, s.fileName)
}

// Walk the source directory tree, returning the files that will need to be
// symlinked
func readSourceDir(dir string) ([]sourceFile, error) {
	var files []sourceFile
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, f := range fileInfo {
		fullPath := path.Join(dir, f.Name())
		files = append(files, sourceFile{path: fullPath, fileName: f.Name(), fileInfo: f})
	}
	return files, nil
}

// This will read a directory and only return the symlinks
func symlinks(dir string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, f := range fileInfo {
		fullPath := path.Join(dir, f.Name())
		fi, err := os.Lstat(fullPath)
		if err != nil {
			log.Println(err)
			continue
		}
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			files = append(files, fullPath)
		}
	}
	return files, nil
}

func stow(cfg *config) error {
	files, err := readSourceDir(cfg.dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		log.Printf("stow: file=%s symlink=%s", file.path, file.DestSymlink(cfg))
	}
	return nil
}

func destow(cfg *config) error {
	links, err := symlinks(cfg.target)
	if err != nil {
		return err
	}
	for _, link := range links {
		log.Printf("destow: symlink=%s\n", link)
	}
	return nil
}

func main() {
	var (
		flagVersion bool
		flagDestow  bool
	)
	cfg := &config{}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gstow version %s\n\nOPTIONS:\n\n", version)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n\nReport bugs to: https://github.com/sleach/gstow")
	}
	flag.BoolVar(&flagVersion, "version", false, "Print version information and exit")
	flag.BoolVar(&flagDestow, "D", false, "Unstow files, deleting symlinks")
	flag.BoolVar(&cfg.verbose, "v", false, "Print extended information")
	flag.BoolVar(&cfg.dryRun, "n", false, "Dry-run - combined with -v prints out actions but does not run them")
	flag.StringVar(&cfg.dir, "d", "", "Source directory to read files")
	flag.StringVar(&cfg.target, "t", "", "Target directory for where symlinks will be created")
	flag.Var(&cfg.ignore, "ignore", "Regular express match for files to ignore (can be multiple)")
	flag.Parse()

	if flagVersion {
		log.Fatal("gstow version 1.0.0")
	}

	if cfg.dir == "" {
		log.Fatal("-d is a required flag")
	}
	fmt.Println(cfg.target)
	if cfg.target == "" {
		// this defaults to the parent directory of where this is being ran
		wd, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		cfg.target = filepath.Dir(wd)
	}

	if flagDestow {
		err := destow(cfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := stow(cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
