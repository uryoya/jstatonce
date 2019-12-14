package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

/* jstatonce "[java execute command]" [jstat arguments ( without vmid )]
 *
 * e.g.
 * $ jstatonce
 */
func main() {
	javaFields := strings.Fields(os.Args[1])
	javaCmd := javaFields[0]
	javaArgs := javaFields[1:]
	jstatArgs := os.Args[2:]

	jvmCmd := exec.Command(javaCmd, javaArgs...)
	err := jvmCmd.Start()
	if err != nil {
		fmt.Errorf("failed execute java: %s\n", err)
	}

	replacedArgs := replaceVmid(jstatArgs, jvmCmd.Process.Pid)
	jstatCmd := exec.Command("jstat", replacedArgs...)
	jstatStdout, err := jstatCmd.StdoutPipe()
	if err != nil {
		fmt.Errorf("failed pipe jstat stdout: %s\n", err)
	}
	defer jstatStdout.Close()
	jstatStderr, err := jstatCmd.StderrPipe()
	if err != nil {
		fmt.Errorf("failed pipe jstat stderr: %s\n", err)
	}
	defer jstatStderr.Close()

	go func() {
		io.Copy(os.Stdout, jstatStdout)
		io.Copy(os.Stderr, jstatStderr)
	}()

	err = jstatCmd.Run()
	if err != nil {
		fmt.Errorf("failed execute jstat: %s\n", err)
	}

	jvmCmd.Wait()
}
