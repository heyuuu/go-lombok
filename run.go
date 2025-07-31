package main

import (
	"flag"
	"fmt"
	"github.com/heyuuu/go-lombok/internal/lombok"
	"log"
	"os"
	"path/filepath"
)

func run(args []string) {
	if len(args) < 2 {
		log.Fatalln("Args 不可为空")
	}

	opts := parseOpts(args)
	fmt.Printf("%+v\n", opts)
	switch opts.cmd {
	case "gen", "generate":
		lombok.RunTask(lombok.TaskGenerate, opts.dir)
	case "clear":
		lombok.RunTask(lombok.TaskClear, opts.dir)
	case "":
		log.Fatalln("命令不可为空")
	default:
		log.Fatalln("未定义命令: " + opts.cmd)
	}
}

type optsType struct {
	cmd string
	dir string
}

func parseOpts(args []string) (opts optsType) {
	// options
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.StringVar(&opts.cmd, "cmd", "", "command")
	flagSet.StringVar(&opts.dir, "d", "", "dir")
	_ = flagSet.Parse(args[1:])

	// workdir
	if opts.dir == "" || opts.dir[0] != '/' {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		opts.dir = filepath.Join(wd, opts.dir)
	}

	return
}
