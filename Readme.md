# SpriteAtlas for Go - currently to read the sprite atlas and provide sufficient data to read frames

Changed the spriteatlas path for public use
still not ready and is subject to changes, but can be used as example till ready

example use at the moment

## in go.mod

module YourApp

go 1.25.2

require github.com/CrazyHat62/SpriteAtlas v0.1.1


## in main.go

package main

import (
	"os"

	sa "github.com/CrazyHat62/spriteatlas"
)

func main() {

	page, region, err := sa.Spriteatlas("", "atiles.atlas")
	if err != nil {
		os.Exit(1)
	}
	println(page.PageToStr())
	println(region.RegionToStr())
	println(err)

}

## Use main.go and Run go mod tidy

as it says the version may not work as will change as time goes
