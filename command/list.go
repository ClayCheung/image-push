package command

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all docker images locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdline := exec.Command("docker", "images")
		var out bytes.Buffer
		cmdline.Stdout = &out
		err := cmdline.Run()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out.String())
		return nil
	},
}
