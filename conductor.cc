/*
 * libjingle
 * Copyright 2012, Google Inc.
 * Copyright 2014, Salman Aljammaz
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  1. Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *  2. Redistributions in binary form must reproduce the above copyright notice,
 *     this list of conditions and the following disclaimer in the documentation
 *     and/or other materials provided with the distribution.
 *  3. The name of the author may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO
 * EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
 * OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
 * OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
 * ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#include "talk/app/webrtc/videosourceinterface.h"
#include "talk/media/devices/devicemanager.h"

#include "conductor.h"

class DummySetSessionDescriptionObserver :
public webrtc::SetSessionDescriptionObserver {
public:
	static DummySetSessionDescriptionObserver* Create() {
		return new talk_base::RefCountedObject<DummySetSessionDescriptionObserver>();
	}
	virtual void OnSuccess() {
		printf("DummySetSessionDescriptionObserver succeeded\n");
	}
	virtual void OnFailure(const std::string& error) {
		printf("DummySetSessionDescriptionObserver failed\n");
	}

protected:
	DummySetSessionDescriptionObserver() {}
	~DummySetSessionDescriptionObserver() {}
};

bool Conductor::connection_active() const {
	return peer_connection_.get() != NULL;
}

void Conductor::Close() {
	peer_connection_ = NULL;
	peer_connection_factory_ = NULL;
}

bool Conductor::InitializePeerConnection() {
	peer_connection_factory_  = webrtc::CreatePeerConnectionFactory();

	if (!peer_connection_factory_.get()) {
		printf("Failed to initialize PeerConnectionFactory\n");
		Close();
		return false;
	}

	webrtc::PeerConnectionInterface::IceServers servers;
	webrtc::PeerConnectionInterface::IceServer server;
	server.uri = "stun:stun.l.google.com:19302";
	servers.push_back(server);
	peer_connection_ = peer_connection_factory_->CreatePeerConnection(servers, NULL, NULL, this);
	if (!peer_connection_.get()) {
		printf("CreatePeerConnection failed\n");
		Close();
	}
	AddStreams();
	return peer_connection_.get() != NULL;
}

// PeerConnectionObserver implementation.

void Conductor::OnError() {
}

// Called when a remote stream is added. We don't need this.
void Conductor::OnAddStream(webrtc::MediaStreamInterface* stream) {
	printf("add stream\n");
	stream->AddRef();
}

void Conductor::OnRemoveStream(webrtc::MediaStreamInterface* stream) {
	printf("remove stream\n");
	stream->AddRef();
}

void Conductor::OnIceCandidate(const webrtc::IceCandidateInterface* candidate) {
	printf("got ice\n");
	std::string mid = candidate->sdp_mid();
	int line = candidate->sdp_mline_index();
	std::string sdp;
	if (!candidate->ToString(&sdp)) {
		printf("Failed to serialize candidate\n");
		return;
	}
	//candidate <- {sdp, mid, line}
}

// Here we make an offer...
void Conductor::Connect() {
	if (!InitializePeerConnection()) {
		printf("Failed to initialize PeerConnection\n");
		return;
	}
	peer_connection_->CreateOffer(this, NULL);
}

void Conductor::Answer(std::string sdp) {
	webrtc::SessionDescriptionInterface* session_description(webrtc::CreateSessionDescription("answer", sdp));
	peer_connection_->SetRemoteDescription(DummySetSessionDescriptionObserver::Create(), session_description);
}

void Conductor::AddCandidate(std::string sdp, std::string mid, int line) {
	talk_base::scoped_ptr<webrtc::IceCandidateInterface> candidate(webrtc::CreateIceCandidate(mid, line, sdp));
	if (!candidate.get()) {
		printf("Failed to parse candidate\n");
		return;
	}
	if (!peer_connection_->AddIceCandidate(candidate.get())) {
		printf("Failed to add candidate\n");
		return;
	}
	return;
}

cricket::VideoCapturer* Conductor::OpenVideoCaptureDevice() {
	talk_base::scoped_ptr<cricket::DeviceManagerInterface> dev_manager(
	cricket::DeviceManagerFactory::Create());
	if (!dev_manager->Init()) {
		printf("Can't create device manager\n");
		return NULL;
	}
	std::vector<cricket::Device> devs;
	if (!dev_manager->GetVideoCaptureDevices(&devs)) {
		printf("Can't enumerate video devices\n");
		return NULL;
	}
	std::vector<cricket::Device>::iterator dev_it = devs.begin();
	cricket::VideoCapturer* capturer = NULL;
	for (; dev_it != devs.end(); ++dev_it) {
	capturer = dev_manager->CreateVideoCapturer(*dev_it);
	if (capturer != NULL)
		break;
	}
	return capturer;
}

void Conductor::AddStreams() {
	talk_base::scoped_refptr<webrtc::AudioTrackInterface> audio_track(
		peer_connection_factory_->CreateAudioTrack(
			"audio_label", peer_connection_factory_->CreateAudioSource(NULL)));

	talk_base::scoped_refptr<webrtc::VideoTrackInterface> video_track(
		peer_connection_factory_->CreateVideoTrack(
			"video_label", peer_connection_factory_->CreateVideoSource(OpenVideoCaptureDevice(), NULL)));

	talk_base::scoped_refptr<webrtc::MediaStreamInterface> stream =
		peer_connection_factory_->CreateLocalMediaStream("stream_label");

	stream->AddTrack(audio_track);
	stream->AddTrack(video_track);
	if (!peer_connection_->AddStream(stream, NULL)) {
		printf("Adding stream to PeerConnection failed\n");
	}
}

// CreateSessionDescriptionObserver implementation.

void Conductor::OnSuccess(webrtc::SessionDescriptionInterface* desc) {
	peer_connection_->SetLocalDescription(DummySetSessionDescriptionObserver::Create(), desc);
	std::string sdp;
	desc->ToString(&sdp);
	//offer <- sdp;
}

void Conductor::OnFailure(const std::string& error) {
	printf("SDP Error, oops I can't print what it is because c++ is pants.\n");
}
