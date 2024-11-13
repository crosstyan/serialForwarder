package cmd

import (
	"github.com/crosstyan/serialForwarder/log"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var listCmd = cobra.Command{
	Use:   "list",
	Short: "List available serial ports",
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	var err error
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Sugar().Error(err)
		return
	}
	if len(ports) == 0 {
		log.Sugar().Info("No serial ports found")
		return
	}
	for _, port := range ports {
		log.Sugar().Info(port)
	}
}
