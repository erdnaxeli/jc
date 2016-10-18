package irc

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	irc "github.com/fluffle/goirc/client"

	"git.iiens.net/morignot2011/jc"
)

const (
	DEFAULT_NICK = "jc"
	DEFAULT_SSL  = false
)

type Transport struct {
	jc.BaseTransport

	cfg          map[string]interface{}
	channels     []string
	client       *irc.Conn
	realNicks    map[string]string
	userClients  map[string]*irc.Conn
	userChannels map[string][]string

	connectionError chan error
}

func New(name string, cfg map[string]interface{}) (jc.Transport, error) {
	t := &Transport{
		cfg:             cfg,
		connectionError: make(chan error),
		userClients:     make(map[string]*irc.Conn),
		userChannels:    make(map[string][]string),

		BaseTransport: jc.BaseTransport{
			Events: make(chan interface{}),
			End:    make(chan bool),

			Logger: log.New(os.Stdout, name+": ", log.LstdFlags),
		},
	}

	ircCfg, err := t.newIrcConfig(DEFAULT_NICK, cfg)
	if err != nil {
		return nil, err
	}

	t.client = t.getIrcClient(ircCfg)

	channels, ok := cfg["channels"]
	if !ok {
		return nil, jc.ConfigError{"channels"}
	}
	for _, channel := range channels.([]interface{}) {
		t.channels = append(t.channels, channel.(string))
	}

	return t, nil
}

func (t *Transport) Run() error {
	if err := t.client.Connect(); err != nil {
		return jc.ConnectionError{err}
	}

	if err, ok := <-t.connectionError; ok {
		return jc.ConnectionError{err}
	}

	return nil
}

func (t *Transport) newIrcConfig(nick string, cfg map[string]interface{}) (*irc.Config, error) {
	server, ok := cfg["server"]
	if !ok {
		return nil, jc.ConfigError{"server"}
	}

	port, ok := cfg["port"]
	if !ok {
		return nil, jc.ConfigError{"port"}
	}

	ssl, ok := cfg["ssl"]
	if !ok {
		ssl = DEFAULT_SSL
	}

	insecure, ok := cfg["ssl_insecure"]
	if !ok {
		insecure = false
	}

	ircCfg := irc.NewConfig(nick)
	ircCfg.SSL = ssl.(bool)
	if ircCfg.SSL {
		ircCfg.SSLConfig = &tls.Config{
			ServerName:         server.(string),
			InsecureSkipVerify: insecure.(bool),
		}
	}
	ircCfg.Server = fmt.Sprintf("%s:%d", server, port)
	ircCfg.NewNick = t.newNick

	return ircCfg, nil
}

func (t *Transport) newNick(nick string) string {
	log.Printf("new nick %s", nick)
	realNick, ok := t.realNicks[nick]
	if !ok {
		realNick = nick
	} else {
		delete(t.realNicks, nick)
	}

	newNick := nick + "_"
	t.realNicks[newNick] = realNick
	return newNick
}

func (t *Transport) getIrcClient(cfg *irc.Config) *irc.Conn {
	client := irc.Client(cfg)
	client.EnableStateTracking()
	client.HandleFunc(irc.CONNECTED, t.connected)
	client.HandleFunc(irc.DISCONNECTED, t.disconnected)
	client.HandleFunc(irc.JOIN, t.join)
	client.HandleFunc(irc.NICK, t.nick)
	client.HandleFunc(irc.PART, t.part)
	client.HandleFunc(irc.PRIVMSG, t.privmsg)
	client.HandleFunc(irc.QUIT, t.quit)

	return client
}

func (t *Transport) isUserDistant(user string) bool {
	for k, _ := range t.userClients {
		if k == user {
			return true
		}
	}

	return false
}
