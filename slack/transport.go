package slack

import (
	"log"
	"os"

	"github.com/nlopes/slack"

	"git.iiens.net/morignot2011/jc"
)

type Transport struct {
	jc.BaseTransport

	api          *slack.Client
	rtm          *slack.RTM
	channelIDs   map[string]string
	channelNames map[string]string
	userNames    map[string]string
}

func New(name string, cfg map[string]interface{}) (jc.Transport, error) {
	token, ok := cfg["token"]
	if !ok {
		return nil, jc.ConfigError{"token"}
	}

	api := slack.New(token.(string))
	rtm := api.NewRTM()
	t := &Transport{
		api:          api,
		rtm:          rtm,
		channelIDs:   make(map[string]string),
		channelNames: make(map[string]string),
		userNames:    make(map[string]string),

		BaseTransport: jc.BaseTransport{
			Events: make(chan interface{}),
			End:    make(chan bool),

			Logger: log.New(os.Stdout, name+": ", log.LstdFlags),
		},
	}

	channels, err := t.api.GetChannels(true)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if !channel.IsMember {
			continue
		}

		t.channelIDs[channel.Name] = channel.ID
		t.channelNames[channel.ID] = channel.Name
	}

	return t, nil

}
func (t *Transport) Run() error {
	go t.rtm.ManageConnection()
	go t.dispatchEvents()

	return nil
}

func (t *Transport) dispatchEvents() {
	for {
		select {
		case msg := <-t.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ChannelJoinedEvent:
				// TODO
			case *slack.ChannelRenameEvent:
				// TODO
			case *slack.ConnectedEvent:
				t.connected(ev)
			case *slack.MessageEvent:
				t.message(ev)
			case *slack.InvalidAuthEvent:
				t.invalidAuth(ev)
			case *slack.RTMError:
				// Log
			}
		}
	}
}
