package main

import (
	"net"
)

type PrivateMessage struct {
        Messages []string
}

type ChannelUser struct {
        Mode string
}

type Channel struct {
        Mode string
        Topic string
        Users map[string]ChannelUser
        Messages []string
}

type Client struct {
        Nick string
        Secure bool
        ServerHost string
        ServerPort string
        ServerPass string
        ServerPrefix string
        NetIncoming chan string
        NetOutgoing chan string
        Conn net.Conn
        Quit chan bool
        ServerMessages []string
        Channels map[string]Channel
        PrivateMessages map[string]PrivateMessage

}