package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

const (
	usage = `push [image]:[tag] to [registry]/[project]
	      push [image]:[tag] to [registry]/[project]/[image]:[tag]`
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: usage,
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 || args[1] != "to" {
			fmt.Printf(usage)
		}

		oldImage := getImage(args[0])
		newImage, err := completeNewImage(args[2], oldImage)
		if err != nil {
			return err
		}
		fmt.Printf("docker tag %s %s\n", args[0], newImage)
		cmdline := exec.Command("docker", "tag", args[0], newImage)
		var out bytes.Buffer
		cmdline.Stdout = &out
		err = cmdline.Run()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out.String())
		fmt.Printf("%s\n", "Done.")

		fmt.Printf("docker push %s\n", newImage)
		cmdline = exec.Command("docker", "push", newImage)
		cmdline.Stdout = &out
		err = cmdline.Run()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out.String())

		return nil
	},
}

func getImage(image string) string {
	img := strings.Split(image, "/")
	return img[len(img)-1]
}

func completeNewImage(registryProject, oldImage string) (string, error) {

	p := strings.Split(registryProject, "/")

	if len(p) < 2 || len(p) > 3 {

		return "", errors.New("Usage errors")
	}
	if len(p) == 3 {
		return registryProject, nil
	}
	return registryProject + "/" + oldImage, nil
}
