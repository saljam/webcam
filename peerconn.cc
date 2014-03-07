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

class SetSDPHandler :
public webrtc::SetSessionDescriptionObserver {
public:
	static SetSDPHandler* Create() {
		return new talk_base::RefCountedObject<SetSDPHandler>();
	}
	virtual void OnSuccess() {
		Debug("SetSDPHandler succeeded");
	}
	virtual void OnFailure(const std::string& error) {
		Debug("SetSDPHandler failed");
		Debug((char*)error.c_str());
	}
protected:
	SetSDPHandler() {}
	~SetSDPHandler() {}
};

void Conductor::Close() {
	pc = NULL;
}

void Conductor::Offer() {
	// TODO add servers argument.
	webrtc::PeerConnectionInterface::IceServers servers;
	webrtc::PeerConnectionInterface::IceServer server;
	server.uri = "stun:stun.l.google.com:19302";
	servers.push_back(server);
	pc = peerConnectionFactory->CreatePeerConnection(servers, NULL, NULL, this);
	if (pc.get() == NULL) {
		Debug("CreatePeerConnection failed");
		Close();
		return;
	}
	AddStreams();
	pc->CreateOffer(this, NULL);
}

void Conductor::AddAnswer(std::string& sdp) {
	webrtc::SessionDescriptionInterface* session_description(webrtc::CreateSessionDescription("answer", sdp));
	if (!session_description) {
		return;
	}
	pc->SetRemoteDescription(SetSDPHandler::Create(), session_description);
}

void Conductor::AddCandidate(std::string& sdp, std::string& mid, int line) {
	talk_base::scoped_ptr<webrtc::IceCandidateInterface> candidate(webrtc::CreateIceCandidate(mid, line, sdp));
	if (!candidate.get()) {
		Debug("Failed to parse candidate");
		return;
	}
	if (!pc->AddIceCandidate(candidate.get())) {
		Debug("Failed to add candidate");
		return;
	}
	return;
}


void Conductor::AddStreams() {
	talk_base::scoped_refptr<webrtc::MediaStreamInterface> stream =
		peerConnectionFactory->CreateLocalMediaStream("stream_label");
	
	// We don't care about audio for now.	
	//talk_base::scoped_refptr<webrtc::AudioTrackInterface> audio_track(
	//	peerConnectionFactory->CreateAudioTrack(
	//		"audio_label", peerConnectionFactory->CreateAudioSource(NULL)));
	//stream->AddTrack(audio_track);

	talk_base::scoped_refptr<webrtc::VideoTrackInterface> video_track(
		peerConnectionFactory->CreateVideoTrack(
			"video_label", videoSource));
	stream->AddTrack(video_track);

	if (!pc->AddStream(stream, NULL)) {
		Debug("Adding stream to PeerConnection failed");
	}
}

// PeerConnectionObserver implementation.

void Conductor::OnIceCandidate(const webrtc::IceCandidateInterface* candidate) {
	std::string mid = candidate->sdp_mid();
	int line = candidate->sdp_mline_index();
	std::string sdp;
	if (!candidate->ToString(&sdp)) {
		Debug("Failed to serialize candidate");
		return;
	}
	callbackCandidate(this, (char*)sdp.c_str(), (char*)mid.c_str(), line);
}

// CreateSessionDescriptionObserver implementation.

void Conductor::OnSuccess(webrtc::SessionDescriptionInterface* desc) {
	pc->SetLocalDescription(SetSDPHandler::Create(), desc);
	std::string sdp;
	desc->ToString(&sdp);
	Debug("offer callback");
	callbackOffer(this, strdup(sdp.c_str()));
}

void Conductor::OnFailure(const std::string& error) {
	Debug((char *)error.c_str());
}

// C Interface

cricket::VideoCapturer* OpenVideoCaptureDevice() {
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


void init() {
	talk_base::InitializeSSL();
	
	Debug("Initializing PeerConnectionFactory");
	peerConnectionFactory  = webrtc::CreatePeerConnectionFactory();
	if (peerConnectionFactory.get() == NULL) {
		Debug("Failed to initialize PeerConnectionFactory");
		return;
	}
	
	videoSource = peerConnectionFactory->CreateVideoSource(OpenVideoCaptureDevice(), NULL);
}

void* Offer() {
	Conductor *pc = new Conductor();
	pc->Offer();
	return (void*)pc;
}

void AddAnswer(void* pc, char* sdp) {
	Conductor* cpc = (Conductor*)pc;
	std::string csdp = std::string(sdp);
	cpc->AddAnswer(csdp);
}

void AddCandidate(void* pc, char* sdp,char* mid, int line) {
	Conductor* cpc = (Conductor*)pc;
	std::string csdp = std::string(sdp);
	std::string cmid = std::string(mid);
	cpc->AddCandidate(csdp, cmid, (int)line);
}
