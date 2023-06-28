package videoasr

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
)

//var ErrJobNonSelect = errors.New("only last job can be insert/update/delete")

type Conf struct {
	DSN          string
	WriterType   string
	Debug        bool
	Asr          MyAsrConf
	SubtitleConf SubtitleConf

	TranslatorConf TranslatorConf
	DefaultLocale  string

	InputVideoFile string

	Mode string // skip: 需要生成文件已存在，就不再生成了; normal： 已存在就正常报错; defult normal
}

type Asr struct {
	conf       *Conf
	db         *sqlx.DB
	myasr      MyAsr
	ff         *Ffmpeg
	subt       *Subtitle
	translator Translator
}

func (c *Conf) New() (*Asr, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	myasr, err := c.Asr.New()
	if err != nil {
		return nil, err
	}
	subt, err := c.SubtitleConf.New()
	if err != nil {
		return nil, err
	}

	translator, err := c.TranslatorConf.New()
	if err != nil {
		return nil, err
	}
	ff, err := NewFfmpeg(c.Debug)
	if err != nil {
		return nil, err
	}
	asr := &Asr{
		conf:       c,
		myasr:      myasr,
		subt:       subt,
		ff:         ff,
		translator: translator,
	}

	if c.DSN != "" {
		db, err := Open(c.DSN)
		if err != nil {
			return nil, err
		}
		asr.db = db
	}

	return asr, nil
}

func (c *Conf) Validate() error {
	return nil
}

func (a *Asr) Close() {
	a.translator.Close()
}

func (a *Asr) Run() (err error) {
	videof := a.conf.InputVideoFile
	if !isFileExisted(videof) {
		return fmt.Errorf("videof file does not existed: %s", videof)
	}

	audiof := genFilename(videof, "mp3")
	videooutputf := genFilename(videof, "mp4")
	srtf := genFilename(videof, "srt0")
	transSrtFile := genFilename(videof, "srt")
	realVideoOutputF := genFilenameV2(videof, "mp4")

	ctx, cancel := context.WithTimeout(context.Background(), 480*time.Second)
	defer cancel()

	err = a.ff.ExtractAudio(videof, audiof)
	if a.isSeriousError(err) {
		return
	}
	var msg string
	msg, err = a.myasr.GenSrtAndTranslate(ctx, srtf, transSrtFile, audiof, a.translator, a.conf.DefaultLocale)
	if a.isSeriousError(err) {
		return
	}
	_ = msg
	area, err := a.subt.Find(videof)
	if a.isSeriousError(err) {
		return
	}
	err = a.ff.RemoveSubtitles(area.x, area.y, area.w, area.h, videof, videooutputf)
	if a.isSeriousError(err) {
		return
	}
	err = a.ff.AddSubtitles(videooutputf, transSrtFile, realVideoOutputF, area)
	if a.isSeriousError(err) {
		return
	}
	return
}

func (a *Asr) isSeriousError(err error) bool {
	if err == nil {
		return false
	}
	if a.conf.Mode == "skip" {
		if err == ErrFileExisted {
			return false
		}
	}
	return true
}

func genFilenameV2(videof string, ext string) string {
	dir := filepath.Dir(videof)
	base := filepath.Base(videof)
	name, _ := parseFilename(base)
	name = fmt.Sprintf("%s_outputv2.%s", name, ext)
	return filepath.Join(dir, name)
}

func genFilename(videof string, ext string) string {
	dir := filepath.Dir(videof)
	base := filepath.Base(videof)
	name, _ := parseFilename(base)
	name = fmt.Sprintf("%s_output.%s", name, ext)
	return filepath.Join(dir, name)
}

func parseFilename(f string) (string, string) {
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '.' {
			return f[0:i], f[i:]
		}
	}
	return f, ""
}
