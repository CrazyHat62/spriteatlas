package spriteatlas

import (
	"fmt"
	"testing"
)

var (
	pageAns1    = "page atiles.bmp has alphacolor=255,0,255,255 and sheetSize 1729 874 with margins 1,1,1,1. region player_walk has 4 animations"
	pageAns2    = "page atiles.bmp 255,0,255,255 true 48,48 0,0,0,0"
	pageAns3    = "page atiles.bmp has alphacolor=255,0,255,255 and sheetSize 48 48 with margins 0,0,0,0"
	pageStr1    = "page atiles.bmp 255,0,255,255 true 48,48 0,0,0,0\n\r"
	pageStrArr1 = []string{"atiles.bmp", "255,0,255,255", "true", "48,48", "0,0,0,0"}
	regionStr1  = []string{"player_walk", "1,148,384,243", "48,48", "north,1,1,4", "west,5,1,4", "south,1,2,4", "east,5,2,4"}
)

// TestReadAtlasRow is a test to verify file open and a general Parse the of specific atlas
func TestOpenAtlas(t *testing.T) {

	page, region, err := Spriteatlas("", "atiles_test.atlas")
	if err != nil {
		t.Errorf("got open file error %q", err)
	}
	want := pageAns1
	got := page.PageToStr() + ". " + region.RegionToStr()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestStripAtlasLine(t *testing.T) {
	b := []byte(pageStr1)
	got := StripAtlasLine(b)

	want := pageAns2
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestParsePageStr(t *testing.T) {
	var page Page
	err := page.ParsePageStr(pageStrArr1)
	got := page.PageToStr()
	want := pageAns3
	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}

func TestParseRegionStr(t *testing.T) {
	var reg Region
	err := reg.ParseRegionStr(regionStr1)
	got := reg.RegionToStr()
	want := "region player_walk has 4 animations"
	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}

func TestAnimKeys(t *testing.T) {
	var reg Region
	err := reg.ParseRegionStr(regionStr1)
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
	s := 48
	x := []int{1, 49, 97, 145}
	//x2 := []int{48, 96, 144, 192}
	y := []int{148, 196, 244, 292}
	//y2 := []int{195, 243, 291, 339}
	err := reg.ParseRegionStr(regionStr1)
	if err != nil {
		t.Error(err.Error())
	}

	want := RECT{X: x[0], Y: y[0], Width: s, Height: s}

	got, nextidx, err := reg.GetAnimation("north", 0)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[1], Y: y[0], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[2], Y: y[0], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[3], Y: y[0], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[0], Y: y[0], Width: 48, Height: 48}
	got, nextidx, err = reg.GetAnimation("north", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[0], Y: y[1], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("south", 0)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[1], Y: y[1], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("south", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[2], Y: y[1], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("south", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[3], Y: y[1], Width: s, Height: s}
	got, nextidx, err = reg.GetAnimation("south", nextidx)

	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	want = RECT{X: x[0], Y: y[1], Width: 48, Height: 48}
	got, nextidx, err = reg.GetAnimation("south", nextidx)

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
