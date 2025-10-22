package main

import "testing"

func TestReadAtlasRow(t *testing.T) {

	got, err := ReadAtlas("atiles_test.atlas")
	want := "page \"atiles.bmp\" 255,0,255,255 true 48,48 0,0,0,0"

	if got != want {
		t.Errorf("got %q want %q with error %q", got, want, err)
	}
}
