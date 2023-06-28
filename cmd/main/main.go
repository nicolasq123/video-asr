package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/nicolasq123/videoasr"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	c := parseConf()
	asr, err := c.New()
	if err != nil {
		panic(err)
	}
	err = asr.Run()

	if err != nil {
		panic(err)
	}
}

func parseConf() *videoasr.Conf {
	c := &videoasr.Conf{}
	confPath := flag.String("conf", "./conf.yml", "path to config")
	flag.Parse()
	b, err := ioutil.ReadFile(*confPath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, c)

	if err != nil {
		panic(err)
	}

	log.Printf("conf is: %v \n", c)

	out, _ := yaml.Marshal(c)
	os.Stdout.Write(out)
	return c
}
