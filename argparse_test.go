package main

import (
	"testing"
)

func Test_argparse_help(t *testing.T) {
	tests := []struct {
		args   []string
		expect bool
	}{
		{[]string{"jstatonce", "-h"}, true},
		{[]string{"jstatonce", "--help"}, true},
		{[]string{"jstatonce", "-o", "app.log"}, false},
		{[]string{"jstatonce", "-o", "app.log", "--help"}, true},
	}

	for _, test := range tests {
		err, opts := argparse(test.args)
		if err != nil {
			t.Errorf("予期しないエラーの発生!: %v", err)
		}

		if opts.needHelp != test.expect {
			t.Fatalf("期待値: `%v` 取得値: `%v`", test.expect, opts.needHelp)
		}
	}
}

func Test_argparse_javaOutput(t *testing.T) {
	tests := []struct {
		args   []string
		expect string
		err    bool
	}{
		{[]string{"jstatonce", "-h"}, "/dev/null", false},
		{[]string{"jstatonce", "-o", "app.log"}, "app.log", false},
		{[]string{"jstatonce", "--java-output", "app.log"}, "app.log", false},
		{[]string{"jstatonce", "-o"}, "/dev/null", true},
	}

	for _, test := range tests {
		err, opts := argparse(test.args)
		if !test.err && err != nil {
			t.Errorf("予期しないエラーの発生!: %v", err)
		}

		if test.err && err == nil {
			t.Fatal("エラーにしてください")
		}

		if opts.javaOutput != test.expect {
			t.Fatalf("期待値: `%s` 取得値: `%s`", test.expect, opts.javaOutput)
		}
	}
}
