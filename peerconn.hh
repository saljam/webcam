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
	bool connection_active() const;
	void Connect();
	void Answer(std::string& sdp);
	void AddCandidate(std::string& sdp,std::string& mid, int line);
	void Close();

protected:
	talk_base::scoped_refptr<webrtc::PeerConnectionInterface> peer_connection_;
	talk_base::scoped_refptr<webrtc::PeerConnectionFactoryInterface> peer_connection_factory_;

	bool InitializePeerConnection();
	void AddStreams();
	cricket::VideoCapturer* OpenVideoCaptureDevice();

	// PeerConnectionObserver implementation.
	virtual void OnError();
	virtual void OnStateChange(webrtc::PeerConnectionObserver::StateType state_changed);
	virtual void OnAddStream(webrtc::MediaStreamInterface* stream);
	virtual void OnRemoveStream(webrtc::MediaStreamInterface* stream);
	virtual void OnRenegotiationNeeded();
	virtual void OnIceChange();
	virtual void OnIceCandidate(const webrtc::IceCandidateInterface* candidate);

	// CreateSessionDescriptionObserver implementation.
	virtual void OnSuccess(webrtc::SessionDescriptionInterface* desc);
	virtual void OnFailure(const std::string& error);
};
