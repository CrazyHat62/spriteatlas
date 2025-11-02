package spriteatlas

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type XY struct {
	X int
	Y int
}

type RECT struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type Anim struct {
	Pos0  XY
	Count int
}

type Region struct {
	Name         string
	RegionBounds RECT
	TileSize     XY
	Anims        map[string]Anim
}

func (r *Region) ParseRegionStr(values []string) error {
	var err error
	var errstr string = ""
	r.Anims = make(map[string]Anim)
	//player_walk 1,148,384,244 48,48 north,1,1,4 west,1,5,4 south,2,1,4 east,2,5,4
	r.Name = values[0]
	p := strings.Split(values[1], ",")
	px1, err := strconv.ParseInt(p[0], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size X1 failed, "
	}
	py1, err := strconv.ParseInt(p[1], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size Y1 failed, "
	}
	px2, err := strconv.ParseInt(p[2], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size X2 failed, "
	}
	py2, err := strconv.ParseInt(p[3], 0, 0)
	if err != nil {
		errstr = errstr + "Parse Region rect size Y2 failed, "
	}
	r.RegionBounds.X1 = int(px1)
	r.RegionBounds.X2 = int(px2)
	r.RegionBounds.Y1 = int(py1)
	r.RegionBounds.Y2 = int(py2)

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

	for _, item := range values[3:] {
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
		anim := Anim{Pos0: pos, Count: count}
		r.Anims[i[0]] = anim
	}

	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

func (r *Region) RegionToStr() string {
	return fmt.Sprintf("region %s has %d animations", r.Name, len(r.Anims))
}

func (r *Region) AnimKeys() []string {

	var keys []string
	for key := range r.Anims {
		keys = append(keys, key)
	}
	return keys
}

type Page struct {
	name             string
	alpha_color      string
	imageRegionMarks bool
	sheetSize        XY
	margin           RECT
}

func (p *Page) ParsePageStr(values []string) error {
	var err error
	var errstr string = ""
	p.name = values[0]
	p.alpha_color = values[1]
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
	p.margin.X1 = int(mx1)
	p.margin.X2 = int(mx2)
	p.margin.Y1 = int(my1)
	p.margin.Y2 = int(my2)

	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

func (p *Page) PageToStr() string {
	return fmt.Sprintf("page %s has alphacolor=%s and sheetSize %d %d with margins %d,%d,%d,%d", p.name, p.alpha_color, p.sheetSize.X, p.sheetSize.Y, p.margin.X1, p.margin.Y1, p.margin.X2, p.margin.Y2)
}

var page Page
var region Region

// Spriteatlas reads a atlas at Path + Name ~ use forward slash(s) in path.
// reads the spritesheet as an image and will 'overwrite' alpha color if specified and found.
// alpha-color is NOT the same as pre-multiplied-alpha
func Spriteatlas(filePath string, fileName string) (*Page, *Region, error) {

	if len(filePath) > 0 && filePath[len(filePath)-1] != '/' {
		filePath = filePath + "/"
	}

	fileBuf, err := os.ReadFile(filePath + fileName)
	if err != nil {
		return &page, &region, err
	}

	err1 := ParseAtlas(fileBuf)

	return &page, &region, err1
}

func ParseAtlas(fileBytes []byte) error {
	var str string
	var err error

	for line := range bytes.Lines(fileBytes) {

		str = StripAtlasLine(line)
		if str == "" || str[:1] == "#" || str[:2] == "//" {
			continue
		}
		a := strings.Split(str, " ")
		for i := 0; i < len(a); i++ {
			a[i] = strings.Trim(a[i], " ")
		}
		if a[0] == "page" {
			err = page.ParsePageStr(a[1:])
		}
		if a[0] == "region" {
			err = region.ParseRegionStr(a[1:])
		}

	}

	//println(page.PageToStr(), region.RegionToStr(), err)

	return err
}

func StripAtlasLine(line []byte) string {
	//remove newline and carrige return if it exists from line
	line = bytes.ReplaceAll(line, []byte("\r"), []byte(""))
	line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
	//line = bytes.ReplaceAll(line, []byte(" "), []byte(""))

	s := string(line)
	s = strings.Trim(s, " ")
	return s
}

// GetAnimation gets the RECT for the given animation name and frame index in that animation
func GetAnimation(region *Region, animName string, idx int) (RECT, error) {
	anim, ok := region.Anims[animName]
	var rect RECT

	if !ok {
		return rect, errors.New("animation %q not found in region " + region.Name)
	}
	// get correct frame number in animation
	idx = idx % anim.Count

	//rectangle within region for 1st frame + plus region offset
	rect.X1 = (anim.Pos0.X-1)*region.TileSize.X + region.RegionBounds.X1
	rect.Y1 = (anim.Pos0.Y-1)*region.TileSize.Y + region.RegionBounds.Y1
	//calculate start position of frame within animation
	startPosX := idx * region.TileSize.X
	rect.X1 += startPosX
	rect.X2 = rect.X1 + region.TileSize.X
	rect.Y2 = rect.Y2 + region.TileSize.Y

	// idxOffsetX := idx*region.TileSize.X
	// //
	// startPosX := (anim.Pos0.X * idx) * region.TileSize.X
	// startPosY := (anim.Pos0.Y - 1) * region.TileSize.Y

	return rect, nil
}
