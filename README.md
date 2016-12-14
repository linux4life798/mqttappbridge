# Description

This is a simple program that maps a set of mqtt source topics to a set
of destination topics.
This is done at the mqtt application/client level. It subscribes to all
source topics and publishes any messages it receives to it's associated
destination topic.

# Example Usage

mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 simulation/src1 simulation/dst1 simulation/src2 simulation/dst2 simulation/src3 simulation/dst3