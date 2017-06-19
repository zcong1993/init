package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/colorstring"
	"github.com/xtaci/goeval"
	"io"
	"os"
	"path/filepath"
)

const (
	// ExitCodeOK is exit code 0
	ExitCodeOK int = 0
	// ExitCodeError is exit code 1
	ExitCodeError = 1
)

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
		help    bool
		install bool
		force   bool
	)
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&help, "help", false, "")
	flags.BoolVar(&help, "h", false, "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	flags.BoolVar(&install, "install", false, "")
	flags.BoolVar(&install, "i", false, "")

	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}
	if help {
		fmt.Fprint(cli.errStream, helpText)
		return 0
	}
	if version {
		ShowVersion()
		return 0
	}
	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		PrintRedf(cli.errStream,
			"Invalid argument: you use as 'init [option] REPO DESTPATH'")
		return ExitCodeError
	}
	url := parsedArgs[0]
	outPut := parsedArgs[1]
	gc, err := normalizeURL(url)
	if err != nil {
		PrintRedf(cli.errStream,
			err.details)
		return err.errorCode
	}
	dest := filepath.Join(TemplateHome, gc.Repo.Owner, gc.Repo.Name)
	if install {
		dest, err = DownloadAndExtract(gc)
		if err != nil {
			PrintRedf(cli.errStream,
				err.details)
			return err.errorCode
		}
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		PrintRedf(cli.errStream,
			"Cache not exists, please run with option '-install' or '-i'")
		return ExitCodeError
	}
	data := map[string]interface{}{}
	configFile := filepath.Join(dest, "init.json")
	cfg := Cfg{}
	if _, err := os.Stat(configFile); err == nil {
		c, err := NewConfig(configFile)
		if err != nil {
			PrintRedf(cli.errStream,
				err.Error())
			return ExitCodeError
		}
		d, err := c.GetPrompts()
		if err != nil {
			PrintRedf(cli.errStream,
				err.Error())
			return ExitCodeError
		}
		data = d
		cfg = *c
	}
	if force {
		if _, err := os.Stat(outPut); err == nil {
			err := os.RemoveAll(outPut)
			if err != nil {
				PrintRedf(cli.errStream,
					err.Error())
				return ExitCodeError
			}
		}
	}
	src := filepath.Join(dest, "template")
	if _, err := os.Stat(src); os.IsNotExist(err) {
		PrintRedf(cli.errStream,
			"Repo not have template folder, is not a init template repo")
		return ExitCodeError
	}
	sandbox := goeval.Scope{}
	if len(cfg.Config.Filters) != 0 {
		s := *EvalWithVals(data)
		sandbox = s
	}
	err1 := CopyDirWithData(src, outPut, data, &cfg, &sandbox, src)
	if err1 != nil {
		PrintRedf(cli.errStream,
			err1.Error())
		return ExitCodeError
	}
	PrintBluef(cli.outStream, "Success, all done!")
	return ExitCodeOK
}

// PrintRedf is helper function printf with color red
func PrintRedf(w io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]", format)
	fmt.Fprint(w,
		colorstring.Color(fmt.Sprintf(format, args...)))
	fmt.Println()
}

// PrintBluef is helper function printf with color blue
func PrintBluef(w io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[blue]%s[reset]", format)
	fmt.Fprint(w,
		colorstring.Color(fmt.Sprintf(format, args...)))
	fmt.Println()
}

var helpText = `
	Usage: init [options...] REPO DSTPATH
init is a tool to help you init project from git repo template, it will generate custom template by inquiring some question.

 Usage:

 	init [options] REPO DSTPATH

 Options:
 	-install, -i				Install repo from Github, will replace local cache

 	-force, -f					Replace dest path if exists

 Example:

 	init -i -f zcong1993/test ./test

`
