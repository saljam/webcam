package webcam

var template = `<polymer-element name="webrtc-stream">
<template>
<video id=cam autoplay>
</template>
<style>
video {
	width: 100%%;
}
</style>
<script>
function error(err) {
	console.log("err:", err)
}

Polymer('webrtc-stream', {
	ready: function() {
		delete URL; // Polymer's platform.js overwrites URL with one that doesn't have createObjectURL().
		
		var ws = new WebSocket('ws://' + location.host + %s)
		var cfg = {"iceServers": [{"url": "stun:stun.l.google.com:19302"}]};
		pc = new RTCPeerConnection(cfg, {optional: [{RtpDataChannels: true}]});
		
		ws.onmessage = function(m) {
			var msg = JSON.parse(m.data);
			if (msg.type === 'offer') {
				pc.setRemoteDescription(new RTCSessionDescription(msg), function(){}, error);
				pc.createAnswer(function(desc) {
					pc.setLocalDescription(desc);
					ws.send(JSON.stringify(desc));
				}, error);
			} else if (msg.type === 'candidate') {
				pc.addIceCandidate(new RTCIceCandidate(msg));
			} else {
				console.log("what's this?", msg)
			}
		}
		
		pc.onicecandidate = function(e) {
			if (e.candidate) {
				e.candidate.type="candidate";
				console.log(e.candidate);
				ws.send(JSON.stringify(e.candidate));
			}
		}
		pc.onaddstream = function (e) {
			var vid = document.getElementById("cam");
			attachMediaStream(vid, e.stream);
		};
	}
});
</script>
</polymer-element>


<title>webcam</title>
<body>
`