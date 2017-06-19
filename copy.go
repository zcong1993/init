package main

import (
	"fmt"
	"github.com/zcong1993/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/xtaci/goeval"
	"github.com/bmatcuk/doublestar"
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
func CopyDirWithData(src string, dst string, data map[string]interface{}, cfg *Cfg, sandbox *goeval.Scope) (err error) {
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
			err = CopyDirWithData(srcPath, dstPath, data, cfg, sandbox)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if len(cfg.Config.Filters) != 0 {
				for key, val := range cfg.Config.Filters {
					v, err := sandbox.Eval(val)
					if err != nil {
						return err
					}
					fmt.Println(v, v.(bool))
					isMatch, err := doublestar.PathMatch(filepath.Join(src, key), srcPath)
					if err != nil {
						return err
					}
					fmt.Println(filepath.Join(src, key), isMatch, srcPath)
					if isMatch {
						if v.(bool){
							err = CopyFileWithData(srcPath, dstPath, data)
							if err != nil {
								return err
							}
						}
					} else {
						err = CopyFileWithData(srcPath, dstPath, data)
						if err != nil {
							return err
						}
					}
				}
			} else {
				err = CopyFileWithData(srcPath, dstPath, data)
				if err != nil {
					return
				}
			}
		}
	}

	return
}
