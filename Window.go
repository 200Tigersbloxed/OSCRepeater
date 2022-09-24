package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var (
	a           fyne.App
	c           *fyne.Container
	updateLabel *widget.Label
)

func serverSettingsWindow(mainWindow fyne.Window) {
	window := a.NewWindow("Server Settings")
	clientHost := widget.NewEntry()
	clientHost.SetText(oscConfig.ClientHost)
	clientHost.SetPlaceHolder(oscConfig.ClientHost)
	clientListenToPort := widget.NewEntry()
	clientListenToPort.SetText(strconv.Itoa(oscConfig.ClientListenToPort))
	clientListenToPort.SetPlaceHolder(strconv.Itoa(oscConfig.ClientListenToPort))
	clientSendToPort := widget.NewEntry()
	clientSendToPort.SetText(strconv.Itoa(oscConfig.ClientSendToPort))
	clientSendToPort.SetPlaceHolder(strconv.Itoa(oscConfig.ClientSendToPort))
	window.SetContent(container.NewVBox(
		widget.NewLabel("ClientHost"),
		widget.NewLabel("Defines the host to repeat all OSC messages to (127.0.0.1)"),
		clientHost,
		widget.NewLabel("ClientListenToPort"),
		widget.NewLabel("Defines the port to listen for OSC messages coming from the server endpoint (9001)"),
		clientListenToPort,
		widget.NewLabel("ClientSendToPort"),
		widget.NewLabel("Defines the port to send repeated OSC messages to (9000)"),
		clientSendToPort,
		widget.NewButton("APPLY", func() {
			didChangeOne := false
			if len(clientHost.Text) > 0 {
				oscConfig.ClientHost = clientHost.Text
				didChangeOne = true
			}
			if len(clientListenToPort.Text) > 0 {
				oscConfig.ClientListenToPort, _ = strconv.Atoi(clientListenToPort.Text)
				didChangeOne = true
			}
			if len(clientSendToPort.Text) > 0 {
				oscConfig.ClientSendToPort, _ = strconv.Atoi(clientSendToPort.Text)
				didChangeOne = true
			}
			if didChangeOne {
				SaveConfig(oscConfig)
				reloadConfigFromWindow(mainWindow)
				window.Close()
			}
		}),
	))
	window.Resize(fyne.NewSize(400, 0))
	window.Show()
}

func getButtonWidgets(mainWindow fyne.Window) []*widget.Button {
	buttons := make([]*widget.Button, len(oscConfig.Routes))
	for i := 0; i < len(oscConfig.Routes); i++ {
		route := oscConfig.Routes[i]
		buttons[i] = widget.NewButton(route.ApplicationName, func() {
			window := a.NewWindow(route.ApplicationName)
			appwindow := AppWindow{}.CreateAndFill(route)
			window.SetContent(container.NewVBox(
				widget.NewLabel("ApplicationName"),
				widget.NewLabel("The name of your application"),
				appwindow.applicationName,
				widget.NewLabel("ListenAddress"),
				widget.NewLabel("The Address to Listen for OSC Messages (127.0.0.1)"),
				appwindow.listenAddress,
				widget.NewLabel("ListenPort"),
				widget.NewLabel("The Port to Listen for OSC Messages to be repeated to the SendHost"),
				appwindow.listenPort,
				widget.NewLabel("SendHost"),
				widget.NewLabel("The Host to open for Repeating"),
				appwindow.sendHost,
				widget.NewButton("SAVE", func() {
					if len(appwindow.applicationName.Text) > 0 {
						oscConfig.Routes[GetRouteIndexByName(appwindow.applicationName.PlaceHolder, oscConfig)].ApplicationName = appwindow.applicationName.Text
					}
					if len(appwindow.listenAddress.Text) > 0 {
						oscConfig.Routes[GetRouteIndexByName(appwindow.applicationName.PlaceHolder, oscConfig)].ListenAddress = appwindow.listenAddress.Text
					}
					if len(appwindow.listenPort.Text) > 0 {
						oscConfig.Routes[GetRouteIndexByName(appwindow.applicationName.PlaceHolder, oscConfig)].ListenPort, _ = strconv.Atoi(appwindow.listenPort.Text)
					}
					if len(appwindow.sendHost.Text) > 0 {
						oscConfig.Routes[GetRouteIndexByName(appwindow.applicationName.PlaceHolder, oscConfig)].SendHost = appwindow.sendHost.Text
					}
					SaveConfig(oscConfig)
					reloadConfigFromWindow(mainWindow)
					window.Close()
				}),
				widget.NewButton("DELETE", func() {
					oscConfig.Routes = FilterRoutesByName(appwindow.applicationName.PlaceHolder, oscConfig)
					SaveConfig(oscConfig)
					reloadConfigFromWindow(mainWindow)
					window.Close()
				}),
			))
			window.Resize(fyne.NewSize(400, 0))
			window.Show()
		})
	}
	return buttons
}

func reloadConfigFromWindow(window fyne.Window) {
	if didStart {
		StopAll()
	}
	oscConfig = LoadConfig()
	clearAndRegenContainer(window)
}

func getTextFromStatus(status bool) string {
	if status == true {
		return "Running"
	}
	return "Stopped"
}

func clearAndRegenContainer(w fyne.Window) {
	if c != nil {
		c.RemoveAll()
	}
	c = container.NewVBox(
		widget.NewLabel("OSCRouter - "+version),
		widget.NewSeparator(),
	)
	c.Add(widget.NewButton("Server Settings", func() {
		serverSettingsWindow(w)
	}))
	c.Add(widget.NewSeparator())
	c.Add(widget.NewLabel("Applications"))
	appButtons := getButtonWidgets(w)
	if len(appButtons) <= 0 {
		c.Add(widget.NewLabel("I looked far and wide, alas, no Applications :("))
	} else {
		for i := 0; i < len(appButtons); i++ {
			c.Add(appButtons[i])
		}
	}
	c.Add(widget.NewSeparator())
	c.Add(widget.NewLabel("Actions"))
	c.Add(widget.NewButton("Start", func() {
		StartOSCClient()
		StartOSCRoutes(oscConfig.Routes)
		didStart = true
		updateLabel.SetText("Status: " + getTextFromStatus(didStart))
	}))
	c.Add(widget.NewButton("Stop", func() {
		if didStart {
			StopAll()
		}
	}))
	c.Add(widget.NewButton("Create Application", func() {
		window := a.NewWindow("New Application Wizard")
		appwindow := AppWindow{}.CreateTemplate()
		window.SetContent(container.NewVBox(
			widget.NewLabel("ApplicationName"),
			widget.NewLabel("The name of your application"),
			appwindow.applicationName,
			widget.NewLabel("ListenAddress"),
			widget.NewLabel("The Address to Listen for OSC Messages (127.0.0.1)"),
			appwindow.listenAddress,
			widget.NewLabel("ListenPort"),
			widget.NewLabel("The Port to Listen for OSC Messages to be repeated to the SendHost"),
			appwindow.listenPort,
			widget.NewLabel("SendHost"),
			widget.NewLabel("The Host to open for Repeating"),
			appwindow.sendHost,
			widget.NewButton("SAVE", func() {
				if len(appwindow.applicationName.Text) > 0 && len(appwindow.listenAddress.Text) > 0 &&
					len(appwindow.listenPort.Text) > 0 && len(appwindow.sendHost.Text) > 0 {
					port, _ := strconv.Atoi(appwindow.listenPort.Text)
					route := Route{
						ApplicationName: appwindow.applicationName.Text,
						ListenAddress:   appwindow.listenAddress.Text,
						ListenPort:      port,
						SendHost:        appwindow.sendHost.Text,
					}
					oscConfig.Routes = AddRouteToRoutes(route, oscConfig.Routes)
					SaveConfig(oscConfig)
					reloadConfigFromWindow(w)
					window.Close()
				}
			}),
		))
		window.Resize(fyne.NewSize(200, 0))
		window.Show()
	}))
	c.Add(widget.NewButton("Reload Config", func() {
		reloadConfigFromWindow(w)
	}))
	c.Add(widget.NewSeparator())
	updateLabel = widget.NewLabel("Status: " + getTextFromStatus(didStart))
	c.Add(updateLabel)
	w.SetContent(c)
}

func window() {
	a = app.New()
	w := a.NewWindow("OSCRepeater")
	clearAndRegenContainer(w)
	w.Resize(fyne.NewSize(400, 0))
	w.ShowAndRun()
}

type AppWindow struct {
	route           Route
	applicationName *widget.Entry
	listenAddress   *widget.Entry
	listenPort      *widget.Entry
	sendHost        *widget.Entry
}

func (AppWindow) CreateAndFill(route Route) AppWindow {
	appWindow := AppWindow{
		route:           route,
		applicationName: widget.NewEntry(),
		listenAddress:   widget.NewEntry(),
		listenPort:      widget.NewEntry(),
		sendHost:        widget.NewEntry(),
	}
	appWindow.applicationName.SetPlaceHolder(route.ApplicationName)
	appWindow.listenAddress.SetPlaceHolder(route.ListenAddress)
	appWindow.listenPort.SetPlaceHolder(strconv.Itoa(route.ListenPort))
	appWindow.sendHost.SetPlaceHolder(route.SendHost)
	return appWindow
}

func (AppWindow) CreateTemplate() AppWindow {
	appWindow := AppWindow{
		route:           Route{},
		applicationName: widget.NewEntry(),
		listenAddress:   widget.NewEntry(),
		listenPort:      widget.NewEntry(),
		sendHost:        widget.NewEntry(),
	}
	appWindow.applicationName.SetPlaceHolder("ApplicationName")
	appWindow.listenAddress.SetPlaceHolder("ListenAddress")
	appWindow.listenPort.SetPlaceHolder("ListenPort")
	appWindow.sendHost.SetPlaceHolder("SendHost")
	return appWindow
}
