package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

func replaceVmid(jstatArgs []string, vmid int) []string {
	replacedArgs := []string{}
	for _, arg := range jstatArgs {
		if arg == "vmid" {
			replacedArgs = append(replacedArgs, strconv.Itoa(vmid))
		} else {
			replacedArgs = append(replacedArgs, arg)
		}
	}

	return replacedArgs
}

func help() string {
	return `使い方: jstatonce [jstatonce options...] "[java args]" "[jstat args]"

jstatonce option
-h --help                             : このヘルプを表示
-o --java-output "path/to/outputfile" : java command の出力先
                                        指定しない場合は /dev/null が指定されます

java args
- java の実行パスを含めてください
e.g. "java Application -J-Xms512m -J-Xmx512m"

jstat arguments
- jstat の実行パスを **含めないで** ください
e.g. "-gc -h5 vmid 1000"
              ^^^^
              jstat コマンドで vmid を挿入する箇所に "vmid" と描いてください。
　　　　　　　実際の java アプリケーションの PID が置換されます。

使用例
$ jstatonce -o app.log "java -cp app.jar com.example.App -J-Xms100m -J-Xmx100m" "-gc -h5 vmid 1000ms"
`
}

func main() {
	err, opts := argparse(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if opts.needHelp {
		fmt.Print(help())
		os.Exit(0)
	}

	javaCmd := opts.javaArgs[0]
	javaArgs := opts.javaArgs[1:]

	jvmOutFile, err := os.Create(opts.javaOutput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open file: %s, couse %s\n", opts.javaOutput, err)
		os.Exit(1)
	}
	defer jvmOutFile.Close()

	jvmCmd := exec.Command(javaCmd, javaArgs...)
	jvmCmdStdout, err := jvmCmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed pipe jvm stdout: %s\n", err)
		os.Exit(1)
	}

	go func() {
		io.Copy(jvmOutFile, jvmCmdStdout)
	}()

	err = jvmCmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed execute java: %s\n", err)
		os.Exit(1)
	}

	replacedArgs := replaceVmid(opts.jstatArgs, jvmCmd.Process.Pid)
	jstatCmd := exec.Command("jstat", replacedArgs...)
	jstatStdout, err := jstatCmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed pipe jstat stdout: %s\n", err)
		os.Exit(1)
	}
	defer jstatStdout.Close()
	jstatStderr, err := jstatCmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed pipe jstat stderr: %s\n", err)
		os.Exit(1)
	}
	defer jstatStderr.Close()

	go func() {
		io.Copy(os.Stdout, jstatStdout)
		io.Copy(os.Stderr, jstatStderr)
	}()

	err = jstatCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed execute jstat: %s\n", err)
		os.Exit(1)
	}

	jvmCmd.Wait()
}
