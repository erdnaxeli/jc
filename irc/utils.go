package irc

func FindChannel(channels []string, needle string) int {
	for i, channel := range channels {
		if channel == needle {
			return i
		}
	}

	return -1
}
