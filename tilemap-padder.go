package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func main() {

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "show help",
	}

	app := &cli.App{
		Name:      "tilemap-padder",
		Usage:     "adds padding to tilemaps",
		UsageText: "tilemap-padder -w $WIDTH -h $HEIGHT -p $PADDING -i $FILENAME.png -o $FILENAME_padding.png",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:     "tile-width",
				Usage:    "width of tile in pixels",
				Aliases:  []string{"w"},
				Value:    32,
				Required: true,
			},
			&cli.UintFlag{
				Name:     "tile-height",
				Usage:    "height of tile in pixels",
				Aliases:  []string{"h"},
				Value:    32,
				Required: true,
			},
			&cli.UintFlag{
				Name:     "padding",
				Usage:    "padding in pixels",
				Aliases:  []string{"p"},
				Value:    0,
				Required: true,
			},
			&cli.StringFlag{
				Name:     "input",
				Usage:    "input image (PNG)",
				Aliases:  []string{"i"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "output image (PNG)",
				Aliases:  []string{"o"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			file, err := os.Open(c.String("input"))
			if err != nil {
				return err
			}
			tilesetPNG, err := png.Decode(file)
			if err != nil {
				return err
			}

			tileWidth := int(c.Uint("tile-width"))
			tileHeight := int(c.Uint("tile-height"))
			padding := int(c.Uint("padding"))

			n := tilesetPNG.Bounds().Dx() / tileWidth
			m := tilesetPNG.Bounds().Dy() / tileHeight

			newBounds := image.Rect(0, 0, n*(tileWidth+padding*2), m*(tileHeight+padding*2))
			rgba := image.NewRGBA(newBounds.Bounds())

			xOffset := 0
			yOffset := 0

			for y := 0; y < m; y++ {
				yOffset += padding
				for x := 0; x < n; x++ {
					xOffset += padding

					for i := xOffset - padding; i < xOffset; i++ {
						// Left
						draw.Draw(rgba,
							image.Rect(
								i+x*tileWidth,
								yOffset+y*tileHeight,
								i+(x+1)*tileWidth,
								yOffset+(y+1)*tileHeight),
							tilesetPNG,
							image.Point{x * tileWidth, y * tileHeight},
							draw.Src)
					}
					for i := yOffset - padding; i < yOffset; i++ {
						// Up
						draw.Draw(rgba,
							image.Rect(
								xOffset+x*tileWidth,
								i+y*tileHeight,
								xOffset+(x+1)*tileWidth,
								i+(y+1)*tileHeight),
							tilesetPNG,
							image.Point{x * tileWidth, y * tileHeight},
							draw.Src)
					}

					for i := xOffset + padding + 1; i > xOffset; i-- {
						// Right
						draw.Draw(rgba,
							image.Rect(
								i+x*tileWidth,
								yOffset+y*tileHeight,
								i+(x+1)*tileWidth,
								yOffset+(y+1)*tileHeight),
							tilesetPNG,
							image.Point{x * tileWidth, y * tileHeight},
							draw.Src)
					}

					for i := yOffset + padding + 1; i > yOffset; i-- {
						// Down
						draw.Draw(rgba,
							image.Rect(
								xOffset+x*tileWidth,
								i+y*tileHeight,
								xOffset+(x+1)*tileWidth,
								i+(y+1)*tileHeight),
							tilesetPNG,
							image.Point{x * tileWidth, y * tileHeight},
							draw.Src)
					}

					// Center
					draw.Draw(rgba,
						image.Rect(
							xOffset+x*tileWidth,
							yOffset+y*tileHeight,
							xOffset+(x+1)*tileWidth,
							yOffset+(y+1)*tileHeight),
						tilesetPNG,
						image.Point{x * tileWidth, y * tileHeight},
						draw.Src)
					xOffset += padding
				}
				yOffset += padding
				xOffset = 0
			}

			f, err := os.Create(c.String("output"))
			if err != nil {
				return err
			}

			err = png.Encode(f, rgba)
			if err != nil {
				return err
			}

			return nil
		},
	}

	red := color.New(color.FgRed).PrintfFunc()

	err := app.Run(os.Args)
	if err != nil {
		red("Error: %s\n", err)
		os.Exit(1)
	}
}
