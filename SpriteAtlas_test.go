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

// func TestGetAnimation(t *testing.T) {

// 	var reg Region
// 	err := reg.ParseRegionStr([]string{"player_walk", "1,148,384,244", "47,47", "north,1,1,4", "west,1,5,4", "south,2,1,4", "east,2,5,4"})
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	got1, err := reg.GetAnimation("north", 0)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	got2, err := reg.GetAnimation("north", 1)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	got3, err := reg.GetAnimation("north", 2)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	got4, err := reg.GetAnimation("north", 3)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	got5, err := reg.GetAnimation("north", 4)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	want1 := RECT{X: 1, Y: 148, Width: 47, Height: 47}
// 	want2 := RECT{X: 48, Y: 148, Width: 47, Height: 47}
// 	want3 := RECT{X: 95, Y: 148, Width: 47, Height: 47}
// 	want4 := RECT{X: 142, Y: 148, Width: 47, Height: 47}
// 	want5 := RECT{X: 1, Y: 148, Width: 47, Height: 47}
// 	if got1 != want1 {
// 		t.Errorf("got %v want %v", got1, want1)
// 	}
// 	if got2 != want2 {
// 		t.Errorf("got %v want %v", got2, want2)
// 	}
// 	if got3 != want3 {
// 		t.Errorf("got %v want %v", got3, want3)
// 	}
// 	if got4 != want4 {
// 		t.Errorf("got %v want %v", got4, want4)
// 	}
// 	if got5 != want5 {
// 		t.Errorf("got %v want %v", got5, want5)
// 	}
// }

func ExampleStripAtlasLine() {
	// the last comment actualy causes this to be run ~ not just compiler tested
	b := []byte("abc\n\r")
	got := StripAtlasLine(b)

	fmt.Println(got)
	// Output: abc
}
