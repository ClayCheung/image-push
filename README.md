# himg
## Feature
Do this below by one command  :

pull image from Registry-A and rename it, then push it to another Registry-B.

## Usage

```shell
➜ himg --help                                                                                                                       
NO more docker pull/tag/push

Usage:
  himg [flags]
  himg [command]

Available Commands:
  help        Help about any command
  list        list all docker images locally
  push        push [image]:[tag] to [registry]/[project]
              push [image]:[tag] to [registry]/[project]/[image]:[tag]

Flags:
  -h, --help   help for himg

Use "himg [command] --help" for more information about a command.

```