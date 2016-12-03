package slack

import (
	"fmt"

	"github.com/nlopes/slack"

	"git.iiens.net/morignot2011/jc"
)

func (t *Transport) connected(ev *slack.ConnectedEvent) {
	t.Logger.Printf("Connected event")
	events := []*jc.JoinEvent{}

	for name, id := range t.channelIDs {
		channel, err := t.api.GetChannelInfo(id)
		if err != nil {
			t.Logger.Print(err)
			return
		}

		for _, id := range channel.Members {
			user, err := t.api.GetUserInfo(id)
			if err != nil {
				t.Logger.Printf(err)
				return
			}
			t.userNames[id] = user.Name

			events = append(events, &jc.JoinEvent{
				Nick:    user.Name,
				Channel: name,
			})
		}
	}

	fmt.Printf("joins: %d\n", len(events))
	for _, event := range events {
		t.Events <- event
	}
}

func (t *Transport) invalidAuth(ev *slack.InvalidAuthEvent) {
}

func (t *Transport) message(ev *slack.MessageEvent) {
	t.Logger.Printf("Message event")
	channel, ok := t.channelNames[ev.Channel]
	if !ok {
		// probably we got invited on a channel after starting the bot
		return
	}

	nick, ok := t.userNames[ev.User]
	if !ok {
		// probably a user which have join after connection
		return
	}

	t.Events <- &jc.MessageEvent{
		Nick:    nick,
		Channel: channel,
		Text:    ev.Text,
	}
}
