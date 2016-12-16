# Description

This is a simple program that maps a set of mqtt source topics to a set
of destination topics. Additionally, you can select parts of the source
topic's message using path expressions, assuming the message is JSON.
This is done at the mqtt application/client level. It subscribes to all
source topics and publishes any messages it receives to it's associated
destination topic. If a path expression is given, it will select the
relevant part of the message using [gjson][gjson]'s GetMany method. The output
message will be a JSON array containing the selected parts.

To specify a JSON path expression for a source, you append the [gjson][gjson]
path expression to the end of a source topic with two commas, ",,",
separating them. Multiple selections can be made by simply appending
them with the ",," separator between them.

I understand the ",," is ugly, but it is pretty clear that ",," doesn't
appear in english names or sentence, isn't used as a wildcard, nor doesn
it appear in [gjson][gjson]'s path expression description.
I am open to suggestions about a better separator.

# Example Usage

```
mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 simulation/src1 simulation/dst1 simulation/src2 simulation/dst2 simulation/src3 simulation/dst3
```
```
mqttappbridge -mqtt_user me -mqtt_pass somepassword -mqtt_server tls://mqtt.example.com:1883 simulation/src1 simulation/dst1 simulation/src2,,name.first,,name.last simulation/dst2 simulation/src3,,1.id simulation/dst3
```

# Dependencies

go get -u [github.com/tidwall/gjson][gjson]

[gjson]: http://github.com/tidwall/gjson