/*
 * Copyright 2014-2015, Salman Aljammaz
 * Copyright 2012, Google Inc.
 *
 *  Use of this source code is governed by a BSD-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */

#include "webrtc/base/nethelpers.h"
#include "webrtc/base/physicalsocketserver.h"
#include "webrtc/base/scoped_ptr.h"
#include "webrtc/base/signalthread.h"
#include "webrtc/base/sigslot.h"

#include "webrtc/base/common.h"
#include "webrtc/base/nethelpers.h"

#include "talk/app/webrtc/test/fakeconstraints.h"

#include "talk/app/webrtc/videosourceinterface.h"
#include "talk/media/devices/devicemanager.h"
#include "base/ssladapter.h"

#include "peerconn.hh"

extern "C" {
#include "_cgo_export.h"
#include "peerconn.h"
}

class SetSDPHandler :
public webrtc::SetSessionDescriptionObserver {
public:
	static SetSDPHandler* Create() {
		return new rtc::RefCountedObject<SetSDPHandler>();
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

void WebcamConductor::Close() {
	pc = NULL;
}

void WebcamConductor::Offer() {
	// TODO add servers argument.
	webrtc::PeerConnectionInterface::IceServers servers;
	webrtc::PeerConnectionInterface::IceServer server;
	server.uri = "stun:stun.l.google.com:19302";
	servers.push_back(server);
	webrtc::FakeConstraints constraints;
	pc = peerConnectionFactory->CreatePeerConnection(servers, &constraints, NULL, NULL, this);
	if (pc.get() == NULL) {
		Debug("CreatePeerConnection failed");
		Close();
		return;
	}
	AddStreams();
	pc->CreateOffer(this, NULL);
}

void WebcamConductor::AddAnswer(std::string& sdp) {
	webrtc::SdpParseError error;
	webrtc::SessionDescriptionInterface* session_description(webrtc::CreateSessionDescription("answer", sdp, &error));
	// Let's ignore the error for now. Boo.
	if (!session_description) {
		return;
	}
	pc->SetRemoteDescription(SetSDPHandler::Create(), session_description);
}

void WebcamConductor::AddCandidate(std::string& sdp, std::string& mid, int line) {
	webrtc::SdpParseError error;
	rtc::scoped_ptr<webrtc::IceCandidateInterface> candidate(webrtc::CreateIceCandidate(mid, line, sdp, &error));
	// Let's ignore the error here too...
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


void WebcamConductor::AddStreams() {
	rtc::scoped_refptr<webrtc::MediaStreamInterface> stream =
		peerConnectionFactory->CreateLocalMediaStream("stream_label");
	
	// We don't care about audio for now.	
	//rtc::scoped_refptr<webrtc::AudioTrackInterface> audio_track(
	//	peerConnectionFactory->CreateAudioTrack(
	//		"audio_label", peerConnectionFactory->CreateAudioSource(NULL)));
	//stream->AddTrack(audio_track);

	rtc::scoped_refptr<webrtc::VideoTrackInterface> video_track(
		peerConnectionFactory->CreateVideoTrack(
			"video_label", videoSource));
	stream->AddTrack(video_track);

	if (!pc->AddStream(stream)) {
		Debug("Adding stream to PeerConnection failed");
	}
}

// PeerConnectionObserver implementation.

void WebcamConductor::OnIceCandidate(const webrtc::IceCandidateInterface* candidate) {
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

void WebcamConductor::OnSuccess(webrtc::SessionDescriptionInterface* desc) {
	pc->SetLocalDescription(SetSDPHandler::Create(), desc);
	std::string sdp;
	desc->ToString(&sdp);
	Debug("offer callback");
	callbackOffer(this, strdup(sdp.c_str()));
}

void WebcamConductor::OnFailure(const std::string& error) {
	Debug((char *)error.c_str());
}

int WebcamConductor::AddRef() const { return 1; }

int WebcamConductor::Release() const { return 1; }

// C Interface

cricket::VideoCapturer* OpenVideoCaptureDevice() {
	rtc::scoped_ptr<cricket::DeviceManagerInterface> dev_manager(
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
	rtc::InitializeSSL();
	
	Debug("Initializing PeerConnectionFactory");
	peerConnectionFactory = webrtc::CreatePeerConnectionFactory();
	if (peerConnectionFactory.get() == NULL) {
		Debug("Failed to initialize PeerConnectionFactory");
		return;
	}
	
	videoSource = peerConnectionFactory->CreateVideoSource(OpenVideoCaptureDevice(), NULL);
}

void* Offer() {
	WebcamConductor *pc = new WebcamConductor();
	pc->Offer();
	return (void*)pc;
}

void AddAnswer(void* pc, char* sdp) {
	WebcamConductor* cpc = (WebcamConductor*)pc;
	std::string csdp = std::string(sdp);
	cpc->AddAnswer(csdp);
}

void AddCandidate(void* pc, char* sdp,char* mid, int line) {
	WebcamConductor* cpc = (WebcamConductor*)pc;
	std::string csdp = std::string(sdp);
	std::string cmid = std::string(mid);
	cpc->AddCandidate(csdp, cmid, (int)line);
}
