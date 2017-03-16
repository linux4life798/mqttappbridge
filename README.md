# Description

This is a simple program that maps a set of mqtt source topics to a set
of destination topics. Additionally, you can select parts of the source
topic's message using path expressions, assuming the message is JSON.
This is done at the mqtt application/client level. It subscribes to all
source topics and publishes any messages it receives to it's associated
destination topic. If a [path expression][gjson_path] is given, it will select the
relevant part of the message using [gjson][gjson]'s GetMany method. The
output message will be a JSON array containing the selected parts, unless
the "simple_single" flag is specified. If this flag is used, single
JSON Paths will present just the bare value, instead of JSON array
wrapped single value.

To specify a JSON path expression for a source, you append the
[gjson path expression][gjson_path] to the end of a source topic with
two commas, ",,", separating them. Multiple selections can be made
by simply appending them with the ",," separator between them.

I understand the ",," is ugly, but it is pretty clear that ",," doesn't
appear in english names or sentence, isn't used as a wildcard, nor doesn
it appear in [gjson][gjson]'s path expression description.
I am open to suggestions about a better separator.

# Usage
```
mqttappbridge [options] <source_topic1> <destination_topic1> [<source_topic2> <destination_topic2> [...]]

option:
  -mqtt_pass string
    	Sets the MQTT password
  -mqtt_server string
    	Sets the MQTT server (default "tcp://localhost:1883")
  -mqtt_user string
    	Sets the MQTT username
  -simple_single
    	Removes JSON array brackets from single values
```

# Example Usage

```
mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 simulation/src1 simulation/dst1 simulation/src2 simulation/dst2 simulation/src3 simulation/dst3
```
```
mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 simulation/src1 simulation/dst1 simulation/src2,,name.first,,name.last simulation/dst2 simulation/src3,,1.id simulation/dst3
```
```
mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 -simple_single simulation/+/rx simulation/allrx simulation/+/rx,,data simulation/allrxdata
```

# Dependencies

go get -u [github.com/tidwall/gjson][gjson]

[gjson]: http://github.com/tidwall/gjson
[gjson_path]: https://github.com/tidwall/gjson#path-syntax
