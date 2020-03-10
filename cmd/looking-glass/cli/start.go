package cli

import (
	"log"

	"github.com/simplifi/looking-glass/pkg/looking-glass/agent"
	"github.com/simplifi/looking-glass/pkg/looking-glass/config"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Looking Glass agent",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	startCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"/etc/looking-glass.yml",
		"the full path to the yaml config file, default: /etc/looking-glass.yml")
	rootCmd.AddCommand(startCmd)
}

// Starts up the agent
func start() {
	var exit = make(chan bool)

	log.Printf("INFO: Starting looking-glass")

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Panicf("ERROR: Failed to load config: %v", err)
	}

	for _, agtConfig := range cfg.Agents {
		agt, err := agent.New(cfg.Artifactory, agtConfig)
		if err != nil {
			log.Panicf("ERROR: Failed to start agent '%v': %v", agtConfig.Name, err)
		}
		go agt.Start()
	}
	// Block until something kills the process
	<-exit
}
