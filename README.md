# Web Browser ← Camera

A one-way camera stream using the WebRTC native libraries.

## Building

To use this package you have to run `make` first. That'll checkout the WebRTC source and build it. Make sure you've got the [prerequisite software](http://www.webrtc.org/native-code/development/prerequisite-sw). This usually takes a while.

Once that's done go build should work as usual.

## Using

Set up the stream handler by calling webcam.NewWebcam(), then register that with http.Handle().

	http.HandleFunc(<prefix>, webcam.HandleHTTP)

You can copy the JS code in webcam/ui/plainindex.html. The webcam handler also serves a Polymer element  to make the client side easier. Any request URL path ending with "webrtc-stream.html" will be given the Polymer element declaration. For example, a complete index.html could be:

	<!DOCTYPE html>
	<html>
	<script src="platform.js"></script>
	<link rel="import" href="polymer.html">
	<link rel="import" href="<prefix>/webrtc-stream.html">
	<title>webcam</title>
	<body>
	<webrtc-stream></webrtc-stream>

Where <prefix> is the prefix for the stream handler. Just make sure to import platform.js and polymer.html before you import webrtc-stream.
