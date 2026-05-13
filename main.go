package main

import "github.com/alexflint/go-arg"

func main() {
	var args ArgsModel

	arg.MustParse(&args)

	if args.WorkDir == "" {
		args.WorkDir = "."
	}

	if args.Size {
		runCountSize(&args)
	}

	if args.Lines {
		runCountLines(&args)
	}
}
