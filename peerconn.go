package webcam

// #cgo CFLAGS: -I/Library/Frameworks/OpenWebRTC.framework/Headers
// #cgo LDFLAGS: -framework OpenWebRTC
// #include <owr/owr.h>
// #include <owr/owr_session.h>
// #include <owr/owr_media_session.h>
//
// OwrMediaSession* new_session();
import "C"
import (
	"unsafe"
)

type PeerConn struct {
	Candidate chan candidateMsg
	// pointer to the session object in C-land
	session   *C.OwrMediaSession
}

func NewPeerConn() PeerConn{
	session := C.new_session()
	pc := PeerConn{
		Candidate: make(chan candidateMsg),
		session:   session,
	}
	return pc
}

func (pc PeerConn) GenerateOffer() string {
	// generate sdp from pc.session
	return ""
}

func (pc PeerConn) GenerateAnswer(offer string) string {
	// generate sdp from pc.session
	return ""
}

func (pc PeerConn) Answer(sdp string) {
	// C.AddAnswer(pc.Pointer, C.CString(sdp))
}

func (pc PeerConn) AddCandidate(sdp, mid string, line int) {
	// C.owr_session_add_remote_candidate()
}

// Map of sessions to find the correct channels to use for the candidate callbacks.
var sessions = make(map[*C.OwrSession]chan candidateMsg)

//export cbCandidate
func cbCandidate(session *C.OwrSession, candidate *C.OwrCandidate, pc unsafe.Pointer) {
	ch := sessions[session]
	ch <- candidateMsg{
		Sdp:  "",
		Mid:  "",
		Line: 0,
	}
}

// I *think* this is what they mean...
//export cbCandidatesDone
func cbCandidatesDone(session *C.OwrSession, pc unsafe.Pointer) {
	ch := sessions[session]
	close(ch)
}
