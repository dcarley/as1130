package cmd

import (
	"fmt"
	"os"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var (
	Device  string
	Address int
)

var RootCmd = &cobra.Command{
	Use:   "as1130ctl",
	Short: "CLI for controlling an AS1130 driver",
	Long:  "CLI for controlling an AS1130 driver",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Device, "device", "d", as1130.DeviceDefault, "I2C device path")
	RootCmd.PersistentFlags().IntVarP(&Address, "address", "a", 0, fmt.Sprintf(
		"I2C device address (default 0x%02X)", as1130.AddressDefault,
	))
}
