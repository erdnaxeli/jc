package slack

import (
	"fmt"

	"git.iiens.net/morignot2011/jc"
)

func (t *Transport) Join(ev *jc.JoinEvent) {
	text := fmt.Sprintf("%s has joined the channel", ev.Nick)
	channel := t.channelIDs[ev.Channel]
	t.rtm.SendMessage(t.rtm.NewOutgoingMessage(text, channel))
}

func (t *Transport) Message(ev *jc.MessageEvent) {
	text := fmt.Sprintf("<%s> %s", ev.Nick, ev.Text)
	channel := t.channelIDs[ev.Channel]
	t.rtm.SendMessage(t.rtm.NewOutgoingMessage(text, channel))
}

func (t *Transport) Nick(ev *jc.NickEvent) {
	text := fmt.Sprintf("%s is now known as%s", ev.OldNick, ev.NewNick)

	for _, channel := range t.channelIDs {
		t.rtm.SendMessage(t.rtm.NewOutgoingMessage(text, channel))
	}
}

func (t *Transport) PrivMessage(ev *jc.PrivMessageEvent) {
	// TODO
}

func (t *Transport) Part(ev *jc.PartEvent) {
	text := fmt.Sprintf("%s has left", ev.Nick)
	channel := t.channelIDs[ev.Channel]
	t.rtm.SendMessage(t.rtm.NewOutgoingMessage(text, channel))
}

func (t *Transport) Quit(ev *jc.QuitEvent) {
	text := fmt.Sprintf("%s has quit", ev.Nick)

	for _, channel := range t.channelIDs {
		t.rtm.SendMessage(t.rtm.NewOutgoingMessage(text, channel))
	}
}
