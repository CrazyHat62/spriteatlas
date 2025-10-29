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
	Name       string
	RegionSize RECT
	TileSize   XY
	Anims      map[string]Anim
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
	r.RegionSize.X1 = int(px1)
	r.RegionSize.X2 = int(px2)
	r.RegionSize.Y1 = int(py1)
	r.RegionSize.Y2 = int(py2)

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
	tile_size        XY
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

	p.tile_size.X = int(px)
	p.tile_size.Y = int(py)
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
	return fmt.Sprintf("page %s has default tilesize %d %d", p.name, p.tile_size.X, p.tile_size.Y)
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
