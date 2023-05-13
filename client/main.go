package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cmd := exec.Command("qtdeploy", "build", "desktop", "asset/cfc-proxyNet")
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
	err = run()
	if err != nil {
		panic(err)
	}
}

func CopyFile(from, to string) error {
	var err error
	f, err := os.Stat(from)
	if err != nil {
		return err
	}
	fn := func(fromFile string) error {
		rel, err := filepath.Rel(from, fromFile)
		if err != nil {
			return err
		}
		toFile := filepath.Join(to, rel)
		if err = os.MkdirAll(filepath.Dir(toFile), 0777); err != nil {
			return err
		}
		file, err := os.Open(fromFile)
		if err != nil {
			return err
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		out, err := os.Create(toFile)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, bufReader)
		return err
	}
	pwd, _ := os.Getwd()
	if !filepath.IsAbs(from) {
		from = filepath.Join(pwd, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}
	if f.IsDir() {
		return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return fn(path)
			} else {
				if err = os.MkdirAll(path, 0777); err != nil {
					return err
				}
			}
			return err
		})
	} else {
		return fn(from)
	}
}
