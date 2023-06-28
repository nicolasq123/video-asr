package videoasr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type Translator interface {
	Translate(ctx context.Context, inputs []string, locale string) (rt []string, err error)
	Close()
}

type TranslatorConf struct {
	GoogleTransConf *GoogleTransConf
}

func (c *TranslatorConf) New() (Translator, error) {
	return c.GoogleTransConf.New()
}

type GoogleTransConf struct {
	CredentialFile string
}

func (c *GoogleTransConf) New() (*GoogleTrans, error) {
	abs, err := filepath.Abs(c.CredentialFile)
	if err != nil {
		return nil, err
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", abs)
	if err != nil {
		return nil, err
	}

	client, err := translate.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &GoogleTrans{client}, nil
}

type GoogleTrans struct {
	client *translate.Client
}

func (g *GoogleTrans) Translate(ctx context.Context, texts []string, targetLang string) (rt []string, err error) {
	lang, err := language.Parse(targetLang)
	if err != nil {
		err = fmt.Errorf("language.Parse: %v", err)
		return
	}
	rt, err = g.translateHelp(ctx, texts, lang)
	return
}

func (g *GoogleTrans) Close() {
	g.client.Close()
}

func (g *GoogleTrans) translateHelp(ctx context.Context, texts []string, lang language.Tag) (rt []string, err error) {
	length := len(texts)
	opts := translate.Options{Format: translate.Text}
	for i := 0; i < length; i += 100 {
		end := min(i+100, length)
		var resp []translate.Translation
		resp, err = g.client.Translate(ctx, texts[i:end], lang, &opts)
		if err != nil {
			err = fmt.Errorf("Translate: %v", err)
			return
		}
		if len(resp) != end-i {
			err = fmt.Errorf("Translate returned error: %d != %d", len(resp), end-i)
			return
		}
		for _, v := range resp {
			rt = append(rt, v.Text)
		}
	}
	return
}
