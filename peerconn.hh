/*
 * Copyright 2014, Salman Aljammaz
 * Copyright 2012, Google Inc.
 */

#include "talk/app/webrtc/peerconnectioninterface.h"

namespace webrtc {
class VideoCaptureModule;
}

class Conductor :
public webrtc::PeerConnectionObserver,
public webrtc::CreateSessionDescriptionObserver {
public:
	void Offer();
	void AddAnswer(std::string& sdp);
	void AddCandidate(std::string& sdp,std::string& mid, int line);
	void Close();

protected:
	talk_base::scoped_refptr<webrtc::PeerConnectionInterface> pc;

	void AddStreams();

	// PeerConnectionObserver implementation.
	virtual void OnError(){}
	virtual void OnStateChange(webrtc::PeerConnectionObserver::StateType state_changed){}
	virtual void OnAddStream(webrtc::MediaStreamInterface* stream){}
	virtual void OnRemoveStream(webrtc::MediaStreamInterface* stream){}
	virtual void OnRenegotiationNeeded(){}
	virtual void OnIceChange(){}
	virtual void OnIceCandidate(const webrtc::IceCandidateInterface* candidate);

	// CreateSessionDescriptionObserver implementation.
	virtual void OnSuccess(webrtc::SessionDescriptionInterface* desc);
	virtual void OnFailure(const std::string& error);
	virtual int AddRef(){} // We manage own memory.
	virtual int Release(){}
};

cricket::VideoCapturer* OpenVideoCaptureDevice();
talk_base::scoped_refptr<webrtc::VideoSourceInterface> videoSource;
talk_base::scoped_refptr<webrtc::PeerConnectionFactoryInterface> peerConnectionFactory;
