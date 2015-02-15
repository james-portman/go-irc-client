package main

import (
	"strings"
	"fmt"
)

func (client *Client) incomingParser() {

	// for, select, case stops it smashing cpu
    for {
    select {
    case data := <-client.NetIncoming:

    	// split message into command etc
        words := strings.Fields(data)

        // check for prefix (who)
        // prefix =  servername / ( nickname [ [ "!" user ] "@" host ] )
        var prefix string
        if words[0][0] == ':' {
            prefix = words[0]
            prefix = prefix[1:]

            if client.ServerPrefix == prefix || client.ServerHost == prefix {
            	// is the server
            } else { // might be user so split the nick out
            	// cut off @~~~
            	at := strings.Index(prefix,"@")
            	if at != -1 {
            		prefix = prefix[0:(at)]
            		// cut off !~~~
            		excl := strings.Index(prefix,"!")
            		if excl != -1 {
            			prefix = prefix[0:(excl)]
            		}
            	}
            }
            fmt.Println("prefix: "+prefix+" ")

            // remove prefix from the rest:
            data = removeFirstWord(data)
            words = strings.Fields(data)
        }

        command := words[0]
        stringparams := removeFirstWord(data)
        params := strings.Fields(stringparams)

        if params[0] == "lol" {
        	// dont moan
        }

        switch command {

        default:
            client.addServerMessage(data)


        // leave anything commented to just log to ServerMessages

    	case "001": // RPL_WELCOME - parse nick out?
    		client.ServerPrefix = prefix
            client.addServerMessage(data)
    		Log("saved server prefix from this line as "+prefix)
        // case "002": // RPL_YOURHOST
        // case "003": // RPL_CREATED
        // case "004": // RPL_MYINFO
        // case "005": // RPL_BOUNCE

        //302    RPL_USERHOST
        //303    RPL_ISON
        // case "301": // RPL_AWAY "<nick> :<away message>"
       // 305    RPL_UNAWAY ":You are no longer marked as being away"
       // 306    RPL_NOWAWAY ":You have been marked as being away"
       // 311    RPL_WHOISUSER "<nick> <user> <host> * :<real name>"
       // 312    RPL_WHOISSERVER "<nick> <server> :<server info>"
       // 313    RPL_WHOISOPERATOR "<nick> :is an IRC operator"
       // 317    RPL_WHOISIDLE "<nick> <integer> :seconds idle"
       // 318    RPL_ENDOFWHOIS "<nick> :End of WHOIS list"
       // 319    RPL_WHOISCHANNELS "<nick> :*( ( "@" / "+" ) <channel> " " )"
       // 314    RPL_WHOWASUSER "<nick> <user> <host> * :<real name>"
       // 369    RPL_ENDOFWHOWAS "<nick> :End of WHOWAS"
       // 321    RPL_LISTSTART Obsolete. Not used.
       // 322    RPL_LIST "<channel> <# visible> :<topic>"
       // 323    RPL_LISTEND ":End of LIST"
       // 325    RPL_UNIQOPIS "<channel> <nickname>"
       // 324    RPL_CHANNELMODEIS "<channel> <mode> <mode params>"
       // 331    RPL_NOTOPIC "<channel> :No topic is set"
        case "332": //    RPL_TOPIC "<channel> :<topic>"
            channel := params[1]
            topic := strings.Join(params[2:]," ")
            client.gotChannelTopic(channel,topic)
            
       // 341    RPL_INVITING "<channel> <nick>"
       // 342    RPL_SUMMONING "<user> :Summoning user to IRC"
       // 346    RPL_INVITELIST "<channel> <invitemask>"
       // 347    RPL_ENDOFINVITELIST "<channel> :End of channel invite list"
       // 348    RPL_EXCEPTLIST "<channel> <exceptionmask>"
       // 349    RPL_ENDOFEXCEPTLIST "<channel> :End of channel exception list"
       // 351    RPL_VERSION "<version>.<debuglevel> <server> :<comments>"
       // 352    RPL_WHOREPLY
       //        "<channel> <user> <host> <server> <nick>
       //        ( "H" / "G" > ["*"] [ ( "@" / "+" ) ]
       //        :<hopcount> <real name>"
       // 315    RPL_ENDOFWHO "<name> :End of WHO list"

       case "353"://    RPL_NAMREPLY
            // params[0] will be my nick
            // params[1] will be @/*/= channel type
            // params[2] will be channel
            if params[3][0] == ':' {
                params[3] = params[3][1:]
            }
            for i := range(params) {
                if i < 3 { continue }
                mode := ""
                if params[i][0] == '@' || params[i][0] == '+' {
                    mode = string(params[i][0])
                    params[i] = params[i][1:]
                }
                client.addUserToChannel(params[2],params[i],mode)
            }

       // 366    RPL_ENDOFNAMES "<channel> :End of NAMES list"
       // 364    RPL_LINKS "<mask> <server> :<hopcount> <server info>"
       // 365    RPL_ENDOFLINKS "<mask> :End of LINKS list"
       // 367    RPL_BANLIST "<channel> <banmask>"
       // 368    RPL_ENDOFBANLIST "<channel> :End of channel ban list"
       // 371    RPL_INFO ":<string>"
       // 374    RPL_ENDOFINFO ":End of INFO list"
       // 375    RPL_MOTDSTART ":- <server> Message of the day - "
       // 372    RPL_MOTD ":- <text>"
       // 376    RPL_ENDOFMOTD ":End of MOTD command"
       // 381    RPL_YOUREOPER ":You are now an IRC operator"
       // 382    RPL_REHASHING "<config file> :Rehashing"
       // 383    RPL_YOURESERVICE "You are service <servicename>"
       // 391    RPL_TIME "<server> :<string showing server's local time>"
       // 392    RPL_USERSSTART ":UserID   Terminal  Host"
       // 393    RPL_USERS ":<username> <ttyline> <hostname>"
       // 394    RPL_ENDOFUSERS ":End of users"
       // 395    RPL_NOUSERS ":Nobody logged in"
       //  200    RPL_TRACELINK
       //        "Link <version & debug level> <destination>
       //         <next server> V<protocol version>
       //         <link uptime in seconds> <backstream sendq>
       //         <upstream sendq>"
       // 201    RPL_TRACECONNECTING "Try. <class> <server>"
       // 202    RPL_TRACEHANDSHAKE "H.S. <class> <server>"
       // 203    RPL_TRACEUNKNOWN "???? <class> [<client IP address in dot form>]"
       // 204    RPL_TRACEOPERATOR "Oper <class> <nick>"
       // 205    RPL_TRACEUSER "User <class> <nick>"
       // 206    RPL_TRACESERVER
       //        "Serv <class> <int>S <int>C <server>
       //         <nick!user|*!*>@<host|server> V<protocol version>"
       // 207    RPL_TRACESERVICE "Service <class> <name> <type> <active type>"
       // 208    RPL_TRACENEWTYPE "<newtype> 0 <client name>"
       // 209    RPL_TRACECLASS "Class <class> <count>"
       // 210    RPL_TRACERECONNECT Unused.
       // 261    RPL_TRACELOG "File <logfile> <debug level>"
       // 262    RPL_TRACEEND "<server name> <version & debug level> :End of TRACE"
       // 211    RPL_STATSLINKINFO
       //        "<linkname> <sendq> <sent messages>
       //         <sent Kbytes> <received messages>
       //         <received Kbytes> <time open>"
       // 212    RPL_STATSCOMMANDS "<command> <count> <byte count> <remote count>"
       // 219    RPL_ENDOFSTATS "<stats letter> :End of STATS report"
       // 242    RPL_STATSUPTIME ":Server Up %d days %d:%02d:%02d"
       // 243    RPL_STATSOLINE "O <hostmask> * <name>"
       // 221    RPL_UMODEIS "<user mode string>"
       // 234    RPL_SERVLIST "<name> <server> <mask> <type> <hopcount> <info>"
       // 235    RPL_SERVLISTEND "<mask> <type> :End of service listing"
       // 251    RPL_LUSERCLIENT
       //        ":There are <integer> users and <integer>
       //         services on <integer> servers"
       // 252    RPL_LUSEROP "<integer> :operator(s) online"
       // 253    RPL_LUSERUNKNOWN "<integer> :unknown connection(s)"
       // 254    RPL_LUSERCHANNELS "<integer> :channels formed"
       // 255    RPL_LUSERME ":I have <integer> clients and <integer> servers"
       // 256    RPL_ADMINME "<server> :Administrative info"
       // 257    RPL_ADMINLOC1 ":<admin info>"
       // 258    RPL_ADMINLOC2 ":<admin info>"
       // 259    RPL_ADMINEMAIL ":<admin info>"
       // 263    RPL_TRYAGAIN "<command> :Please wait a while and try again."


// 5.2 Error Replies
       // 401    ERR_NOSUCHNICK "<nickname> :No such nick/channel"
       // 402    ERR_NOSUCHSERVER "<server name> :No such server"
       // 403    ERR_NOSUCHCHANNEL "<channel name> :No such channel"
       // 404    ERR_CANNOTSENDTOCHAN "<channel name> :Cannot send to channel"
       // 405    ERR_TOOMANYCHANNELS "<channel name> :You have joined too many channels"
       // 406    ERR_WASNOSUCHNICK "<nickname> :There was no such nickname"
       // 407    ERR_TOOMANYTARGETS "<target> :<error code> recipients. <abort message>"
       // 408    ERR_NOSUCHSERVICE "<service name> :No such service"
       // 409    ERR_NOORIGIN ":No origin specified"
       // 411    ERR_NORECIPIENT ":No recipient given (<command>)"
       // 412    ERR_NOTEXTTOSEND ":No text to send"
       // 413    ERR_NOTOPLEVEL "<mask> :No toplevel domain specified"
       // 414    ERR_WILDTOPLEVEL "<mask> :Wildcard in toplevel domain"
       // 415    ERR_BADMASK "<mask> :Bad Server/host mask"
       // 421    ERR_UNKNOWNCOMMAND "<command> :Unknown command"
       // 422    ERR_NOMOTD ":MOTD File is missing"
       // 423    ERR_NOADMININFO "<server> :No administrative info available"
       // 424    ERR_FILEERROR ":File error doing <file op> on <file>"
       // 431    ERR_NONICKNAMEGIVEN ":No nickname given"
       // 432    ERR_ERRONEUSNICKNAME "<nick> :Erroneous nickname"
       // 433    ERR_NICKNAMEINUSE "<nick> :Nickname is already in use"
       // 436    ERR_NICKCOLLISION "<nick> :Nickname collision KILL from <user>@<host>"
       // 437    ERR_UNAVAILRESOURCE "<nick/channel> :Nick/channel is temporarily unavailable"
       // 441    ERR_USERNOTINCHANNEL "<nick> <channel> :They aren't on that channel"
       // 442    ERR_NOTONCHANNEL "<channel> :You're not on that channel"
       // 443    ERR_USERONCHANNEL "<user> <channel> :is already on channel"
       // 444    ERR_NOLOGIN "<user> :User not logged in"
       // 445    ERR_SUMMONDISABLED ":SUMMON has been disabled"
       // 446    ERR_USERSDISABLED":USERS has been disabled"
       // 451    ERR_NOTREGISTERED ":You have not registered"
       // 461    ERR_NEEDMOREPARAMS "<command> :Not enough parameters"
       // 462    ERR_ALREADYREGISTRED ":Unauthorized command (already registered)"
       // 463    ERR_NOPERMFORHOST ":Your host isn't among the privileged"
       // 464    ERR_PASSWDMISMATCH ":Password incorrect"
       // 465    ERR_YOUREBANNEDCREEP ":You are banned from this server"
       // 466    ERR_YOUWILLBEBANNED
       // 467    ERR_KEYSET "<channel> :Channel key already set"
       // 471    ERR_CHANNELISFULL "<channel> :Cannot join channel (+l)"
       // 472    ERR_UNKNOWNMODE "<char> :is unknown mode char to me for <channel>"
       // 473    ERR_INVITEONLYCHAN "<channel> :Cannot join channel (+i)"
       // 474    ERR_BANNEDFROMCHAN "<channel> :Cannot join channel (+b)"
       // 475    ERR_BADCHANNELKEY "<channel> :Cannot join channel (+k)"
       // 476    ERR_BADCHANMASK "<channel> :Bad Channel Mask"
       // 477    ERR_NOCHANMODES "<channel> :Channel doesn't support modes"
       // 478    ERR_BANLISTFULL "<channel> <char> :Channel list is full"
       // 481    ERR_NOPRIVILEGES ":Permission Denied- You're not an IRC operator"
       // 482    ERR_CHANOPRIVSNEEDED "<channel> :You're not channel operator"
       // 483    ERR_CANTKILLSERVER ":You can't kill a server!"
       // 484    ERR_RESTRICTED ":Your connection is restricted!"
       // 485    ERR_UNIQOPPRIVSNEEDED ":You're not the original channel operator"
       // 491    ERR_NOOPERHOST ":No O-lines for your host"
       // 501    ERR_UMODEUNKNOWNFLAG ":Unknown MODE flag"
       // 502    ERR_USERSDONTMATCH ":Cannot change mode for other users"



    	case "JOIN":
    		params[0] = params[0][1:]
            client.addChannelMessage(params[0],"* "+prefix+" joined")
    		client.addUserToChannel(params[0],prefix,"")

        // case "NOTICE":
        //     if pararms[0] == client.Nick {
        //             // private
        //         } else {
        //             // channel
        //         }

        case "PART":
          if prefix == client.Nick {
            client.removeChannel(params[0])
          } else {
            client.addChannelMessage(params[0],"* "+prefix+" left")
        	client.removeUserFromChannel(params[0],prefix)
          }

        case "PING":
            //Log("PING received, replying")
            client.NetOutgoing <- "PONG "+words[1]


    	case "PRIVMSG": // chan or priv message
            message := strings.Join(params," ")
            message = removeFirstWord(message)
            if message[0] == ':' {
                message = message[1:]
            }

            if params[0] == client.Nick {
                client.privateMessageReceived(prefix,message)

                // private message so presume a command
                client.parseCommandMessage(message,"",prefix)
                

    		} else {
                client.addChannelMessage(params[0],prefix+": "+message)

                if params[1] == ":"+client.Nick || params[1] == ":"+client.Nick+":" {
                    client.parseCommandMessage(message,params[0],prefix)
                }
    		}




        }



    }} // closing for,select,case
}

