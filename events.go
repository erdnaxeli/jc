package jc

type JoinEvent struct {
	Nick    string
	Channel string
}

type MessageEvent struct {
	Nick    string
	Channel string
	Text    string
}

type NickEvent struct {
	OldNick string
	NewNick string
}

type PrivMessageEvent struct {
	Nick    string
	Channel string
	Text    string
}

type PartEvent struct {
	Nick    string
	Channel string
}

type QuitEvent struct {
	Nick string
}
