package command

import (
	"fmt"
	"io/ioutil"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/apps/models/learn/frameworks"

	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	composePath   = ""
	composeOpt    = ""
	dockerService = ""
	shell         = ""
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

var dockerLogsCmd = &cobra.Command{
	Use:   "logs [t]",
	Short: "Show Logs, use t for tailing",
	Run: func(cmd *cobra.Command, args []string) {
		action := "logs"
		if len(args) > 0 {
			action += " -f"
		}
		tools.LiveCommand(getDockerCmd(action))
	},
}

var dockerBuildCmd = &cobra.Command{
	Use:   "build [imageName]",
	Short: "Rebuild Docker Image without Cache",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dockerService = args[0]
		tools.LiveCommand(fmt.Sprintf("cd %v/..;docker build ./%v --no-cache -t %v:latest", composePath, args[0], args[0]))
		tools.LiveCommand(getDockerCmd("restart"))
	},
}

var dockerLoginCmd = &cobra.Command{
	Use:   "login [svcName] [#Container]",
	Short: "Login to Specified Docker Compose Container",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(fmt.Sprintf("docker exec -it compose_%v_%v %v", args[0], args[1], shell))
	},
}

var dockerRunCmd = &cobra.Command{
	Use:   "run [svcName] [#Container] [cmd]",
	Short: "Run a command in Specified Docker Compose Container",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(fmt.Sprintf("docker exec compose_%v_%v %v -c \"%v\"", args[0], args[1], shell, args[2]))
	},
}

var dockerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Docker Compose",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand("docker build")
		tools.LiveCommand(getDockerCmd("stop"))
	},
}

var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Docker Compose",
	Run: func(cmd *cobra.Command, args []string) {
		tools.LiveCommand(getDockerCmd("up -d"))
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
		})
		util.PrintGreen(fmt.Sprintf("Written Config: %v\n\n%v", DOCKER_CONFIG, string(bytes)))
		err = ioutil.WriteFile(DOCKER_CONFIG, bytes, util.DEFAULT_PERM)
		return
	},
}

func init() {
	RootCmd.AddCommand(dockerCmd)
	dockerCmd.PersistentFlags().StringVarP(&composePath, "path", "p", "/Users/amanpreet.singh/IdeaProjects/Go/go-fun/Docker/compose/", "Compose Path for Docker")
	dockerCmd.PersistentFlags().StringVarP(&dockerService, "svc", "s", "", "Specify Service to Act On")
	dockerCmd.PersistentFlags().StringVarP(&composeOpt, "opt", "o", "", "Compose Options.Eg: --scale target=3")
	dockerCmd.PersistentFlags().StringVarP(&shell, "shell", "l", "bash", "Login Shell")

	dockerCmd.AddCommand(dockerSetCmd)
	dockerCmd.AddCommand(dockerPsCmd)
	dockerCmd.AddCommand(dockerLoginCmd)
	dockerCmd.AddCommand(dockerRunCmd)
	dockerCmd.AddCommand(dockerLogsCmd)

	dockerCmd.AddCommand(dockerStartCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	dockerCmd.AddCommand(dockerRestartCmd)

	dockerCmd.AddCommand(dockerKillCmd)
	dockerCmd.AddCommand(dockerResetCmd)
	dockerCmd.AddCommand(dockerBuildCmd)
}

func getDockerCmd(action string) (cmd string) {
	var dockerConfig frameworks.DockerConfig
	bytes, _ := ioutil.ReadFile(DOCKER_CONFIG)
	_ = yaml.Unmarshal(bytes, &dockerConfig)
	cmd = fmt.Sprintf("docker-compose %v %v %v %v", dockerConfig.Path, action, dockerService, composeOpt)
	return
}
