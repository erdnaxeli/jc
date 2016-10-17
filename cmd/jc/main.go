package main

import (
	"flag"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"git.iiens.net/morignot2011/jc"
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

var transports = make(map[string]jc.Transport)
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
			t, err = irc.New(transport)
		case "slack":
			t, err = slack.New(transport)
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

	for _, link := range cfg.Links {
		filter, ok := link["filter"]
		if !ok {
			filter = []string{}
		}

		delete(link, "filter")

		for a, b := range link {
			cfgA := strings.Split(a, "@")
			cfgB := strings.Split(b.(string), "@")

			if len(cfgA) != 2 || len(cfgB) != 2 {
				log.Fatalf("Invalid link '%s => %s', it must be at the format 'channelA@transportA => channelB@transportB'", a, b.(string))
			}

			if _, ok := transports[cfgA[1]]; !ok {
				log.Fatalf("Unknown transport %s in link '%s => %s'", cfgA[1], a, b.(string))
			}

			if _, ok := transports[cfgB[1]]; !ok {
				log.Fatalf("Unknown transport %s in link '%s => %s'", cfgB[1], a, b.(string))
			}

			var filters []string
			for _, f := range filter.([]interface{}) {
				filters = append(filters, f.(string))
			}

			if links[cfgA[1]] == nil {
				links[cfgA[1]] = make(Link)
			}
			links[cfgA[1]][cfgA[0]] = LinkDest{
				transport: cfgB[1],
				channel:   cfgB[0],
				filter:    filters,
			}

			if links[cfgB[1]] == nil {
				links[cfgB[1]] = make(Link)
			}
			links[cfgB[1]][cfgB[0]] = LinkDest{
				transport: cfgA[1],
				channel:   cfgA[0],
				filter:    filters,
			}
		}
	}

	// remove sensitive data
	cfg = Cfg{}

	cases := make([]reflect.SelectCase, len(transports)*2)
	i := 0
	for _, t := range transports {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.GetEvents())}
		i++
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.GetEnd())}
		i++
	}

	for {
		idx, value, ok := reflect.Select(cases)

		i := 0
		var name string
		for name, _ = range transports {
			if i == idx || i+1 == idx {
				break
			}

			i += 2
		}

		if !ok {
			log.Fatalf("Transport %s close one of its channel!", name)
		}

		switch ev := value.Interface().(type) {
		case *jc.JoinEvent:
			dest, ok := links[name][ev.Channel]
			if !ok {
				continue
			}

			if isFiltered(dest.filter, ev.Nick) {
				continue
			}

			t := transports[dest.transport]

			go t.Join(&jc.JoinEvent{
				Nick:    ev.Nick + "_jc",
				Channel: dest.channel,
			})
		case *jc.MessageEvent:
			dest, ok := links[name][ev.Channel]
			if !ok {
				continue
			}

			if isFiltered(dest.filter, ev.Nick) {
				continue
			}

			t := transports[dest.transport]

			go t.Message(&jc.MessageEvent{
				Nick:    ev.Nick + "_jc",
				Channel: dest.channel,
				Text:    ev.Text,
			})
		case bool:
			log.Fatalf("Disconnected from %s", name)
			return
		}
	}
}

func isFiltered(filter []string, nick string) bool {
	for _, v := range filter {
		if v == nick {
			return true
		}
	}

	return false
}
