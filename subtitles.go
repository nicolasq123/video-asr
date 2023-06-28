package videoasr

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Area struct {
	w, h, x, y  int
	VideoHeight int
	VideoWidth  int
}

// 距离底部的距离
func (c *Area) BottomDistance() int {
	return max(c.VideoHeight-c.h, 0)
}

type SubtitleConf struct {
	PythonPathF string
	Arg         string
	Debug       bool
}

func (c *SubtitleConf) New() (*Subtitle, error) {
	if !isFileExisted(c.PythonPathF) {
		err := fileStat(c.PythonPathF)
		fmt.Println("err: ", err.Error())
		return nil, fmt.Errorf("PythonPathF does not existed: %s", c.PythonPathF)
	}
	return &Subtitle{
		c: c,
	}, nil
}

type Subtitle struct {
	c *SubtitleConf
}

func (s *Subtitle) Find(videof string) (*Area, error) {
	cmd := exec.Command(s.c.PythonPathF, s.c.Arg, "--input", videof)
	if s.c.Debug {
		log.Println("cmd is: ", s.c.PythonPathF, s.c.Arg, "--input", videof)
	}
	bs, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("CombinedOutput error: %s", string(bs))
		return nil, err
	}
	bss := bytes.Split(bs, []byte{'('})
	nbs := bss[len(bss)-1]
	nbss := bytes.Split(nbs, []byte{','})
	if len(nbss) != 6 {
		fmt.Println("len nbss", len(nbss), string(nbs))
		err = fmt.Errorf("result is error: %s", string(bs))
		return nil, err
	}
	area, err := parseArea(nbss)
	return area, err
}

func parseArea(bs [][]byte) (*Area, error) {
	do := func(b []byte) (int, error) {
		s := strings.TrimSpace(string(b))
		s = strings.TrimPrefix(s, "(")
		s = strings.TrimSuffix(s, ")")
		s = strings.TrimSpace(s)
		return strconv.Atoi(s)
	}

	x, err := do(bs[0])
	if err != nil {
		return nil, err
	}

	y, err := do(bs[1])
	if err != nil {
		return nil, err
	}

	w, err := do(bs[2])
	if err != nil {
		return nil, err
	}

	h, err := do(bs[3])
	if err != nil {
		return nil, err
	}

	height, err := do(bs[4])
	if err != nil {
		return nil, err
	}

	width, err := do(bs[3])
	if err != nil {
		return nil, err
	}

	return &Area{
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		VideoHeight: height,
		VideoWidth:  width,
	}, nil
}
