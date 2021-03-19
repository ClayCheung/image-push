package command

import (
	"errors"
	"fmt"
	"github.com/ClayCheung/image-push/util"
	"github.com/spf13/cobra"
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

		fmt.Printf("docker pull %s \n", args[0])

		if err := util.DoCmd("docker", "pull", args[0]); err != nil {
			fmt.Println("Find NO image in remote, use local image.")
		}

		fmt.Printf("docker tag %s %s\n", args[0], newImage)
		if err := util.DoCmd("docker", "tag", args[0], newImage); err != nil {
			return err
		}

		fmt.Printf("%s\n", "Done.")

		fmt.Printf("docker push %s\n", newImage)
		if err := util.DoCmd("docker", "push", newImage); err != nil {
			return err
		}

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
