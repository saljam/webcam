# Web Browser ← Camera

A one-way camera stream using WebRTC.

## Building

To build this you'll first need to checkout and build the WebRTC native libraries into
the directory webrtc. (Or bind/symlink/etc. that to where it is.)

http://www.webrtc.org/reference/getting-started can guide you through that.

Once that's done go install should work as expected to produce the webcam binary.

## Using

Set up the stream handler by calling webcam.NewWebcam(), then register that with http.Handle().

	stream := webcam.NewWebcam()
	http.Handle(<prefix>, stream)

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

## Todo

 - Split out the WebRTC wrapper parts into their own package.
