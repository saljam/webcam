package webrtc

import (
	"encoding/json"
	"testing"
)

func TestCandidates(t *testing.T) {
	s := NewSession()
	ncandidates := 0
	for c := range s.candidatec {
		t.Log(c)
		ncandidates++
	}
	s.Close()
	if ncandidates == 0 {
		t.Fatalf("got zero ICE candidates for new session")
	}
}

func TestDescription(t *testing.T) {
	s := NewSession()
	desc, err := s.Description()
	if err != nil {
		t.Fatalf("error marshalling description json: %v", err)
	}
	t.Logf("session description: %s", desc)
	s.Close()
}

func TestUnmarshalDesc(t *testing.T) {
	in := `{"version":0,"originator":{"username":"mozilla...THIS_IS_SDPARTA-42.0","sessionId":"1727442231039758385","sessionVersion":0,"netType":"IN","addressType":"IP4","address":"0.0.0.0"},"sessionName":"-","startTime":0,"stopTime":0,"mediaDescriptions":[{"type":"video","port":39906,"protocol":"RTP/SAVPF","netType":"IN","addressType":"IP4","address":"212.102.22.168","mode":"recvonly","payloads":[{"type":120,"encodingName":"VP8","clockRate":90000,"nack":true,"nackpli":true,"ccmfir":true,"ericscream":false,"parameters":{"maxFs":12288,"maxFr":60}},{"type":126,"encodingName":"H264","clockRate":90000,"nack":true,"nackpli":true,"ccmfir":true,"ericscream":false,"parameters":{"profileLevelId":"42e01f","levelAsymmetryAllowed":1,"packetizationMode":1}},{"type":97,"encodingName":"H264","clockRate":90000,"nack":true,"nackpli":true,"ccmfir":true,"ericscream":false,"parameters":{"profileLevelId":"42e01f","levelAsymmetryAllowed":1}}],"rtcp":{"netType":"IN","port":53045,"addressType":"IP4","address":"212.102.22.168","mux":true},"ssrcs":[3660879275],"cname":"{3423fb43-b1ac-db4d-8fba-249a935fdfe2}","ice":{"ufrag":"bcee3fc6","password":"82ad6ed4e4eb6fd93f3ad448e80d20f6","iceOptions":{"trickle":true},"candidates":[{"foundation":"0","componentId":1,"transport":"UDP","priority":2130379007,"address":"192.168.30.63","port":61124,"type":"host"},{"foundation":"0","componentId":2,"transport":"UDP","priority":2130379006,"address":"192.168.30.63","port":63125,"type":"host"},{"foundation":"1","componentId":1,"transport":"UDP","priority":1694179327,"address":"212.102.22.168","port":39906,"type":"srflx","relatedAddress":"192.168.30.63","relatedPort":61124},{"foundation":"1","componentId":2,"transport":"UDP","priority":1694179326,"address":"212.102.22.168","port":53045,"type":"srflx","relatedAddress":"192.168.30.63","relatedPort":63125}]},"dtls":{"fingerprintHashFunction":"sha-256","fingerprint":"02:D5:A7:E9:95:DF:53:8B:AC:2F:9A:84:FC:D9:D0:00:E7:0B:B4:6A:60:F1:94:B7:70:0D:E3:2E:62:EE:E4:65","setup":"actpass"}}]}`
	var desc sessionDesc
	err := json.Unmarshal([]byte(in), &desc)
	if err != nil {
		t.Fatalf("Couldn't unmarshal session description: %v", err)
	}
	t.Log(desc)
}
