package webcam

// #cgo CFLAGS: -I/Library/Frameworks/OpenWebRTC.framework/Headers
// #cgo LDFLAGS: -framework OpenWebRTC
// #include <owr/owr.h>
// #include <owr/owr_session.h>
// #include <owr/owr_media_session.h>
//
// void init();
// OwrSession* new_session();
import "C"
import "unsafe"

func init() {
	C.init()
}

type Session struct {
	Candidates chan Candidate
	// pointer to the session object in C-land
	session *C.OwrSession
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

func (pc Session) AddCandidate(sdp, mid string, line int) {
	// C.owr_session_add_remote_candidate()
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

const offerTemplate = `v=0
o=- 6909453319602664734 2 IN IP4 127.0.0.1
s=-
t=0 0
a=msid-semantic: WMS
`
