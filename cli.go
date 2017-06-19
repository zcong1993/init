package main

import (
	"io"
	"flag"
	"fmt"
	"github.com/mitchellh/colorstring"
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

func (cli *CLI) Run(args []string) int {
	var (
		version bool
		help bool
		install bool
	)
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&help, "help", false, "")
	flags.BoolVar(&help, "h", false, "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	flags.BoolVar(&install, "install", false, "")
	flags.BoolVar(&install, "i", false, "")

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
	gc, err := normalizeUrl(url)
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
	if _, err := os.Stat(configFile); err == nil {
		d, err := GetConfig(configFile)
		if err != nil {
			PrintRedf(cli.errStream,
				err.Error())
			return ExitCodeError
		}
		data = d
	}
	fmt.Printf("%+v", data)
	err1 := CopyDir(filepath.Join(dest, "template"), outPut, data)
	if err1 != nil {
		PrintRedf(cli.errStream,
			err.Error())
		return ExitCodeError
	}
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
	Usage: init [options...] REPO
rls is a tool to create Release on Github.

`
