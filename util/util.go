package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func DoCmd(name string, args ...string) error {
	cmdline := exec.Command(name , args... )
	var out bytes.Buffer
	cmdline.Stdout = &out
	err := cmdline.Run()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", out.String())
	return nil
}