package main

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"fmt"
	"github.com/zcong1993/utils"
)

// inspired by https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04

// CopyFileWithData can compile src file with data to dst
func CopyFileWithData(src, dst string, data map[string]interface{}) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	tpl, err := ioutil.ReadFile(src)
	if err != nil {
		return
	}
	err = utils.Compile(out, string(tpl), data)
	return
}

// CopyDirWithData can compile all src folder files with data to dst
func CopyDirWithData(src string, dst string, data map[string]interface{}) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination path already exists, delete it first or use option '-force, -f'")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirWithData(srcPath, dstPath, data)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFileWithData(srcPath, dstPath, data)
			if err != nil {
				return
			}
		}
	}

	return
}
