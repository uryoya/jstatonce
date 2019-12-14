package main

import (
	"strconv"
	"testing"
)

func Test_replaceVmid(t *testing.T) {
	args := []string{"-gc", "-h5", "vmid", "1000"}
	vmid := 12345
	expect := strconv.Itoa(vmid)

	result := replaceVmid(args, vmid)
	if result[2] != expect {
		t.Fatalf("expected `%s` but got `%s`", expect, result[2])
	}
}
