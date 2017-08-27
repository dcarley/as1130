package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read register subregister",
	Short: "Read data from a register",
	Long: `Read data from a register

This can be useful for debugging a running device.
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		register, err := strconv.ParseInt(args[0], 0, 16)
		if err != nil {
			log.Fatal(err)
		}
		subregister, err := strconv.ParseInt(args[1], 0, 16)
		if err != nil {
			log.Fatal(err)
		}

		if err := read(byte(register), byte(subregister), os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
}

func read(register, subregister byte, out io.Writer) error {
	as, err := as1130.NewAS1130(Device, Address)
	if err != nil {
		return err
	}
	defer as.Close()

	data, err := as.Read(register, subregister)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "bin: %08b\n", data)
	fmt.Fprintf(out, "hex: 0x%02x\n", data)
	fmt.Fprintf(out, "dec: %d\n", data)

	return nil
}
