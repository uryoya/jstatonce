package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/* jstatonce "[java execute command]" [jstat arguments ( without vmid )]
 *
 * e.g.
 * $ jstatonce
 */
func main() {
	javaFields := strings.Fields(os.Args[1])
	javaCmd := javaFields[0]
	javaArgs := javaFields[1:]

	jvmCmd := exec.Command(javaCmd, javaArgs...)
	err := jvmCmd.Start()
	if err != nil {
		fmt.Errorf("failed execute java: %s\n", err)
	}

	vmid := strconv.Itoa(jvmCmd.Process.Pid)
	jstatCmd := exec.Command("jstat", "-gcutil", vmid, "1000")
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
	// go func() {
	// 	stdoutScn := bufio.NewScanner(jstatStdout)
	// 	stderrScn := bufio.NewScanner(jstatStderr)
	// 	for stdoutScn.Scan() {
	// 		fmt.Printf("%s", stdoutScn.Text())
	// 	}
	// 	for stderrScn.Scan() {
	// 		fmt.Printf("%s", stderrScn.Text())
	// 	}
	// }()

	err = jstatCmd.Run()
	if err != nil {
		fmt.Errorf("failed execute jstat: %s\n", err)
	}

	jvmCmd.Wait()
}
