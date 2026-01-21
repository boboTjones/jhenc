# jhEnc

This is a small Go library to encode `image.Image` to JSON and HTML tables for displaying images (like TOTP QR codes) inline, inspired by how GitHub delivers TOTP QR codes.

Why? Because you never want the image just sitting around and you definitely don't want to be sending keys to some third party API like some authentication providers I won't mention.

## Features

- **JSON encoding**: Encodes images as compact RGB arrays `[[r,g,b], ...]`
- **HTML encoding**: Two modes available
  - Inline styles: Each pixel gets `style='background-color:rgb(...)'`
  - CSS classes: More compact HTML using generated class names
- **Alpha channel support**: Transparent pixels are blended with white background
- **Modern HTML/CSS**: No deprecated attributes, uses `border-spacing` instead of `cellspacing/cellpadding`

## Usage

```go
import "github.com/yourusername/jhenc"

// JSON encoding
err := jhenc.JSONEncode(writer, img)

// HTML encoding with inline styles
err := jhenc.HTMLEncode(writer, img)

// HTML encoding with CSS classes (more compact)
err := jhenc.HTMLEncodeWithClasses(writer, img)
css, err := jhenc.GenerateCSS(img)
// Include the CSS in your <style> tag
```

## API

- `JSONEncode(w io.Writer, img image.Image)` - Encode as JSON with RGB arrays
- `HTMLEncode(w io.Writer, img image.Image)` - Encode as HTML table with inline styles
- `HTMLEncodeWithClasses(w io.Writer, img image.Image)` - Encode as HTML table with CSS classes
- `GenerateCSS(img image.Image)` - Generate CSS for use with `HTMLEncodeWithClasses` 
