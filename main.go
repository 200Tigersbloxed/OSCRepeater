package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

type Router struct {
	Route  Route
	Client *osc.Client
	Server *osc.Server
}

var (
	oscClient *osc.Client = nil
	oscServer *osc.Server = nil
	oscRoutes []Router
	oscConfig = Config{}.Default()
	didStart  bool
)

const (
	version = "pre0.1"
)

func main() {
	oscConfig = LoadConfig()
	//StartOSCClient()
	//StartOSCRoutes(oscConfig.Routes)
	window()
}

func StartOSCClient() {
	oscClient = osc.NewClient(oscConfig.ClientHost, oscConfig.ClientSendToPort)
	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/", func(msg *osc.Message) {
		for i := 0; i < len(oscRoutes); i++ {
			currentRouter := oscRoutes[i]
			message := osc.NewMessage(msg.Address)
			for x := 0; x < len(msg.Arguments); x++ {
				message.Append(msg.Arguments[x])
			}
			currentRouter.Client.Send(message)
		}
	})
	fmt.Println(rune(oscConfig.ClientListenToPort))
	oscServer = &osc.Server{
		Addr:       oscConfig.ClientHost + ":" + string(rune(oscConfig.ClientListenToPort)),
		Dispatcher: d,
	}
	go oscServer.ListenAndServe()
}

func StartOSCRoutes(JSONRoutes []Route) {
	oscRoutes = make([]Router, len(JSONRoutes))
	for i := 0; i < len(JSONRoutes); i++ {
		currentRoute := JSONRoutes[i]
		client := osc.NewClient(currentRoute.ListenAddress, currentRoute.ListenPort)
		d := osc.NewStandardDispatcher()
		d.AddMsgHandler("/", func(msg *osc.Message) {
			message := osc.NewMessage(msg.Address)
			for i := 0; i < len(msg.Arguments); i++ {
				message.Append(msg.Arguments[i])
			}
			oscClient.Send(message)
		})
		server := &osc.Server{
			Addr:       currentRoute.SendHost,
			Dispatcher: d,
		}
		go server.ListenAndServe()
		router := Router{
			Route:  currentRoute,
			Client: client,
			Server: server,
		}
		oscRoutes[i] = router
		add := "Added Route " + currentRoute.ApplicationName + " Listening on " + currentRoute.ListenAddress + ":"
		fmt.Printf("\n")
		fmt.Print(add)
		fmt.Print(currentRoute.ListenPort)
		fmt.Printf("\n")
	}
}

func StopAll() {
	for i := 0; i < len(oscRoutes); i++ {
		currentRouter := oscRoutes[i]
		currentRouter.Server.CloseConnection()
	}
	oscClient = nil
	oscServer.CloseConnection()
	oscServer = nil
	oscRoutes = []Router{}
	didStart = false
	updateLabel.SetText("Status: " + getTextFromStatus(didStart))
}
