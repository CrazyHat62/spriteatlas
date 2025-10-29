package spriteatlas

import (
	"fmt"
	"testing"
)

// TestReadAtlasRow is a test to verify file open and a general Parse the of specific atlas
func TestOpenAtlas(t *testing.T) {

	page, region, err := Spriteatlas("", "atiles_test.atlas")
	if err != nil {
		t.Errorf("got open file error %q", err)
	}
	want := "page atiles.bmp has default tilesize 48 48. region player_walk has 4 animations"
	got := page.PageToStr() + ". " + region.RegionToStr()
	if got != want {
		t.Errorf("got %q want %q", got, want)
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

func TestParsePageStr(t *testing.T) {
	var page Page
	err := page.ParsePageStr([]string{"atiles.bmp", "255,0,255,255", "true", "48,48", "0,0,0,0"})
	got := page.PageToStr()
	want := "page atiles.bmp has default tilesize 48 48"
	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}

func TestParseRegionStr(t *testing.T) {
	var reg Region
	err := reg.ParseRegionStr([]string{"player_walk", "1,148,384,244", "48,48", "north,1,1,4", "west,1,5,4", "south,2,1,4", "east,2,5,4"})
	got := reg.RegionToStr()
	want := "region player_walk has 4 animations"
	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}

func TestAnimKeys(t *testing.T) {
	var reg Region
	err := reg.ParseRegionStr([]string{"player_walk", "1,148,384,244", "48,48", "north,1,1,4", "west,1,5,4", "south,2,1,4", "east,2,5,4"})
	for _, key := range reg.AnimKeys() {
		got := key
		switch {
		case key == "north" || key == "south" || key == "east" || key == "west":
			continue
		default:
			want := "north, south, east, or west"
			t.Errorf("got %q want %q with error %q", got, want, err)
		}
	}
}

func ExampleStripAtlasLine() {
	// the last comment actualy causes this to be run ~ not just compiler tested
	b := []byte("abc\n\r")
	got := StripAtlasLine(b)

	fmt.Println(got)
	// Output: abc
}
