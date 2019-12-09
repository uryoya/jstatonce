package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	javaCmd := os.Args[1]
	javaArgs := strings.Fields(os.Args[2])

	jvmCmd := exec.Command(javaCmd, javaArgs...)
	err := jvmCmd.Start()
	if err != nil {
		fmt.Errorf("failed execute java: %s\n", err)
	}

	fmt.Println(jvmCmd.Process.Pid)
}
