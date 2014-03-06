/*
 * Copyright 2014, Salman Aljammaz
 * Copyright 2012, Google Inc.
 */


#include "talk/app/webrtc/videosourceinterface.h"
#include "talk/media/devices/devicemanager.h"
#include "talk/base/ssladapter.h"

#include "peerconn.hh"

extern "C" {
#include "_cgo_export.h"
#include "peerconn.h"
}

class SetSDHandler :
public webrtc::SetSessionDescriptionObserver {
public:
	static SetSDHandler* Create() {
		return new talk_base::RefCountedObject<SetSDHandler>();
	}
	virtual void OnSuccess();
	virtual void OnFailure(const std::string& error);

protected:
	SetSDHandler() {}
	~SetSDHandler() {}
};

void SetSDHandler::OnSuccess() {
		Debug("SetSDHandler succeeded");
}

void SetSDHandler::OnFailure(const std::string& error) {
		Debug("SetSDHandler failed");
		Debug((char*)error.c_str());
}

bool Conductor::connection_active() const {
	return peer_connection_.get() != NULL;
}

void Conductor::Close() {
//	peer_connection_ = NULL;
//	peer_connection_factory_ = NULL;
}

bool Conductor::InitializePeerConnection() {
	talk_base::InitializeSSL();
	Debug("Initializing Peer Connection");
	peer_connection_factory_  = webrtc::CreatePeerConnectionFactory();

	if (!peer_connection_factory_.get()) {
		Debug("Failed to initialize PeerConnectionFactory\n");
		Close();
		return false;
	}

	webrtc::PeerConnectionInterface::IceServers servers;
	webrtc::PeerConnectionInterface::IceServer server;
	server.uri = "stun:stun.l.google.com:19302";
	servers.push_back(server);
	peer_connection_ = peer_connection_factory_->CreatePeerConnection(servers, NULL, NULL, this);
	if (!peer_connection_.get()) {
		Debug("CreatePeerConnection failed");
		Close();
	}
	AddStreams();
	return peer_connection_.get() != NULL;
}

// PeerConnectionObserver implementation.

void Conductor::OnError() {
	Debug("ERROR");
}

// Called when a remote stream is added. We don't need this.
void Conductor::OnAddStream(webrtc::MediaStreamInterface* stream) {
	Debug("add stream\n");
}

void Conductor::OnRemoveStream(webrtc::MediaStreamInterface* stream) {
	Debug("remove stream\n");
}

void Conductor::OnStateChange(webrtc::PeerConnectionObserver::StateType state_changed) {}
void Conductor::OnRenegotiationNeeded() {}
void Conductor::OnIceChange() {}

void Conductor::OnIceCandidate(const webrtc::IceCandidateInterface* candidate) {
	Debug("got ice");
	std::string mid = candidate->sdp_mid();
	int line = candidate->sdp_mline_index();
	std::string sdp;
	if (!candidate->ToString(&sdp)) {
		Debug("Failed to serialize candidate\n");
		return;
	}
	RegisterCandidate((char*)sdp.c_str(), (char*)mid.c_str(), line);
}

// Here we make an offer...
void Conductor::Connect() {
	if (!InitializePeerConnection()) {
		Debug("Failed to initialize PeerConnection\n");
		return;
	}
	peer_connection_->CreateOffer(this, NULL);
}

void Conductor::Answer(std::string& sdp) {
	Debug("doing answer");
	webrtc::SessionDescriptionInterface* session_description(webrtc::CreateSessionDescription("answer", sdp));
	if (!session_description) {
		Debug("foooooo");
		return;
	}
	peer_connection_->SetRemoteDescription(SetSDHandler::Create(), session_description);
}

void Conductor::AddCandidate(std::string& sdp, std::string& mid, int line) {
	talk_base::scoped_ptr<webrtc::IceCandidateInterface> candidate(webrtc::CreateIceCandidate(mid, line, sdp));
	if (!candidate.get()) {
		Debug("Failed to parse candidate\n");
		return;
	}
	if (!peer_connection_->AddIceCandidate(candidate.get())) {
		Debug("Failed to add candidate\n");
		return;
	}
	return;
}

cricket::VideoCapturer* Conductor::OpenVideoCaptureDevice() {
	talk_base::scoped_ptr<cricket::DeviceManagerInterface> dev_manager(
	cricket::DeviceManagerFactory::Create());
	if (!dev_manager->Init()) {
		Debug("Can't create device manager");
		return NULL;
	}
	std::vector<cricket::Device> devs;
	if (!dev_manager->GetVideoCaptureDevices(&devs)) {
		Debug("Can't enumerate video devices");
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
	Debug("Adding streams");
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
		Debug("Adding stream to PeerConnection failed\n");
	}
}

// CreateSessionDescriptionObserver implementation. (For CreateOffer)

void Conductor::OnSuccess(webrtc::SessionDescriptionInterface* desc) {
	peer_connection_->SetLocalDescription(SetSDHandler::Create(), desc);
	std::string sdp;
	desc->ToString(&sdp);
	RegisterOffer(strdup(sdp.c_str()));
}

void Conductor::OnFailure(const std::string& error) {
	Debug((char *)error.c_str());
}

// C Interface

talk_base::scoped_refptr<Conductor> pc;

void InitPeerConn() {
	pc  = new talk_base::RefCountedObject<Conductor>();
	pc->Connect();
}

void Answer(char* sdp) {
	Debug("got answer");
	std::string csdp = std::string(sdp);
	pc->Answer(csdp);
}

void Candidate(char* sdp,char* mid, int line) {
	std::string csdp = std::string(sdp);
	std::string cmid = std::string(mid);
	pc->AddCandidate(csdp, cmid, (int)line);
}
