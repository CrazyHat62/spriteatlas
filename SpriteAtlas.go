package spriteatlas

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var page Page

// XY is int Vector
type XY struct {
	X int
	Y int
}

// Rect is int Rect
type RECT struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Animation info from sprite-sheet
type Anim struct {
	Pos   XY
	Count int
	Step  bool
	Loop  bool
}

// Region info from sprite-sheet
type Region struct {
	Name     string
	Pos      XY
	TileSize XY
	Anims    map[string]Anim
}

// Page info from sprite-sheet
type Page struct {
	Name             string
	Alpha_color      string
	imageRegionMarks bool
	sheetSize        XY
	margin           RECT
	Regions          map[string]Region
}

// RecToStr - tests/diagnostics/logs/debugging (unlikely to be used in app)
func (r *RECT) RectToStr() string {
	return fmt.Sprintf("rect is at pos %v, %v with size %v, %v", r.X, r.Y, r.Width, r.Height)
}

// AnimKeys - tests/diagnostics/logs/debugging (can be used in app ~ don't rely on order)
func (r *Region) AnimKeys() []string {

	var keys []string
	for key := range r.Anims {
		keys = append(keys, key)
	}
	return keys
}

// RegionToStr - tests/diagnostics/logs/debugging (unlikely to be used in app)
func (r *Region) RegionToStr() string {
	return fmt.Sprintf("region %s has %d animations", r.Name, len(r.Anims))
}

// *********************PAGE*********************

// Spriteatlas reads a atlas at Path + Name ~ use forward slash(s) in path.
// reads the spritesheet as an image and will 'overwrite' alpha color if specified and found.
// alpha-color is NOT the same as pre-multiplied-alpha
func Spriteatlas(filePath string, fileName string) (*Page, error) {

	if len(filePath) > 0 && filePath[len(filePath)-1] != '/' {
		filePath = filePath + "/"
	}

	fileBuf, err := os.ReadFile(filePath + fileName)
	if err != nil {
		return &page, err
	}

	err1 := ParseAtlas(fileBuf)

	return &page, err1
}

func ParseAtlas(fileBytes []byte) error {
	var str string
	var err error
	var region Region

	for line := range bytes.Lines(fileBytes) {

		str = StripAtlasLine(line)
		if str == "" || str[:1] == "#" || str[:2] == "//" {
			continue
		}

		a := strings.Split(str, " ")
		for i := 0; i < len(a); i++ {
			a[i] = strings.Trim(a[i], " ")
		}

		switch {
		case a[0] == "page":
			err = page.ParsePageStr(a[1:])
		case a[0] == "region":
			region = Region{}
			err = region.ParseRegionStr(a[1:])
			page.Regions[region.Name] = region
		default:
			errstr := parseAnimStr(a[:], "", &region)
			if err != nil {
				errstr = err.Error() + errstr
			}
			err = errors.New(errstr)
		}

	}

	//println(page.PageToStr(), region.RegionToStr(), err)

	return err
}

func StripAtlasLine(line []byte) string {
	//remove newline and carrige return if it exists from line
	line = bytes.ReplaceAll(line, []byte("\r"), []byte(""))
	line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
	//remove mutiple spaces
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAll(line, []byte(" "))

	s := string(result)
	s = strings.Trim(s, " ")
	return s
}

func (p *Page) ParsePageStr(values []string) error {
	var err error
	var errstr string = ""
	p.Name = values[0]
	p.Alpha_color = values[1]
	p.imageRegionMarks, err = strconv.ParseBool(values[2])
	if err != nil {
		errstr = errstr + "Parse Atlas region marks failed, "
	}
	s := strings.Split(values[3], ",")
	m := strings.Split(values[4], ",")

	px, err := strconv.ParseInt(s[0], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas tile size X failed, "
	}
	py, err := strconv.ParseInt(s[1], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas tile size Y failed, "
	}

	p.sheetSize.X = int(px)
	p.sheetSize.Y = int(py)
	mx1, err := strconv.ParseInt(m[0], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas margin X1 failed, "
	}
	my1, err := strconv.ParseInt(m[1], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas margin Y1 failed, "
	}
	mx2, err := strconv.ParseInt(m[2], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas margin X2 failed, "
	}
	my2, err := strconv.ParseInt(m[3], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Atlas margin Y2 failed, "
	}
	p.margin.X = int(mx1)
	p.margin.Width = int(mx2)
	p.margin.Y = int(my1)
	p.margin.Height = int(my2)

	if errstr != "" {
		return errors.New(errstr)
	}

	p.Regions = make(map[string]Region)

	return nil
}

func (p *Page) PageToStr() string {
	return fmt.Sprintf("page %s has alphacolor=%s and sheetSize %d %d with margins %d,%d,%d,%d", p.Name, p.Alpha_color, p.sheetSize.X, p.sheetSize.Y, p.margin.X, p.margin.Y, p.margin.Width, p.margin.Height)
}

// *********************REGION*********************

// GetFrameRect gets the RECT for the given animation name and frame index in that animation
func (r *Region) ParseRegionStr(values []string) error {
	var err error
	var errstr string = ""
	r.Anims = make(map[string]Anim)
	//player_walk 1,148,384,244 48,48 north,1,1,4 west,1,5,4 south,2,1,4 east,2,5,4

	// 1st
	r.Name = values[0]

	// 2nd
	p := strings.Split(values[1], ",")
	px1, err := strconv.ParseInt(p[0], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size X1 failed, "
	}
	py1, err := strconv.ParseInt(p[1], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size Y1 failed, "
	}

	r.Pos.X = int(px1)
	r.Pos.Y = int(py1)

	t := strings.Split(values[2], ",")
	tx, err := strconv.ParseInt(t[0], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region tile size X failed, "
	}
	ty, err := strconv.ParseInt(t[1], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region tile size Y failed, "
	}

	r.TileSize.X = int(tx)
	r.TileSize.Y = int(ty)

	errstr = parseAnimStr(values[3:], errstr, r)

	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

// GetFrameRect relies on a anim struct from GetAnimation and a framenumber to return the wanted Rect X,Y,Width,Height (all int)
func (r *Region) GetFrameRect(anim Anim, frameNumber int) (RECT, error) {

	var rect RECT

	rect = RECT{X: 0, Y: 0, Width: r.TileSize.X, Height: r.TileSize.Y}

	//Change Anim Pos in Region Grid to zero based for calc
	//Adjust Offset position of Animation in Region
	offsetX := (anim.Pos.X - 1) * r.TileSize.X
	offsetY := (anim.Pos.Y - 1) * r.TileSize.Y

	rect.X = frameNumber*r.TileSize.X + offsetX + r.Pos.X
	rect.Y = offsetY + r.Pos.Y

	return rect, nil
}

// *********************ANIMATION*********************

// GetAnimation returns the animation sprite sheet info (used with GetFrameRec)
func (r *Region) GetAnimation(animName string) (Anim, error) {
	anim, ok := r.Anims[animName]
	if !ok {
		return Anim{}, errors.New("animation %q not found in region " + r.Name)
	}
	return anim, nil
}

func parseAnimStr(values []string, errstr string, r *Region) string {
	for _, item := range values {
		i := strings.Split(item, ",")
		r1, err := strconv.ParseInt(i[1], 0, 0)
		row := int(r1)
		if err != nil {
			errstr = errstr + fmt.Sprintf("Parse Anim %q row failed, ", i[0])
		}
		c1, err := strconv.ParseInt(i[2], 0, 0)
		col := int(c1)
		if err != nil {
			errstr = errstr + fmt.Sprintf("Parse Anim %q col failed, ", i[0])
		}
		co, err := strconv.ParseInt(i[3], 0, 0)
		count := int(co)
		if err != nil {
			errstr = errstr + fmt.Sprintf("Parse Anim %q count failed, ", i[0])
		}
		pos := XY{row, col}
		step := false
		loop := true
		if len(i) > 4 {
			if i[4] == "step" {
				step = true
			}
			if i[4] == "once" {
				loop = false
			}
		}
		anim := Anim{Pos: pos, Count: count, Step: step, Loop: loop}

		r.Anims[i[0]] = anim
	}
	return errstr
}
