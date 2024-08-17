/*
Copyright © 2024 Narek Tamazyan <natamazy@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"github.com/xyproto/palgen"
	"github.com/xyproto/png2svg"
)

type Config struct {
	inputFilename         string
	outputFilename        string
	colorPink             bool
	limit                 bool
	singlePixelRectangles bool
	verbose               bool
	palReduction          int
}

var rootCmd = &cobra.Command{
	Use:   "qrdonkey",
	Short: `Tool to generate QR codes.`,
	Long: `
 _____  ___           _                _                   
(  _  )|  _ \        ( )              ( )                  
| ( ) || (_) )      _| |   _     ___  | |/')    __   _   _ 
| | | || ,  /     / _  | / _ \ /  _  \| , <   /__\( ) ( ) |
| (('\|| |\ \    ( (_| |( (_) )| ( ) || |\ \ (  ___/| (_) |
(___\_)(_) (_)    \__,_) \___/ (_) (_)(_) (_) \____) \__, |
                                                    ( )_| |
                                                    \___/'

qrdonkey is a CLI tool for generating QR codes.

This application is a tool to generate
the QR codes in PNG and SVG formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		svg, _ := cmd.Flags().GetBool("s")

		if len(args) != 1 {
			fmt.Println("Usage qrdonkey [--s] link")
			return
		}

		if generateQRPNG(args[0]) != nil {
			fmt.Println("Something went wrong. Please try again or with other link.")
			return
		}

		if svg {
			if generateQRSVG(&Config{
				inputFilename:  "qrdonkey_" + args[0] + ".png",
				outputFilename: "qrdonkey_" + args[0] + ".svg",
			}) != nil {
				fmt.Println("SVG generation went wrong, Donkey is sad. Try again.")
				return
			}
		}

		PrintDonkey()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool("s", false, "For generating SVG too.")
	rootCmd.SetUsageTemplate(`Usage:
	qrdonkey [--s] link
	
Flags:
	--s: For generating QR's SVG too.

`)
}

func generateQRPNG(link string) error {
	err := qrcode.WriteFile(link, qrcode.Medium, 1024, "qrdonkey_"+link+".png")
	if err != nil {
		return err
	}

	return nil
}

func generateQRSVG(c *Config) error {
	var (
		box          *png2svg.Box
		x, y         int
		expanded     bool
		lastx, lasty int
		lastLine     int
		done         bool
	)

	img, err := png2svg.ReadPNG(c.inputFilename, c.verbose)
	if err != nil {
		return err
	}

	if c.palReduction > 0 {
		img, err = palgen.Reduce(img, c.palReduction)
		if err != nil {
			return fmt.Errorf("could not reduce the palette of the given image to a maximum of %d colors", c.palReduction)
		}
	}

	height := img.Bounds().Max.Y - img.Bounds().Min.Y

	pi := png2svg.NewPixelImage(img, c.verbose)
	pi.SetColorOptimize(c.limit)

	percentage := 0
	lastPercentage := 0

	if !c.singlePixelRectangles {
		if c.verbose {
			fmt.Print("Placing rectangles... 0%")
		}

		for !c.singlePixelRectangles && !done {

			x, y = pi.FirstUncovered(lastx, lasty)

			if c.verbose && y != lastLine {
				lastPercentage = percentage
				percentage = int((float64(y) / float64(height)) * 100.0)
				png2svg.Erase(len(fmt.Sprintf("%d%%", lastPercentage)))
				fmt.Printf("%d%%", percentage)
				lastLine = y
			}

			box = pi.CreateBox(x, y)
			expanded = pi.Expand(box)

			pi.CoverBox(box, expanded && c.colorPink, c.limit)

			done = pi.Done(x, y)
		}

		if c.verbose {
			png2svg.Erase(len(fmt.Sprintf("%d%%", lastPercentage)))
			fmt.Println("100%")
		}

	}

	if c.singlePixelRectangles {
		if c.verbose {
			percentage = 0
			lastPercentage = 0
			fmt.Print("Placing 1x1 rectangles... 0%")
			pi.CoverAllPixelsCallback(func(currentIndex, totalLength int) {
				lastPercentage = percentage
				percentage = int((float64(currentIndex) / float64(totalLength)) * 100.0)
				png2svg.Erase(len(fmt.Sprintf("%d%%", lastPercentage)))
				fmt.Printf("%d%%", percentage)

			}, 1024)
			png2svg.Erase(len(fmt.Sprintf("%d%%", lastPercentage)))
			fmt.Println("100%")
		} else {
			pi.CoverAllPixels()
		}
	}

	return pi.WriteSVG(c.outputFilename)
}
