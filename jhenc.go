// Package jhenc provides utilities for encoding images as JSON and HTML tables.
// This is useful for displaying images (like QR codes) without sending them to
// third-party APIs or storing them as files.
package jhenc

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

// JSONResponse represents an image encoded as a grid of RGB colors.
// Each row contains RGB values as [r, g, b] arrays.
type JSONResponse struct {
	Rows [][][]uint8 `json:"rows"`
}

// JSONEncode writes an image to w as JSON, encoding each pixel as an RGB array.
// Alpha channels are blended with a white background.
// Returns an error if the image size is invalid.
func JSONEncode(w io.Writer, img image.Image) error {
	if err := validateImageSize(img); err != nil {
		return err
	}

	min := img.Bounds().Min
	max := img.Bounds().Max
	rows := make([][][]uint8, 0, max.Y-min.Y)

	for y := min.Y; y < max.Y; y++ {
		row := make([][]uint8, 0, max.X-min.X)
		for x := min.X; x < max.X; x++ {
			r, g, b := blendWithWhite(img.At(x, y))
			row = append(row, []uint8{r, g, b})
		}
		rows = append(rows, row)
	}

	response, err := json.Marshal(JSONResponse{Rows: rows})
	if err != nil {
		return err
	}
	_, err = w.Write(response)
	return err
}

// HTMLEncode writes an image to w as an HTML table with inline styles.
// Each pixel becomes a table cell with a background color.
// Alpha channels are blended with a white background.
// Returns an error if the image size is invalid.
func HTMLEncode(w io.Writer, img image.Image) error {
	if err := validateImageSize(img); err != nil {
		return err
	}

	min := img.Bounds().Min
	max := img.Bounds().Max

	var buf strings.Builder
	buf.WriteString("<table>\n<tbody>\n")

	for y := min.Y; y < max.Y; y++ {
		buf.WriteString("<tr>")
		for x := min.X; x < max.X; x++ {
			r, g, b := blendWithWhite(img.At(x, y))
			buf.WriteString(fmt.Sprintf("<td style='background-color:rgb(%d,%d,%d)'></td>", r, g, b))
		}
		buf.WriteString("</tr>\n")
	}
	buf.WriteString("</tbody>\n</table>\n")

	_, err := w.Write([]byte(buf.String()))
	return err
}

// HTMLEncodeWithClasses writes an image to w as an HTML table using CSS classes.
// This generates more compact HTML by using class names instead of inline styles.
// You must include the CSS from GenerateCSS in your page.
// Alpha channels are blended with a white background.
// Returns an error if the image size is invalid.
func HTMLEncodeWithClasses(w io.Writer, img image.Image) error {
	if err := validateImageSize(img); err != nil {
		return err
	}

	min := img.Bounds().Min
	max := img.Bounds().Max

	// Build color class mapping
	colorMap := make(map[[3]uint8]string)
	colorIndex := 0

	var buf strings.Builder
	buf.WriteString("<table>\n<tbody>\n")

	for y := min.Y; y < max.Y; y++ {
		buf.WriteString("<tr>")
		for x := min.X; x < max.X; x++ {
			rgb := [3]uint8(blendWithWhiteArray(img.At(x, y)))
			className, exists := colorMap[rgb]
			if !exists {
				className = fmt.Sprintf("c%d", colorIndex)
				colorMap[rgb] = className
				colorIndex++
			}
			buf.WriteString(fmt.Sprintf("<td class='%s'></td>", className))
		}
		buf.WriteString("</tr>\n")
	}
	buf.WriteString("</tbody>\n</table>\n")

	_, err := w.Write([]byte(buf.String()))
	return err
}

// GenerateCSS generates the CSS needed for HTMLEncodeWithClasses.
// Pass the same image you'll encode to generate matching color classes.
func GenerateCSS(img image.Image) (string, error) {
	if err := validateImageSize(img); err != nil {
		return "", err
	}

	min := img.Bounds().Min
	max := img.Bounds().Max

	// Build color class mapping
	colorMap := make(map[[3]uint8]string)
	colorIndex := 0

	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			rgb := [3]uint8(blendWithWhiteArray(img.At(x, y)))
			if _, exists := colorMap[rgb]; !exists {
				colorMap[rgb] = fmt.Sprintf("c%d", colorIndex)
				colorIndex++
			}
		}
	}

	var buf strings.Builder
	for rgb, className := range colorMap {
		buf.WriteString(fmt.Sprintf(".%s{background-color:rgb(%d,%d,%d)}\n", className, rgb[0], rgb[1], rgb[2]))
	}

	return buf.String(), nil
}

// validateImageSize checks if an image has acceptable dimensions.
func validateImageSize(img image.Image) error {
	imgwidth, imgheight := int64(img.Bounds().Dx()), int64(img.Bounds().Dy())
	if imgwidth <= 0 || imgheight <= 0 || imgwidth >= 1<<32 || imgheight >= 1<<32 {
		return fmt.Errorf("unacceptable image size %d x %d", imgwidth, imgheight)
	}
	return nil
}

// blendWithWhite converts a color to RGB, blending with white if it has transparency.
func blendWithWhite(c color.Color) (r, g, b uint8) {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	if rgba.A == 255 {
		return rgba.R, rgba.G, rgba.B
	}
	// Blend with white background
	alpha := float64(rgba.A) / 255.0
	r = uint8(float64(rgba.R)*alpha + 255.0*(1.0-alpha))
	g = uint8(float64(rgba.G)*alpha + 255.0*(1.0-alpha))
	b = uint8(float64(rgba.B)*alpha + 255.0*(1.0-alpha))
	return
}

// blendWithWhiteArray is like blendWithWhite but returns an array.
func blendWithWhiteArray(c color.Color) [3]uint8 {
	r, g, b := blendWithWhite(c)
	return [3]uint8{r, g, b}
}
