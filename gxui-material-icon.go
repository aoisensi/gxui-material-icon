package materialicon

import (
	"bytes"
	"compress/flate"
	"io/ioutil"

	"github.com/google/gxui"
)

//go:generate go run generator_font.go
//go:generate go run generator_map.go

var (
	// MaterialIcon is material-icon font.
	MaterialIcon = inflate(data)
	// MaterialIconFonts is font data, key is size.
	MaterialIconFonts = map[int]gxui.Font{}
)

// CreateIcon is create icon control
func CreateIcon(theme gxui.Theme, icon rune, size int) gxui.Label {
	if theme == nil || size < 0 {
		return nil
	}
	control := theme.CreateLabel()
	font, ok := MaterialIconFonts[size]
	if !ok {
		var err error
		font, err = theme.Driver().CreateFont(MaterialIcon, size)
		if err != nil {
			return nil
		}
		MaterialIconFonts[size] = font
	}
	control.SetFont(font)
	control.SetText(string(icon))
	return control
}

func inflate(src []byte) []byte {
	r := bytes.NewReader(src)
	b, err := ioutil.ReadAll(flate.NewReader(r))
	if err != nil {
		panic(err)
	}
	return b
}
