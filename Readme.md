# SpriteAtlas for Go - currently to read the sprite atlas and provide sufficient data to read frames

[MIT licence](MIT-LICENCE.md)

Changed the spriteatlas path for public use
still not ready and is subject to changes, but can be used as example till ready

[example](atiles_test.atlas)



## in main.go EXAMPLE

package main

import (
	"os"

	sa "github.com/CrazyHat62/SpriteAtlas"
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
be sure to use SpriteAtlas with the upper-case characters during this alpha phase ~ it may change

## Using the library

in the go.mod file :

    github.com/CrazyHat62/SpriteAtlas v0.1.2

in your main.go

```
import (
	"fmt"
	"os"

	sa "github.com/CrazyHat62/SpriteAtlas"
	rl "github.com/gen2brain/raylib-go/raylib"
)
```

Note 

	I use raylib-go, and have a function that will replace alpha-color in an 
	found at github.com/CrazyHat62/go-rommy which should be public by the time you find this

