package main

import (
  "fmt"
  "net"
  "strings"
  "strconv"
  "time"
)

func timeStamp() string {
	t := time.Now().Local()
	return "["+strconv.Itoa(t.Hour())+":"+strconv.Itoa(t.Minute())+"]"
}

func (client *Client) gotChannelTopic(channel,topic string) {
	c := client.Channels[channel]
	c.Topic = topic
	client.Channels[channel] = c
	client.addChannelMessage(channel,"topic: "+topic)
}

func (client *Client) privateMessageReceived(nick,message string) {
	client.addPrivateMessageUser(nick)
	v := client.PrivateMessages[nick]
	v.Messages = append(client.PrivateMessages[nick].Messages,timeStamp()+" "+nick+": "+message)
	client.PrivateMessages[nick] = v
	Log("Private message from "+nick)
	Log(message)
}

func (client *Client) addPrivateMessageUser(nick string) {
  if _,ok := client.PrivateMessages[nick]; !ok {
	client.PrivateMessages[nick] = PrivateMessage{make([]string,0)}
  }
}

func (client *Client) addServerMessage(message string) {
	client.ServerMessages = append(client.ServerMessages,message)
	fmt.Print("*unhandled/server message: "+message)
	Log(strconv.Itoa(len(client.ServerMessages))+"server messages now")

}


func (client *Client) addUserToChannel(channel,nick,mode string) {
	client.addChannelIfNotExists(channel)
	if _,ok := client.Channels[channel].Users[nick]; !ok {
		// Log("adding "+nick+" to "+channel)
		client.Channels[channel].Users[nick] = ChannelUser{mode}
	}
}

func (client *Client) addChannelMessage(channel,message string) {
	client.addChannelIfNotExists(channel)

	v := client.Channels[channel]
	v.Messages = append(client.Channels[channel].Messages,timeStamp()+" "+message)
	client.Channels[channel] = v
	fmt.Println(channel+" message: "+message)
}

func (client *Client) addChannelIfNotExists(channel string) {
  if _,ok := client.Channels[channel]; !ok {
		// Log("adding channel "+channel)
		client.Channels[channel] = Channel{"","",make(map[string]ChannelUser,0),make([]string,0)}
	}
}

func (client *Client) removeChannel(channel string) {
	delete(client.Channels,channel)
}

func (client *Client) removeUserFromChannel(channel,nick string) {
  if _,ok := client.Channels[channel].Users[nick]; ok {
	// Log("deleting "+nick+" from "+channel)
	delete(client.Channels[channel].Users,nick)
  }
}


func (client *Client) parseCommandMessage(message,channel,user string) {
	reply := ""

	bits := strings.Fields(message)
	if bits[0] == client.Nick || bits[0] == client.Nick+":" {
		message = removeFirstWord(message)
	}

	bits = strings.Fields(message)
	
	switch bits[0] {
	default:
		// show help as default?
		reply = "unknown command"

	case "dig":
		if len(bits) <= 1 {
			reply = "please give something to dig"
		} else {
			result, err := net.LookupHost(bits[1])
			if err != nil {
				reply = "dig failed: "+ err.Error()
			} else {
				for i := range(result) {
					reply = reply + result[i] + " "
				}
			}
		}

	case "rdns":

	case "ping":

	case "help":
		reply = "help yourself"

	}


	client.sendPrivateMessage(user,reply)
}


func (client *Client) sendRaw(message string) {
	client.NetOutgoing <- message
}

func (client *Client) sendMessage(dest,message string) {
	client.NetOutgoing <- "PRIVMSG "+dest+" :"+message
}

func (client *Client) sendChannelMessage(channel,message string) {
	client.sendMessage(channel,message)
	// log it like an incoming message now
	client.addChannelMessage(channel,client.Nick+": "+message);
}

func (client *Client) sendPrivateMessage(user,message string) {
	client.sendMessage(user,message)
	// log it like an incoming message now
	client.addPrivateMessageUser(user)
	v := client.PrivateMessages[user]
	v.Messages = append(client.PrivateMessages[user].Messages,client.Nick+": "+message)
	client.PrivateMessages[user] = v
}


