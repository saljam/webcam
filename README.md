# Web Browser ← Camera

A one-way camera stream using the OpenWebRTC native library.

### Building

This depends on the [OpenWebRTC](http://www.openwebrtc.org/) library. As long as it's installed on your system (there are [binary releases](https://github.com/EricssonResearch/openwebrtc/releases)) `go get 0f.io/webcam` should just work on OS X. Linux/Windows should follow soon.

### Using

The command usage is:

    Usage of webcam:
      -addr string
        	http address to listen on (default ":8003")

For the webrtc package usage refer to the godoc page on [godoc.org](https://godoc.org/0f.io/webcam/webrtc).
