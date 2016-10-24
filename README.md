# JC

Bot to sync Slack or IRC channels together.

## Usage

```
./jc [--config file]
```

## Configuration

By default JC reads the file jc.conf in the current directory.

Example:
```
transports:
    irc.iiens.net:
        type: irc
        server: irc.iiens.net
        port: 7000
        ssl: true
        ssl_insecure: true
        channels:
            - "#list"
            - "#of"
            - "#channels"
    yourteam.slack.com:
        type: slack
        token: secret_bot_token
links:
    - irc.iiens.net:
        - "#channels"
      yourteam.slack.com:
        - aslackchannel
      filters:
        - some
        - nicks
        - to
        - ignore
```

## Build

```
go get git.iiens.net/morignot2011/jc/cmd/jc
cd $GOPATH/src/git.iiens.net/morignot2011/jc/cmd/jc
go build
```
