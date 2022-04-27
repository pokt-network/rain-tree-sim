package main

import "log"

// TODO: Discuss how most things are being passed by value than by reference.

func main() {
	config := LoadConfigFile()

	log.Printf("About to start a RainTree network with the following configuration: %s\n", config.String())

	NewRainTreeNetwork(config)
}
