package spriteatlas

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
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

	return errors.New(errstr)
}

func (p *Page) PageToStr() string {
	return fmt.Sprintf("page %s has default tilesize %d %d", p.name, p.tile_size.X, p.tile_size.Y)
}

// Spriteatlas reads a atlas at Path + Name ~ use forward slashe(s) in path.
// reads the spritesheet as an image and will 'overwrite' alpha color if specified and found.
// alpha-color is NOT the same as pre-multiplied-alpha
func Spriteatlas(filePath string, fileName string) error {
	if filePath[len(filePath)-1] != '/' {
		filePath = filePath + "/"
	}

	file, reader, err := OpenAtlas(filePath + fileName)
	if err != nil && err != io.EOF {
		os.Exit(1)
	}
	defer file.Close()
	var line string
	line, err = ParseAtlas(reader)
	println(line, err)

	return err
}

func OpenAtlas(name string) (*os.File, *bufio.Reader, error) {

	file, err := os.Open(name)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	reader := bufio.NewReader(file)

	return file, reader, err

}

func ParseAtlas(reader *bufio.Reader) (string, error) {
	var str string
	var pgErr error
	var page Page

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		str = StripAtlasLine(line)
		if str == "" || str[:1] == "#" || str[:2] == "//" {
			continue
		}
		a := strings.Split(str, " ")
		for i := 0; i < len(a); i++ {
			a[i] = strings.Trim(a[i], " ")
		}
		if a[0] == "page" {
			pgErr = page.ParsePageStr(a[1:])
		}
		break
	}
	pgStr := page.PageToStr()
	return pgStr, pgErr
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
