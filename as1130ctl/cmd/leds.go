package cmd

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"strconv"
	"strings"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var (
	ledsAllOn bool
)

var ledsCmd = &cobra.Command{
	Use:   "leds x,y ...",
	Short: "Turn on specific LEDs",
	Long: `Turn on specific LEDs

x,y co-ordinates are 1-indexed.
No co-ordinates will result in all LEDs being turned off.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		frame := NewFrame()

		if ledsAllOn {
			draw.Draw(frame, frame.Bounds(), &image.Uniform{as1130.On}, image.ZP, draw.Src)
		}

		points, err := parseArgs(args, frame.Bounds())
		if err != nil {
			log.Fatal(err)
		}

		if err := enableLEDs(points, frame); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	ledsCmd.Flags().BoolVar(&ledsAllOn, "all", false, "Turn all LEDs on")
	RootCmd.AddCommand(ledsCmd)
}

func parseArgs(args []string, bounds image.Rectangle) ([]image.Point, error) {
	var (
		points  []image.Point
		invalid []string
	)

	for _, arg := range args {
		coords := strings.Split(arg, ",")
		if len(coords) != 2 {
			invalid = append(invalid, arg)
			continue
		}
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			invalid = append(invalid, arg)
			continue
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			invalid = append(invalid, arg)
			continue
		}
		if x <= bounds.Min.X || x > bounds.Max.X || y <= bounds.Min.Y || y > bounds.Max.Y {
			invalid = append(invalid, arg)
			continue
		}

		points = append(points, image.Point{x - 1, y - 1})
	}

	if len(invalid) > 0 {
		return []image.Point{}, fmt.Errorf("invalid LED positions: %s", strings.Join(invalid, " "))
	}

	return points, nil
}

func enableLEDs(points []image.Point, frame as1130.Framer) error {
	as, err := as1130.NewAS1130(Device, Address)
	if err != nil {
		return err
	}
	defer as.Close()

	if err := as.Reset(); err != nil {
		return err
	}
	if err := as.Init(1); err != nil {
		return err
	}

	pwm := NewFrame()
	draw.Draw(pwm, pwm.Bounds(), &image.Uniform{as1130.On}, image.ZP, draw.Src)
	if err := as.SetBlinkAndPWMSet(1, NewFrame(), pwm); err != nil {
		return err
	}

	for _, point := range points {
		frame.Set(point.X, point.Y, as1130.On)
	}
	if err := as.SetFrame(1, frame); err != nil {
		return err
	}

	if err := as.SetPicture(as1130.Picture{Display: true}); err != nil {
		return err
	}
	if err := as.Start(); err != nil {
		return err
	}

	return nil
}
