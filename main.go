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

	"strings"

	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"
)

const (
	mqttdefaultserver  = "tcp://localhost:1883"
	mqttclientidprefix = "appbridge"
)

const (
	jsonpathseparator = ",,"
)

/* Options to be filled in by arguments */
var mqttserver string
var mqttuser string
var mqttpass string
var simplesingle bool

/* Generate a random client id for mqtt */
func genclientid() string {
	r, err := rand.Int(rand.Reader, new(big.Int).SetInt64(100000))
	if err != nil {
		log.Fatal("Couldn't generate a random number for MQTT client ID")
	}
	return mqttclientidprefix + r.String()
}

/* Setup argument flags and help prompt */
func init() {
	flag.StringVar(&mqttserver, "mqtt_server", mqttdefaultserver, "Sets the MQTT server")
	flag.StringVar(&mqttuser, "mqtt_user", "", "Sets the MQTT username")
	flag.StringVar(&mqttpass, "mqtt_pass", "", "Sets the MQTT password")
	flag.BoolVar(&simplesingle, "simple_single", false, "Removes JSON array brackets from single values")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <source_topic1> <destination_topic1> [<source_topic2> <destination_topic2> [...]]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\noption:\n")
		flag.PrintDefaults()
	}
}

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
		var parts []string
		var src, dst string
		// make sure each source has a destinations
		if i+1 >= len(topics) {
			log.Fatal("Unmatched source topic " + topics[i])
		}

		parts = strings.Split(topics[i], jsonpathseparator)
		src = parts[0]
		dst = topics[i+1]
		if len(parts) > 1 {
			jpaths := parts[1:]
			log.Println("Map: " + src + " --> " + dst + " | JSON Path: " + fmt.Sprint(jpaths))

			// the most simple way of registering handlers
			c.Subscribe(src, 2, func(client MQTT.Client, message MQTT.Message) {
				results := gjson.GetMany(string(message.Payload()), jpaths...)
				// we want output message to be JSON with proper type representation
				values := make([]interface{}, len(results))
				for i, v := range results {
					values[i] = v.Value()
				}

				var output []byte
				if simplesingle && len(values) == 1 {
					output = []byte(fmt.Sprint(values[0]))
				} else {
					output, _ = json.Marshal(&values)
				}

				client.Publish(dst, message.Qos(), message.Retained(), output)
			})
		} else {
			log.Println("Map: " + src + " --> " + dst)

			// the most simple way of registering handlers
			c.Subscribe(src, 2, func(client MQTT.Client, message MQTT.Message) {
				client.Publish(dst, message.Qos(), message.Retained(), message.Payload())
			})
		}

	}

	/* Wait for SIGINT */
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	<-signals
	log.Println("Shutting down")
}
