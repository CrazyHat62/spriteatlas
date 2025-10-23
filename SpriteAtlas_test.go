package spriteatlas

import (
	"io"
	"strings"
	"testing"
)

// TestReadAtlasRow is a test to verify file open and a general Parse the of specific atlas
func TestOpenAtlas(t *testing.T) {

	file, reader, err1 := OpenAtlas("atiles_test.atlas")
	if err1 != nil && err1 != io.EOF {
		t.Errorf("got open file error %q", err1)
	}
	defer file.Close()

	var str string
	var err error
	var line []byte
	var page Page

	for {
		line, err = reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			str = ""
			break
		}
		str = StripAtlasLine(line)
		if str == "" || str[:1] == "#" || str[:2] == "//" {
			continue
		}
		a := strings.Split(str, " ")
		for i := 0; i < len(a); i++ {
			a[i] = strings.Trim(a[i], " ")
		}
		if a[0] == "page" {
			page.ParsePageStr(a[1:])
			break
		}
		if a[0] != "page" {
			t.Errorf("got %q Not page, ", a[0])
			return
		}
	}

	got := page.PageToStr()

	want := "page atiles.bmp has default tilesize 48 48"

	if got != want {
		if err != nil {
			t.Errorf("got %q want %q with error %q", got, want, err.Error())
		} else {
			t.Errorf("got %q want %q with NO error", got, want)
		}
	}
}

func TestStripAtlasLine(t *testing.T) {
	b := []byte("page atiles.bmp 255,0,255,255 true 48,48 0,0,0,0\n\r")
	got := StripAtlasLine(b)

	want := "page atiles.bmp 255,0,255,255 true 48,48 0,0,0,0"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestPageParse(t *testing.T) {
	var page Page
	err := page.ParsePageStr([]string{"atiles.bmp", "255,0,255,255", "true", "48,48", "0,0,0,0"})
	got := page.PageToStr()
	want := "page atiles.bmp has default tilesize 48 48"
	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}
