package cmd

import (
	"image/color"
	"log"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var (
	flowWidth int
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Display fading pattern using PWM",
	Long: `Display fading pattern using PWM
	
This uses a different memory configuration so you may need to hard reset
the device before or after using it for something different.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if flowWidth > 5 {
			log.Fatal("width cannot be greater than 5: ", flowWidth)
		}

		if err := flow(flowWidth); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	flowCmd.Flags().IntVarP(&flowWidth, "width", "w", 2, "Width of pattern, max 5")
	RootCmd.AddCommand(flowCmd)
}

func flow(frames int) error {
	as, err := as1130.NewAS1130(Device, Address)
	if err != nil {
		return err
	}
	defer as.Close()

	if err := as.Reset(); err != nil {
		return err
	}
	if err := as.Init(5); err != nil {
		return err
	}

	var (
		pwm     uint8
		xTotal  int
		noBlink = NewFrame()
	)

	for i := 1; i <= frames; i++ {
		var (
			index = uint8(i)
			frame = NewFrame()
		)

		midway := (frames * frame.Bounds().Max.X) / 2
		step := uint8(255 / midway)
		for x := 0; x < frame.Bounds().Max.X; x++ {
			xTotal++

			if xTotal <= midway {
				pwm += step
			} else {
				pwm -= step
			}

			for y := 0; y < frame.Bounds().Max.Y; y++ {
				frame.Set(x, y, color.Gray{pwm})
			}
		}

		if err := as.SetBlinkAndPWMSet(index, noBlink, frame); err != nil {
			return err
		}

		frame.SetPWMSet(index)
		if err := as.SetFrame(index, frame); err != nil {
			return err
		}
	}

	spacerFrameIndex := uint8(frames + 1)
	if err := as.SetFrame(spacerFrameIndex, NewFrame()); err != nil {
		return err
	}

	if err := as.SetMovie(as1130.Movie{Display: true}); err != nil {
		return err
	}
	if err := as.SetMovieMode(as1130.MovieMode{Frames: spacerFrameIndex}); err != nil {
		return err
	}
	if err := as.SetFrameTime(as1130.FrameTime{Scrolling: true, Delay: 1}); err != nil {
		return err
	}

	if err := as.Start(); err != nil {
		return err
	}

	return nil
}
