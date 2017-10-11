package jhenc

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

type JSONResponse struct {
	Rows [][]string
}

func JsonEncode(w io.Writer, img image.Image) error {
	var err error
	// Check image size
	imgwidth, imgheight := int64(img.Bounds().Dx()), int64(img.Bounds().Dy())
	if imgwidth <= 0 || imgheight <= 0 || imgwidth >= 1<<32 || imgheight >= 1<<32 {
		return errors.New(fmt.Sprintf("Unacceptable image size %d x %d.\n", imgwidth, imgheight))
	}

	min := img.Bounds().Min
	max := img.Bounds().Max
	rows := [][]string{}

	for y := min.Y; y <= max.Y; y++ {
		row := []string{}
		for x := min.X; x <= max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			row = append(row, fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B))

		}
		rows = append(rows, row)
	}

	response, err := json.Marshal(JSONResponse{Rows: rows})
	if err != nil {
		return err
	}
	w.Write(response)
	return err
}

func HtmlEncode(w io.Writer, img image.Image) error {
	var err error
	// Check image size
	imgwidth, imgheight := int64(img.Bounds().Dx()), int64(img.Bounds().Dy())
	if imgwidth <= 0 || imgheight <= 0 || imgwidth >= 1<<32 || imgheight >= 1<<32 {
		return errors.New(fmt.Sprintf("Unacceptable image size %d x %d.\n", imgwidth, imgheight))
	}

	min := img.Bounds().Min
	max := img.Bounds().Max

	out := "<table cellspacing='0' cellpadding='0'>\n<thead></thead>\n<tbody>\n"
	lines := []string{}

	for y := min.Y; y <= max.Y; y++ {
		line := "<tr>"
		for x := min.X; x <= max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			line += fmt.Sprintf("<td style='background-color:rgb(%d, %d, %d);'></td>\n", c.R, c.G, c.B)
		}
		line += "</tr>\n"
		lines = append(lines, line)
	}
	out += strings.Join(lines, "\n")
	w.Write([]byte(out + "</tbody>\n</table>\n"))
	return err
}
