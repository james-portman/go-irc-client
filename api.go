package main

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ApiChannels struct{
	Client *Client
}
func (apiChannels ApiChannels) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	var channels []string
	for channelName,_ := range(apiChannels.Client.Channels) {
		channels = append(channels,channelName)
	}
	data := map[string][]string{"channels": channels}
	return 200, data, http.Header{"Content-type": {"application/json"}}

}



type ApiChannelUsers struct{
	Client *Client
}
func (apiChannelUsers ApiChannelUsers) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	channelName := values.Get("channel")
	if channelName == "" {
		data := map[string]string{"error": "bad request, missing channel in query string"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	users := make(map[string]string,0)
	for userName,channelUser := range apiChannelUsers.Client.Channels[channelName].Users {
		users[userName] = channelUser.Mode
	}
	data := map[string]map[string]string{"users": users}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}



type ApiChannelMessages struct{
	Client *Client
}
func (apiChannelMessages ApiChannelMessages) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	channelName := values.Get("channel")
	if channelName == "" {
		data := map[string]string{"error": "bad request, missing channel in query string"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	messages := make([]string,0)
	for _,message := range(apiChannelMessages.Client.Channels[channelName].Messages) {
		messages = append(messages,message)
	}

	data := map[string][]string{"messages": messages}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}




type ApiClients struct{
	Clients []*Client
}
func (apiClients ApiClients) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	var clients []map[string]string
	for _,client := range apiClients.Clients {
		x := make(map[string]string)
		x["host"] = client.ServerHost
		x["nick"] = client.Nick
		clients = append(clients,x)
	}
	
	data := map[string][]map[string]string {"clients": clients}
	return 200, data, http.Header{"Content-type": {"application/json"}}

}



type ApiServerMessages struct{
	Client *Client
}
func (apiServerMessages ApiServerMessages) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	data := map[string][]string {"messages": apiServerMessages.Client.ServerMessages}
	return 200, data, http.Header{"Content-type": {"application/json"}}

}

type ApiSendMessage struct{
	Client *Client
}

//post and get are the same for now

func (apiSendMessage ApiSendMessage) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	message := values.Get("message")
	if message == "" {
		data := map[string]string{"error": "bad request, missing message in request"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	user := values.Get("user")
	channel := values.Get("channel")

	if user != "" {
		apiSendMessage.Client.sendMessage(user,message)
	} else if channel != "" {
		if channel[0] != '#' && channel[0] != '&' && channel[0] != '+' && channel[0] != '~' {
			channel = "#"+channel
		}
		apiSendMessage.Client.sendMessage(channel,message)
	} else {
		data := map[string]string{"error": "bad request, must provide user or channel in request"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}
	

	data := map[string]string {"result": "ok"}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}

func (apiSendMessage ApiSendMessage) Post(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	message := values.Get("message")
	if message == "" {
		data := map[string]string{"error": "bad request, missing message in request"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	user := values.Get("user")
	channel := values.Get("channel")

	if user != "" {
		apiSendMessage.Client.sendPrivateMessage(user,message)
	} else if channel != "" {
		if channel[0] != '#' && channel[0] != '&' && channel[0] != '+' && channel[0] != '~' {
			channel = "#"+channel
		}
		apiSendMessage.Client.sendChannelMessage(channel,message)
		
	} else {
		data := map[string]string{"error": "bad request, must provide user or channel in request"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}
	

	data := map[string]string {"result": "ok"}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}




type ApiSendRaw struct{
	Client *Client
}
func (apiSendRaw ApiSendRaw) Post(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	message := values.Get("message")
	if message == "" {
		data := map[string]string{"error": "bad request, missing message in request"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	apiSendRaw.Client.sendRaw(message)

	data := map[string]string {"result": "ok"}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}



/*
	if queryIp == "" {
		data := map[string]string{"error": "bad request, missing ip in query string"}
		return 400, data, http.Header{"Content-type": {"application/json"}}
	}

	if result := includes.GetRdns(queryIp); result == "" {
		data := map[string]string{"error": "failed to get reverse DNS for" + queryIp}
		return 500, data, http.Header{"Content-type": {"application/json"}}

	}
*/


type ApiWebClient struct{
}
func (apiWebClient ApiWebClient) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	content, err := ioutil.ReadFile("webclient.html")
	if err != nil {
		return 500, "error: couldnt get web client content", http.Header{"Content-type": {"text/html"}}
	} else {
		return 200, content, http.Header{"Content-type": {"text/html"}}
	}
}

type ApiWebClientIcon struct{
}
func (apiWebClientIcon ApiWebClientIcon) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {

	content, err := ioutil.ReadFile("irc.png")
	if err != nil {
		return 500, "error: couldnt get web client content", http.Header{"Content-type": {"text/html"}}
	} else {
		return 200, content, http.Header{"Content-type": {"image/png"}}
	}
}