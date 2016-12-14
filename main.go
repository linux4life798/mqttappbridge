// Author: Craig Hesling <craig@hesling.com>

// This is a simple program that copies mqtt messages
// from source topics into destination topics
package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"

	"flag"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttclientidprefix = "appbridge"
)

var mqttserver string
var mqttuser string
var mqttpass string

// generate a random client id for mqtt
func genclientid() string {
	r, err := rand.Int(rand.Reader, new(big.Int).SetInt64(100000))
	if err != nil {
		log.Fatal("Couldn't generate a random number for MQTT client ID")
	}
	return mqttclientidprefix + r.String()
}

func init() {
	flag.StringVar(&mqttserver, "mqtt_server", "localhost", "Sets the MQTT server")
	flag.StringVar(&mqttuser, "mqtt_user", "", "Sets the MQTT username")
	flag.StringVar(&mqttpass, "mqtt_pass", "", "Sets the MQTT password")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <source_topic1> <destination_topic1> [<source_topic2> <destination_topic2> [...]]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\noption:\n")
		flag.PrintDefaults()
	}
}

// test it
func main() {
	// parse arguments
	flag.Parse()

	// get extra non-flag arguments
	topics := flag.Args()

	/* Setup basic MQTT connection */
	opts := MQTT.NewClientOptions().AddBroker(mqttserver)
	opts.SetClientID(genclientid())
	// if username/password given
	if mqttuser != "" {
		opts.SetUsername(mqttuser)
		opts.SetPassword(mqttpass)
	}

	/* Create and start a client using the above ClientOptions */
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Failed to connect: ", token.Error())
	}
	defer c.Disconnect(250)
	log.Println("Connection successful")

	/* Register a handler for each source topic */
	for i := 0; i < len(topics)-1; i += 2 {
		// make sure each source has a destinations
		if i+1 >= len(topics) {
			log.Fatal("Unmatched source topic " + topics[i])
		}
		src := topics[i]
		dst := topics[i+1]
		log.Println("Map: " + src + " --> " + dst)

		// the most simple way of registering handlers
		c.Subscribe(src, 2, func(client MQTT.Client, message MQTT.Message) {
			client.Publish(dst, message.Qos(), message.Retained(), message.Payload())
		})

	}

	/* Wait for SIGINT */
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	<-signals

	log.Println("Shutting down")
}
