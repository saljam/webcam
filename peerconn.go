package webcam

// #cgo CFLAGS: -I/Library/Frameworks/OpenWebRTC.framework/Headers
// #cgo LDFLAGS: -framework OpenWebRTC
// #include <owr/owr.h>
// #include <owr/owr_session.h>
// #include <owr/owr_media_session.h>
//
// void init();
// OwrSession* new_session();
// void del_session(OwrSession* session);
import "C"
import "unsafe"

func init() {
	C.init()
}

type Session struct {
	Candidates chan Candidate
	// Pointer to the Session object in C-land.
	session *C.OwrSession
	// Pointer to the TransportAgent object in C-land.
	transport *C.transport_agent
}

func NewSession() Session {
	session := C.new_session()
	s := Session{
		Candidates: make(chan Candidate),
		session:    session,
	}
	sessions[session] = s.Candidates
	return s
}

// Offer generates a new WebRTC offer.
func (s Session) Description() (offer string) {

	s.WaitForCandodates()
	sinfo := C.get_session_info(s.session)

	desc := jsonSessionDesc{
		MediaDesc: []jsonMediaDesc{{
			Type: sinfo.media_type,
			Rtcp: {Mux: sinfo.rtcp_mux},
			Payloads: []jsonPayload{{
				Type:         sinfo.payload_type,
				ClockRate:    sinfo.clock_rate,
				EncodingName: sinfo.encoding_name,
				Ccmfir:       sinfo.ccm_fir,
				Nackpli:      sinfo.nack_pli,
				// TODO add parameters here for H264
				Channels: sinfo.channels,
			}},
			ICE: {
				Ufrag:    s.Candidates[0].Ufrag,
				Password: s.Candidates[0].Password,
				// TODO candidates marshal xml
				Candidates: s.Candidates,
			},
		}},
	}

}

// Offer generates a new WebRTC offer.
func (s Session) Offer() (offer string) {
	// generate sdp from pc.session
	return ""
}

// Answer generates an answer from an offer.
func (s Session) Answer(offer string) (answer string) {
	// generate sdp from pc.session
	return ""
}

// Accept receives an offer's answer.
func (s Session) Accept(answer string) {
	// C.AddAnswer(pc.Pointer, C.CString(sdp))
}

func (s Session) AddCandidate(sdp, mid string, line int) {
	// C.owr_session_add_remote_candidate()
}

func (s Session) Close() error {
	C.del_session(s.session, s.transport)
	return nil
}

// Map of sessions to find the correct channels to use for the candidate callbacks.
var sessions = make(map[*C.OwrSession]chan Candidate)

type Candidate struct {
	Ufrag, Password string
	CandidateType   int
	Component       int
	Foundation      string
	Priority        int
	Transport       int
	Address         string
	Port            int
	BaseAddress     string
	BasePort        int
}

//export got_candidate_go
func got_candidate_go(session *C.OwrSession, ufrag, password *C.char, candidateType, component int, foundation *C.char, priority int, transportType int, port, basePort int, address, baseAddress *C.char) {
	sessions[session] <- Candidate{
		Ufrag: C.GoString(ufrag), Password: C.GoString(password),
		CandidateType: candidateType, Component: component,
		Foundation: C.GoString(foundation), Priority: priority,
		Transport: transportType,
		Address:   C.GoString(address), Port: port,
		BaseAddress: C.GoString(baseAddress), BasePort: basePort,
	}
}

//export candidate_gathering_done_go
func candidate_gathering_done_go(session *C.OwrSession, pc unsafe.Pointer) {
	ch := sessions[session]
	close(ch)
	delete(sessions, session)
}

// We don't need this...
const sdpTemplate = `v=0
o=- 6909453319602664734 2 IN IP4 127.0.0.1
s=-
t=0 0
a=msid-semantic: WMS
`

const h264params = `parameters: { levelAsymmetryAllowed: 1, packetizationMode: 1 profileLevelId: "42e01f" }`

const descTemplate = `{
"mediaDescriptions": [{
  "type": s.media-type,
  "rtcp": {
    "mux": s.rtcp-mux
  },
  "payloads": [{ 
    "encodingName": s.encoding-name
    "type": s.payload-type
    "clockRate": s.clock-rate
    {{ parameters }}
    "channels": s.channels
	"ccmfir": s.ccm-fir
    "nackpli": s.nack-pli
  }],
  "ice": {
    "ufrag": ice_ufrag,
    "password": ice_password,
    "candidates": [{{candidates}}]
  }
  "dtls": {
    "fingerprintHashFunction": "sha-256"
    "fingerprint": session.fingerprint
    "setup": "active"
  }
}]
}`

const candidateTemplate = `{
"foundation": c.foundation
"componentId": c.Component
"transport": "UDP" || "TCP"
"priority": c.Priority
"address": c.Address
"port": c.Port
"type": candidate_types[candidate_type]
omitempty "relatedAddress": baseaddress
omitempty "relaredPort": baseport
omitempty "tcpType": tcp_types[transport_type]
}`

type jsonCandidate struct {
	CandidateType int    `type`
	Component     int    `componentId`
	Foundation    string `foundation`
	Priority      int    `priority`
	Transport     string `transport`
	Address       string `address`
	Port          int    `port`
	BaseAddress   string `omitempty,relatedAddress`
	BasePort      int    `omitempty,relaredPort`
	TCPType       int    `omitempty,tcpType`
}

type jsonPayload struct {
	Type      string `type`
	Encoding  string `encodingName`
	ClockRate string `jsonPayloadclockRate`
	Channels  string `omitempty,channels`
	Ccmfir    string `omitempty,ccmfir`
	Nackpli   string `omitempty,nackpli`
}

type jsonMediaDesc struct {
	Type string `type`
	Rtcp struct {
		Mux bool `mux`
	} `rtcp`
	Payloads []jsonPayload `payloads`
	ICE      struct {
		Ufrag      string `ufrag`
		Password   string `password`
		Candidates []Candidate
	} `ice`
	// DTLS
}

type jsonSessionDesc struct {
	MediaDesc []jsonMediaDesc `mediaDescriptions`
}

type jsonMessage struct {
	Type        string          `type`
	SessionDesc jsonSessionDesc `sessionDescription`
}
