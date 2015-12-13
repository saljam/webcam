package webcam

// #cgo CFLAGS: -I/Library/Frameworks/OpenWebRTC.framework/Headers
// #cgo LDFLAGS: -framework OpenWebRTC
// #include <owr/owr.h>
// #include <owr/owr_session.h>
// #include <owr/owr_media_session.h>
// #include <owr/owr_transport_agent.h>
// void init();
// OwrSession* new_session();
// void del_session(OwrSession* session);
// void add_candidate(OwrSession *session, char *ufrag, char *password, int type, int component, char *foundation, uint priority, int transport_type, char *address, uint port);
import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"unsafe"
)

// transportMutex protects operations on C.transport_agent.
var transportMutex sync.Mutex

type Session struct {
	// Pointer to the Session object in C-land.
	session         *C.OwrSession
	Candidates      []candidate
	candidatec      chan candidate
	Fingerprint     string
	fingerprintc    chan string
	ufrag, password string
}

func init() { C.init() }

var localSourcesReady = make(chan struct{})

// waitForLocalSources waits until OpenWebRTC has found and initialised the local sources.
func waitForLocalSources() {
	<-localSourcesReady
}

//export got_local_sources_go
func got_local_sources_go() {
	close(localSourcesReady)
}

// Map of C OwrSession pointers to Go Session pointers.
// This should probably be replaced with more locking in C...
var sessions = make(map[*C.OwrSession]*Session)

// The definition in the header files is different than the
// one in client_test.s. Should investigate.
var candidateTypes []string = []string{"host", "srflx", "relay"}
var tcpTypes []string = []string{"", "active", "passive", "so"}

//export got_candidate_go
func got_candidate_go(session *C.OwrSession, ufrag, password *C.char, candidateType, component int, foundation *C.char, priority uint, transportType int, port, basePort uint, address, baseAddress *C.char) {
	sessions[session].ufrag = C.GoString(ufrag)
	sessions[session].password = C.GoString(password)

	c := candidate{
		CandidateType: candidateTypes[candidateType],
		Component:     component,
		Foundation:    C.GoString(foundation), Priority: priority,
		TCPType: tcpTypes[transportType],
		Address: C.GoString(address), Port: port,
	}
	if transportType == 0 {
		c.Transport = "UDP"
	} else {
		c.Transport = "TCP"
	}

	if candidateType != 0 {
		c.BaseAddress = C.GoString(baseAddress)
		c.BasePort = basePort
	}

	sessions[session].candidatec <- c
}

//export got_dtls_certificate_go
func got_dtls_certificate_go(session *C.OwrSession, fingerprint *C.char) {
	sessions[session].Fingerprint = C.GoString(fingerprint)
	close(sessions[session].fingerprintc)
}

//export candidate_gathering_done_go
func candidate_gathering_done_go(session *C.OwrSession, pc unsafe.Pointer) {
	close(sessions[session].candidatec)
}

func NewSession() *Session {
	waitForLocalSources()
	session := C.new_session()
	s := &Session{
		candidatec:   make(chan candidate),
		fingerprintc: make(chan string),
		session:      session,
	}
	sessions[session] = s
	return s
}

// waitForCCallbacks waits for the resources which are tied to callbacks in C.
func (s *Session) waitForCCallbacks() {
	log.Printf("getting candidates")
	for c := range s.candidatec {
		s.Candidates = append(s.Candidates, c)
	}
	log.Printf("got all candidates")
	<-s.fingerprintc
	delete(sessions, s.session)
}

// Description generates the local session description.
func (s *Session) Description() (string, error) {
	s.waitForCCallbacks()

	// payload type: 120^W100
	// clockrate: 90000
	desc := sessionDesc{
		MediaDesc: []mediaDesc{{
			Type: "video",
			Mode: "sendonly",
			Payloads: []payload{{
				Type:      120,
				ClockRate: 90000,
				Encoding:  "VP8",
				Ccmfir:    true,
				Nackpli:   true,
			}},
		}},
	}

	desc.MediaDesc[0].Rtcp.Mux = true
	desc.MediaDesc[0].ICE.Ufrag = s.ufrag
	desc.MediaDesc[0].ICE.Password = s.password
	desc.MediaDesc[0].ICE.Candidates = s.Candidates
	desc.MediaDesc[0].DTLS.HashFunc = "sha-256"
	desc.MediaDesc[0].DTLS.Setup = "active"
	desc.MediaDesc[0].DTLS.Fingerprint = s.Fingerprint

	b, err := json.MarshalIndent(desc, "", "  ")
	return string(b), err
}

// Remote adds the description of the remote peer, including its ICE candidates.
// It expects remoteDescription to be a JSON encoded session description as produced
// by OpenWebRTC's SDP.js.
func (s *Session) Remote(remoteDescription []byte) error {
	var desc sessionDesc
	err := json.Unmarshal(remoteDescription, &desc)
	if err != nil {
		return fmt.Errorf("couldn't set remote description for session: %v", err)
	}
	ice := desc.MediaDesc[0].ICE
	for i := range ice.Candidates {
		log.Println("AddingCandidate")
		s.addCandidate(ice.Candidates[i], ice.Ufrag, ice.Password)
	}
	return nil
}

func (s Session) addCandidate(cand candidate, ufrag, password string) {
	var candidateType, transportType int
	switch cand.CandidateType {
	case "host":
		candidateType = 0
	case "srflx":
		candidateType = 1
	case "relay":
		candidateType = 2
	}

	switch cand.TCPType {
	case "active":
		transportType = 1
	case "passive":
		transportType = 2
	case "so":
		transportType = 3
	default:
		transportType = 0 // UDP
	}

	C.add_candidate(s.session, C.CString(ufrag), C.CString(password),
		C.int(candidateType), C.int(cand.Component),
		C.CString(cand.Foundation), C.uint(cand.Priority),
		C.int(transportType), C.CString(cand.Address), C.uint(cand.Port))

}

// Close frees the resources held by this session.
func (s Session) Close() error {
	transportMutex.Lock()
	C.del_session(s.session)
	transportMutex.Unlock()
	return nil
}

type candidate struct {
	CandidateType string `json:"type"`
	Component     int    `json:"componentId"`
	Foundation    string `json:"foundation"`
	Priority      uint   `json:"priority"`
	Transport     string `json:"transport"`
	Address       string `json:"address"`
	Port          uint   `json:"port"`
	BaseAddress   string `json:"relatedAddress,omitempty"`
	BasePort      uint   `json:"relaredPort,omitempty"`
	TCPType       string `json:"tcpType,omitempty"`
}

type payload struct {
	Type      int    `json:"type"`
	Encoding  string `json:"encodingName"`
	ClockRate uint   `json:"clockRate"`
	Channels  uint   `json:"channels,omitempty"`
	Ccmfir    bool   `json:"ccmfir"`
	Nackpli   bool   `json:"nackpli"`
}

type mediaDesc struct {
	Type string `json:"type"`
	Mode string `json:"mode"`
	Rtcp struct {
		Mux bool `json:"mux"`
	} `json:"rtcp"`
	Payloads []payload `json:"payloads"`
	ICE      struct {
		Ufrag      string      `json:"ufrag"`
		Password   string      `json:"password"`
		ICEOptions struct{}    `json:"iceOptions"`
		Candidates []candidate `json:"candidates"`
	} `json:"ice"`
	DTLS struct {
		HashFunc    string `json:"fingerprintHashFunction"`
		Fingerprint string `json:"fingerprint"`
		Setup       string `json:"setup"`
	} `json:"dtls"`
}

type sessionDesc struct {
	MediaDesc []mediaDesc `json:"mediaDescriptions"`
}

//export Debug
func Debug(s *C.char) {
	log.Println("C:", C.GoString(s))
}
