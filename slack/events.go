package slack

import (
	"fmt"

	"github.com/nlopes/slack"

	"git.iiens.net/morignot2011/jc"
)

func (t *Transport) connected(ev *slack.ConnectedEvent) {
	events := []*jc.JoinEvent{}

	for name, id := range t.channelIDs {
		channel, err := t.api.GetChannelInfo(id)
		if err != nil {
			t.connectionError <- err
			return
		}

		fmt.Printf("%s: %d\n", name, len(channel.Members))
		for _, id := range channel.Members {
			user, err := t.api.GetUserInfo(id)
			if err != nil {
				t.connectionError <- err
				return
			}
			t.userNames[id] = user.Name

			events = append(events, &jc.JoinEvent{
				Nick:    user.Name,
				Channel: name,
			})
		}
	}

	close(t.connectionError)

	fmt.Printf("joins: %d\n", len(events))
	for _, event := range events {
		t.Events <- event
	}
}

func (t *Transport) invalidAuth(ev *slack.InvalidAuthEvent) {
	t.connectionError <- fmt.Errorf("Invalid auth")
}

func (t *Transport) message(ev *slack.MessageEvent) {
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