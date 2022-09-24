package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	ApplicationName string
	ListenAddress   string
	ListenPort      int
	SendHost        string
}

type Config struct {
	ClientHost         string
	ClientListenToPort int
	ClientSendToPort   int
	Routes             []Route
}

func (Config) Default() Config {
	config := Config{
		ClientHost:         "127.0.0.1",
		ClientListenToPort: 9000,
		ClientSendToPort:   9001,
		Routes: []Route{
			{
				ApplicationName: "ExampleApplication1",
				ListenAddress:   "127.0.0.1",
				ListenPort:      7001,
				SendHost:        "127.0.0.1:7000",
			},
			{
				ApplicationName: "ExampleApplication2",
				ListenAddress:   "127.0.0.1",
				ListenPort:      8001,
				SendHost:        "127.0.0.1:8000",
			},
		},
	}
	return config
}

func LoadConfig() Config {
	if _, err := os.Stat("config.json"); err == nil {
		// Read JSON data from file
		fileContent, readErr := os.ReadFile("config.json")
		if readErr != nil {
			fmt.Printf("Could not read config file: %s\n", readErr)
			return Config{}.Default()
		}
		result := Config{}.Default()
		json.Unmarshal(fileContent, &result)
		return result
	} else {
		// Create the file and use default
		jsondata, err := json.MarshalIndent(Config{}.Default(), "", "\t")
		if err != nil {
			fmt.Printf("Could not Marshal JSON: %s\n", err)
			return Config{}.Default()
		}
		if errWrite := os.WriteFile("config.json", jsondata, 0666); errWrite != nil {
			fmt.Printf("Could not write JSON to file: %s\n", errWrite)
		}
		return Config{}.Default()
	}
}

func SaveConfig(config Config) {
	jsondata, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Printf("Could not Marshal JSON: %s\n", err)
	}
	if errWrite := os.WriteFile("config.json", jsondata, 0666); errWrite != nil {
		fmt.Printf("Could not write JSON to file: %s\n", errWrite)
	}
}

func GetRouteIndexByName(applicationName string, config Config) int {
	for i := 0; i < len(config.Routes); i++ {
		if config.Routes[i].ApplicationName == applicationName {
			return i
		}
	}
	panic("Failed to find an Application fro mthe provided Config!")
}

func FilterRoutesByName(applicationName string, config Config) []Route {
	newRoutes := make([]Route, len(config.Routes)-1)
	for x := 0; x < len(newRoutes); x++ {
		for i := 0; i < len(config.Routes); i++ {
			if config.Routes[i].ApplicationName != applicationName {
				newRoutes[x] = config.Routes[i]
			}
		}
	}
	return newRoutes
}

// TODO: Why does this create indexes based on newRoutes? (fills in rest as empty after for loop)
func AddRouteToRoutes(route Route, routes []Route) []Route {
	newRoutes := make([]Route, len(routes)+1)
	for i := 0; i < len(routes); i++ {
		newRoutes[i] = routes[i]
	}
	newRoutes[len(newRoutes)-1] = route
	return newRoutes
}
