package cmd

import (
	"fmt"
	"os"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var (
	Size    SizeVar
	Device  string
	Address int
)

const (
	Size24x5 = iota
	Size12x11
)

type SizeVar int

func (s *SizeVar) String() string {
	return fmt.Sprintf("%d", *s)
}

func (s *SizeVar) Set(val string) error {
	switch val {
	case "24x5":
		*s = Size24x5
	case "12x11":
		*s = Size12x11
	default:
		return fmt.Errorf("invalid size")
	}

	return nil
}

func (s *SizeVar) Type() string {
	return "string"
}

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
	RootCmd.PersistentFlags().VarP(&Size, "size", "s", "size of the frame: 24x5, 12x11 (default 24x5)")
	RootCmd.PersistentFlags().StringVarP(&Device, "device", "d", as1130.DeviceDefault, "I2C device path")
	RootCmd.PersistentFlags().IntVarP(&Address, "address", "a", 0, fmt.Sprintf(
		"I2C device address (default 0x%02X)", as1130.AddressDefault,
	))
}

func NewFrame() as1130.Framer {
	switch Size {
	case Size24x5:
		return as1130.NewFrame24x5()
	case Size12x11:
		return as1130.NewFrame12x11()
	}

	panic("invalid frame size")
}
