package videoasr

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

type Ffmpeg struct {
	Debug bool
}

func NewFfmpeg(debug bool) (*Ffmpeg, error) {
	err := ffmpegCheck()
	if err != nil {
		return nil, err
	}

	return &Ffmpeg{debug}, nil
}

func ffmpegCheck() error {
	cmd := exec.Command("ffmpeg", "-version")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("ffmpeg does not exist!")
	}
	return nil
}

func (f *Ffmpeg) ExtractAudio(videof string, audiof string) error {
	if !isFileExisted(videof) {
		return ErrFileNotExisted
	}

	if isFileExisted(audiof) {
		return ErrFileExisted
	}

	cmd := exec.Command("ffmpeg", "-i", videof, "-ar", "16000", audiof)
	if f.Debug {
		log.Println("ExtractAudio cmd is: ffmpeg", "-i", videof, "-ar", "16000", audiof)
	}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func (f *Ffmpeg) RemoveSubtitles(x, y, w, h int, input, output string) error {
	if !isFileExisted(input) {
		return ErrFileNotExisted
	}

	if isFileExisted(output) {
		return ErrFileExisted
	}

	args := []string{
		"-i",
		input,
		"-filter_complex",
		fmt.Sprintf("[0:v]crop=%d:%d:%d:%d,avgblur=10[fg];[0:v][fg]overlay=%d:%d[v]", w, h, x, y, x, y),
		"-map",
		"[v]",
		"-map",
		"0:a",
		"-c:v",
		"libx264",
		"-c:a",
		"copy",
		"-movflags",
		"+faststart",
		output,
	}

	if f.Debug {
		log.Println("cmd is: ffmpeg ", args)
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("CombinedOutput err: %s", err)
	}
	return nil
}

// ffmpeg -i  ~/test_1.mp4 -vf subtitles=./test_1.srt  -sn ./outfile_test_1.mp4
func (f *Ffmpeg) AddSubtitles(inputVideo, inputSrt, outputVideo string, area *Area) error {
	if !isFileExisted(inputVideo) || !isFileExisted(inputSrt) {
		return ErrFileNotExisted
	}

	if isFileExisted(outputVideo) {
		return ErrFileExisted
	}

	args := []string{
		"-i",
		inputVideo,
		"-vf",
		// fmt.Sprintf("subtitles=%s", inputSrt),
		fmt.Sprintf("subtitles=%s:force_style='Alignment=2,MarginV=%d", inputSrt, area.BottomDistance()),
		"-sn",
		outputVideo,
	}

	if f.Debug {
		log.Println("cmd is: ffmpeg ", args)
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("CombinedOutput err: %s", err)
	}
	return nil
}
