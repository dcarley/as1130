package cmd

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"strings"

	"github.com/dcarley/as1130"
	"github.com/spf13/cobra"
)

var textCmd = &cobra.Command{
	Use:   "text string ...",
	Short: "Scroll text across the display",
	Long:  "Scroll text across the display",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if Size != Size24x5 {
			log.Fatal("only 24x5 size is supported for this command")
		}
		if err := scrollText(strings.Join(args, " ")); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(textCmd)
}

func scrollText(text string) error {
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

	frames, err := as1130.TextFrames(text)
	if err != nil {
		return err
	}

	max, err := as.MaxFrames()
	if err != nil {
		return err
	}
	if count := len(frames); count > int(max) {
		return fmt.Errorf("message requires more than %d frames: %d", max, count)
	}

	pwm := NewFrame()
	draw.Draw(pwm, pwm.Bounds(), &image.Uniform{as1130.On}, image.ZP, draw.Src)
	if err := as.SetBlinkAndPWMSet(1, NewFrame(), pwm); err != nil {
		return err
	}

	for i, frame := range frames {
		err := as.SetFrame(uint8(i+1), frame)
		if err != nil {
			return err
		}
	}

	if err := as.SetMovie(as1130.Movie{Display: true}); err != nil {
		return err
	}
	if err := as.SetMovieMode(as1130.MovieMode{Frames: uint8(len(frames))}); err != nil {
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
