package main

// http://tools.ietf.org/html/rfc2812

// todo:
//
// save links into a db/sqlite at least for history
// optional save all history? sqlite?
// away status

import (
        "bufio"
        "io"
        "fmt"
        "net"
        "time"
        "crypto/tls"
        // "os"
        "strings"
)


func Log(v ...interface{}) {
        fmt.Println(v...)
}

func (client *Client) outgoingSender() {
        Log("outgoingSender started")
        for {
                data := <- client.NetOutgoing
                Log("sending "+data)
                fmt.Fprintf(client.Conn, data+"\r\n")
        }
        Log("outgoingSender stopping")
}

func (client *Client) incomingReceiver() {
        Log("incomingReceiver started")

        reader := bufio.NewReader(client.Conn)
        for {
                line, err := reader.ReadString('\n')
                if err != nil {
                        if err == io.EOF {
                                Log("EOF received")
                        } else {
                                Log(err.Error())
                        }
                        break
                } else {
                        client.NetIncoming <- line
                }
        }

        Log("incomingReceiver stopping")
}




func removeFirstWord(in string) string {
        words := strings.Fields(in)
        in = in[len(words[0])+1:]
        return in        
}


func (client *Client) newClient() error {

        var err error
        if client.Secure {
                client.Conn, err = tls.Dial("tcp", client.ServerHost+":"+client.ServerPort, &tls.Config{InsecureSkipVerify: true})
        } else {
                client.Conn, err = net.Dial("tcp", client.ServerHost+":"+client.ServerPort)
        }
        if err != nil {
                return err
        } else {
                Log("connected OK")
                time.Sleep(1 * time.Second)
                if client.ServerPass != "" {
                        client.NetOutgoing <- "PASS "+client.ServerPass
                }
                client.NetOutgoing <- "NICK "+client.Nick
                client.NetOutgoing <- "User username 8 * :full name"

                go client.incomingReceiver()
                go client.outgoingSender()
                go client.incomingParser()
                return nil
        }
}





func main() {

		// always do these bits no matter how many connections get added
        api := NewAPI()

        apiWebClient := new(ApiWebClient)
        api.AddResource(apiWebClient, "/webclient")

        apiWebClientIcon := new(ApiWebClientIcon)
        api.AddResource(apiWebClientIcon, "/irc.png")

        apiClients := new(ApiClients)
        api.AddResource(apiClients, "/clients")


        ukfastclient := &Client{"nickname",true,"irc.server.org","6667","isthereaserverpassword?","",make(chan string, 32), make(chan string, 32), nil, make(chan bool), make([]string,0), make(map[string]Channel,0), make(map[string]PrivateMessage,0) }
        setupClient(ukfastclient,api,apiClients);


        for {
        	time.Sleep(1 * time.Second);
        }

}


func setupClient(client *Client, api *API, apiClients *ApiClients) {

        apiClients.Clients = append(apiClients.Clients,client)


        apiChannels := new(ApiChannels)
        apiChannels.Client = client
        api.AddResource(apiChannels, "/"+client.ServerHost+"/channels")

        apiChannelUsers := new(ApiChannelUsers)
        apiChannelUsers.Client = client
        api.AddResource(apiChannelUsers, "/"+client.ServerHost+"/channel/users")

        apiChannelMessages := new(ApiChannelMessages)
        apiChannelMessages.Client = client
        api.AddResource(apiChannelMessages, "/"+client.ServerHost+"/channel/messages")

        apiServerMessages := new(ApiServerMessages)
        apiServerMessages.Client = client
        api.AddResource(apiServerMessages, "/"+client.ServerHost+"/messages")

        apiSendRaw := new(ApiSendRaw)
        apiSendRaw.Client = client
        api.AddResource(apiSendRaw, "/"+client.ServerHost)


        apiSendMessage := new(ApiSendMessage)
        apiSendMessage.Client = client
        api.AddResource(apiSendMessage, "/"+client.ServerHost+"/sendmessage")


        // do a link for sendmessage with out the server address too?
        // for git notifications etc?


        err := client.newClient();
        if err != nil {
                Log("error "+err.Error())
        }
}
