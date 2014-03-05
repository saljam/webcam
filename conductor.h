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

#include "talk/app/webrtc/peerconnectioninterface.h"

namespace webrtc {
class VideoCaptureModule;
}

class Conductor :
public webrtc::PeerConnectionObserver,
public webrtc::CreateSessionDescriptionObserver {
public:
	Conductor(){};
	~Conductor(){};

	bool connection_active() const;
	void Connect();
	void Answer(std::string sdp);
	void AddCandidate(std::string sdp,std::string mid, int line);
	void Close();

protected:
	talk_base::scoped_refptr<webrtc::PeerConnectionInterface> peer_connection_;
	talk_base::scoped_refptr<webrtc::PeerConnectionFactoryInterface> peer_connection_factory_;

	bool InitializePeerConnection();
	void AddStreams();
	cricket::VideoCapturer* OpenVideoCaptureDevice();

	// PeerConnectionObserver implementation.
	virtual void OnError();
	virtual void OnStateChange(webrtc::PeerConnectionObserver::StateType state_changed) {}
	virtual void OnAddStream(webrtc::MediaStreamInterface* stream);
	virtual void OnRemoveStream(webrtc::MediaStreamInterface* stream);
	virtual void OnRenegotiationNeeded() {}
	virtual void OnIceChange() {}
	virtual void OnIceCandidate(const webrtc::IceCandidateInterface* candidate);

	// CreateSessionDescriptionObserver implementation.
	virtual void OnSuccess(webrtc::SessionDescriptionInterface* desc);
	virtual void OnFailure(const std::string& error);
};
