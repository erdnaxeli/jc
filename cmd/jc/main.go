package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"git.iiens.net/morignot2011/jc"
	"git.iiens.net/morignot2011/jc/dispatcher"
	"git.iiens.net/morignot2011/jc/irc"
	"git.iiens.net/morignot2011/jc/slack"
)

type Cfg struct {
	Transports map[string]map[string]interface{}
	Links      []map[string]interface{}
}

type Links map[string]Link
type Link map[string]LinkDest

type LinkDest struct {
	transport string
	channel   string
	filter    []string
}

var links = make(Links)

func main() {
	var cfgFile = flag.String("config", "jc.conf", "path to configuration file")
	flag.Parse()

	file, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Cfg
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	transports := make(map[string]jc.Transport)

	for name, transport := range cfg.Transports {
		type_, ok := transport["type"]
		if !ok {
			log.Fatalf("Missing 'type' for transport %s", name)
		}

		delete(transport, "type")

		var t jc.Transport
		var err error

		switch type_.(string) {
		case "irc":
			t, err = irc.New(name, transport)
		case "slack":
			t, err = slack.New(name, transport)
		}

		if err != nil {
			log.Fatalf("Error when creating transport %s: %s", name, err)
		}

		transports[name] = t
		err = t.Run()
		if err != nil {
			log.Fatalf("Error when connecting transport %s: %s", name, err)
		}
	}

	dispatcher := dispatcher.NewDispatcher(cfg.Links, transports)

	// remove sensitive data
	cfg = Cfg{}

	dispatcher.Run()
}
