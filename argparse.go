package main

import (
	"errors"
	"fmt"
	"strings"
)

type jstatonceOpts struct {
	needHelp   bool
	javaOutput string
	javaArgs   []string
	jstatArgs  []string
}

func argparse(args []string) (opts jstatonceOpts, err error) {
	opts.needHelp = false
	opts.javaOutput = "/dev/null"

	arg := ""
	for i := 1; i < len(args); i++ {
		arg = args[i]
		switch arg {
		case "-h", "--help":
			opts.needHelp = true
			return opts, nil

		case "-o", "--java-output":
			i++
			if len(args) <= i {
				err = errors.New(fmt.Sprintf("%s オプションは引数が必要です", args[i-1]))
				return opts, err
			}
			opts.javaOutput = args[i]

		default:
			if i == len(args)-1 { // 引数の末尾は jstat args
				opts.jstatArgs = strings.Fields(args[i])
			} else if i == len(args)-2 { // 引数の末尾から２番目は java args
				opts.javaArgs = strings.Fields(args[i])
			} else {
				err = errors.New(fmt.Sprintf("未知のオプション: %s, ヘルプを参照: jstatonce -h", args[i]))
				return opts, err
			}
		}
	}

	return opts, nil
}
