package command

import (
	"fmt"
	"io/ioutil"

	"github.com/amanhigh/go-fun/kohan/commander/tools"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	composePath   = ""
	composeOpt    = ""
	dockerService = ""
)

const DOCKER_CONFIG = "/tmp/docker-config.yml"

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker Based Commands",
	Args:  cobra.ExactArgs(1),
}

var dockerPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Process Monitor",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(fmt.Sprintf("watch -n1 '%v'", getDockerCmd("ps")))
	},
}

var dockerKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Force kill and Clear Volumes",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(getDockerCmd("rm -svf"))
	},
}

var dockerResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Stop & Clean Containers, Start Fresh",
	Run: func(cmd *cobra.Command, args []string) {
		//Clean old Containers
		tools.PrintCommand("docker-clean stop")

		tools.LiveCommand(getDockerCmd("up -d"))
	},
}

var dockerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Docker Compose",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(getDockerCmd("stop"))
	},
}

var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Docker Compose",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(getDockerCmd("start"))
	},
}

var dockerRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Services",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(getDockerCmd("restart"))
	},
}

var dockerSetCmd = &cobra.Command{
	Use:   "set [files]",
	Short: "Set Docker Config",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		dockerPath := ""
		for _, file := range args {
			dockerPath += fmt.Sprintf("-f %v/%v.yml ", composePath, file)
		}
		fmt.Println(dockerPath, composeOpt)

		bytes, _ := yaml.Marshal(frameworks.DockerConfig{
			Path: dockerPath,
			Opts: composeOpt,
		})
		util.PrintGreen(fmt.Sprintf("Written Config: %v\n\n%v", DOCKER_CONFIG, string(bytes)))
		err = ioutil.WriteFile(DOCKER_CONFIG, bytes, util.DEFAULT_PERM)
		return
	},
}

func init() {
	RootCmd.AddCommand(dockerCmd)
	dockerCmd.PersistentFlags().StringVarP(&composePath, "path", "p", "/Users/amanpreet.singh/IdeaProjects/GoArena/src/github.com/amanhigh/go-fun/Docker/compose", "Compose Path for Docker")
	dockerCmd.PersistentFlags().StringVarP(&dockerService, "svc", "s", "", "Specify Service to Act On")
	dockerSetCmd.Flags().StringVarP(&composeOpt, "opt", "o", "", "Compose Options like Scale")

	dockerCmd.AddCommand(dockerSetCmd)
	dockerCmd.AddCommand(dockerPsCmd)

	dockerCmd.AddCommand(dockerStartCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	dockerCmd.AddCommand(dockerRestartCmd)

	dockerCmd.AddCommand(dockerKillCmd)
	dockerCmd.AddCommand(dockerResetCmd)
}

func getDockerCmd(action string) (cmd string) {
	var dockerConfig frameworks.DockerConfig
	bytes, _ := ioutil.ReadFile(DOCKER_CONFIG)
	_ = yaml.Unmarshal(bytes, &dockerConfig)
	cmd = fmt.Sprintf("docker-compose %v %v %v %v", dockerConfig.Path, action, dockerService, dockerConfig.Opts)
	fmt.Println(cmd)
	return
}
