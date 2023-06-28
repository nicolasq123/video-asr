package videoasr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-speech-sdk-go/asr"
	"github.com/tencentcloud/tencentcloud-speech-sdk-go/common"
)

var (
	EngineType = "16k_zh"
)

type MyAsr interface {
	GenSrt(ctx context.Context, srtfile, audiofule string) (string, error)
	GenSrtAndTranslate(ctx context.Context, srtfile, transsrtfile string, audiofile string, trans Translator, locale string) (string, error)
}

type MyAsrConf struct {
	TxAsrConf
}

func (c *MyAsrConf) New() (MyAsr, error) {
	// 接入其他asr
	return c.TxAsrConf.New()
}

type TxAsrConf struct {
	AppID     string `validate:"nonzero"`
	SecretID  string `validate:"nonzero"`
	SecretKey string `validate:"nonzero"`
}

func (c *TxAsrConf) New() (*TxAsr, error) {
	return &TxAsr{
		conf: c,
	}, nil
}

func (c *TxAsrConf) Validate() error {
	return nil
}

type TxAsr struct {
	conf *TxAsrConf
}

func (a *TxAsr) GenSrt(ctx context.Context, srtfile, audiofule string) (res string, err error) {
	resp, err := a.req(audiofule)
	if err != nil {
		return "", err
	}
	fmt.Printf("request_id: %s\n", resp.RequestId)

	if resp == nil || len(resp.FlashResult) == 0 {
		err = fmt.Errorf("flash result is zero")
		return
	}
	genSrt(srtfile, resp.FlashResult[0].SentenceList)

	bs, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	res = string(bs)
	return
}

func (a *TxAsr) GenSrtAndTranslate(ctx context.Context, srtfile, transsrtfile, audiofule string, trans Translator, locale string) (res string, err error) {
	if isFileExisted(srtfile) || isFileExisted(transsrtfile) {
		return "", ErrFileExisted
	}

	resp, err := a.req(audiofule)
	if err != nil {
		return "", err
	}
	fmt.Printf("request_id: %s\n", resp.RequestId)

	if resp == nil || len(resp.FlashResult) == 0 {
		err = fmt.Errorf("flash result is zero")
		return
	}
	genSrt(srtfile, resp.FlashResult[0].SentenceList)

	texts := []string{}
	for _, sl := range resp.FlashResult[0].SentenceList {
		texts = append(texts, sl.Text)
	}
	rt, err := trans.Translate(ctx, texts, locale)
	if err != nil {
		err = fmt.Errorf("trans.Translate err: %v", err)
		return
	}

	if len(rt) != len(texts) {
		err = fmt.Errorf("trans.Translate length err: %v %v", len(rt), len(texts))
		return
	}

	for i := range resp.FlashResult[0].SentenceList {
		resp.FlashResult[0].SentenceList[i].Text = rt[i]
	}

	genSrt(transsrtfile, resp.FlashResult[0].SentenceList)

	bs, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	res = string(bs)
	return
}

func (a *TxAsr) req(audiofule string) (resp *asr.FlashRecognitionResponse, err error) {
	audio, err := os.Open(audiofule)
	defer audio.Close()
	if err != nil {
		err = fmt.Errorf("open file error: %v\n", err)
		return
	}

	credential := common.NewCredential(a.conf.SecretID, a.conf.SecretKey)
	recognizer := asr.NewFlashRecognizer(a.conf.AppID, credential)
	data, err := ioutil.ReadAll(audio)
	if err != nil {
		err = fmt.Errorf("%s|failed read data, error: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return
	}

	req := new(asr.FlashRecognitionRequest)
	req.EngineType = EngineType
	req.VoiceFormat = "mp3"
	req.SpeakerDiarization = 0
	req.FilterDirty = 0
	req.FilterModal = 0
	req.FilterPunc = 0
	req.ConvertNumMode = 1
	req.FirstChannelOnly = 1
	req.SentenceMaxLength = 10
	req.WordInfo = 3

	resp, err = recognizer.Recognize(req, data)
	if err != nil {
		err = fmt.Errorf("%s|failed do recognize, error: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return
	}
	return
}

func genSrt(f string, sl []*asr.FlashRecognitionSentence) {
	file, e := os.Create(f)
	if e != nil {
		panic(e)
	}
	defer file.Close()
	for i, s := range sl {
		linestr := MakeSubtitleText(i, int64(s.StartTime), int64(s.EndTime), s.Text)
		file.WriteString(linestr)
	}
}

func MakeSubtitleText(index int, startTime int64, endTime int64, text string) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\n")
	content.WriteString(SubtitleTimeMillisecond(startTime))
	content.WriteString(" --> ")
	content.WriteString(SubtitleTimeMillisecond(endTime))
	content.WriteString("\n")
	content.WriteString(text)
	content.WriteString("\n")
	content.WriteString("\n")
	return content.String()
}

func SubtitleTimeMillisecond(time int64) string {
	var miao int64 = 0
	var min int64 = 0
	var hours int64 = 0
	var millisecond int64 = 0

	millisecond = (time % 1000)
	miao = (time / 1000)

	if miao > 59 {
		min = (time / 1000) / 60
		miao = miao % 60
	}
	if min > 59 {
		hours = (time / 1000) / 3600
		min = min % 60
	}

	//00:00:06,770
	var miaoText = RepeatStr(strconv.FormatInt(miao, 10), "0", 2, true)
	var minText = RepeatStr(strconv.FormatInt(min, 10), "0", 2, true)
	var hoursText = RepeatStr(strconv.FormatInt(hours, 10), "0", 2, true)
	var millisecondText = RepeatStr(strconv.FormatInt(millisecond, 10), "0", 3, true)

	return hoursText + ":" + minText + ":" + miaoText + "," + millisecondText
}

func RepeatStr(str string, s string, length int, before bool) string {
	ln := len(str)

	if ln >= length {
		return str
	}

	if before {
		return strings.Repeat(s, (length-ln)) + str
	} else {
		return str + strings.Repeat(s, (length-ln))
	}
}
