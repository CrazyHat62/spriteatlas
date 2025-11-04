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
	want := "page atiles.bmp has alphacolor=255,0,255,255 and sheetSize 1729 874 with margins 1,1,1,1. region player_walk has 4 animations"
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
	want := "page atiles.bmp has alphacolor=255,0,255,255 and sheetSize 48 48 with margins 0,0,0,0"
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

func TestGetAnimation(t *testing.T) {

	var reg Region
	err := reg.ParseRegionStr([]string{"player_walk", "1,148,384,244", "47,47", "north,1,1,4", "west,1,5,4", "south,2,1,4", "east,2,5,4"})
	if err != nil {
		t.Error(err.Error())
	}

	want := RECT{X: 1, Y: 148, Width: 47, Height: 47}
	got, nextidx, err := reg.GetAnimation("north", 0)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: 49, Y: 148, Width: 47, Height: 47}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: 97, Y: 148, Width: 47, Height: 47}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: 145, Y: 148, Width: 47, Height: 47}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: 1, Y: 148, Width: 47, Height: 47}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

}

func ExampleStripAtlasLine() {
	// the last comment actualy causes this to be run ~ not just compiler tested
	b := []byte("abc\n\r")
	got := StripAtlasLine(b)

	fmt.Println(got)
	// Output: abc
}
