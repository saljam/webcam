package webcam

// #cgo CXXFLAGS: -DWEBRTC_POSIX
//
// #cgo CXXFLAGS: -pthread -std=c++11 -Wno-narrowing -Wno-write-strings
// #cgo CXXFLAGS: -Iwebrtc/src -Iwebrtc/src/webrtc
//
// #cgo LDFLAGS: -lstdc++ -lm -lcrypto -lexpat -ldl
// #cgo LDFLAGS: -framework CoreGraphics -framework CoreAudio
// #cgo LDFLAGS: -framework QTKit -framework ApplicationServices
// #cgo LDFLAGS: -framework AVFoundation -framework AudioToolbox
// #cgo LDFLAGS: -framework CoreServices -framework CoreFoundation -framework Foundation
// #cgo LDFLAGS: -L${SRCDIR}/webrtc/src/out/Release -Wl
// #cgo LDFLAGS: -lvideo_render_module_internal_impl -lvideo_render
// #cgo LDFLAGS: -lstunprober -lrtc_xmpp -lrtc_xmllite -lrtc_sound -lmetrics
// #cgo LDFLAGS: -ljingle_peerconnection -ljingle_p2p -ljingle_media
// #cgo LDFLAGS: -lisac_fix -lgmock -lgenperf_libs -lframe_editing_lib
// #cgo LDFLAGS: -lexpat -ldesktop_capture -lcommand_line_parser
// #cgo LDFLAGS: -lchannel_transport -lbwe_tools_util -lbwe_simulator
// #cgo LDFLAGS: -laudioproc_protobuf_utils
// #cgo LDFLAGS: -lwebrtc
// #cgo LDFLAGS: -lvoice_engine -lvideo_render_module -lvideo_processing -lusrsctplib
// #cgo LDFLAGS: -lsrtp -lrtc_p2p -lrtc_event_log -lrtc_base
// #cgo LDFLAGS: -lremote_bitrate_estimator -lpaced_sender
// #cgo LDFLAGS: -ljsoncpp -lgflags -lframe_generator -ldesktop_capture_differ_sse2
// #cgo LDFLAGS: -lboringssl
// #cgo LDFLAGS: -lbitrate_controller
// #cgo LDFLAGS: -laudio_device -laudio_conference_mixer
// #cgo LDFLAGS: -lwebrtc_video_coding -lwebrtc_utility -lwebrtc_i420 -lwebrtc_h264
// #cgo LDFLAGS: -lvideo_processing_sse2 -lrtp_rtcp -lrtc_event_log_proto
// #cgo LDFLAGS: -lmedia_file
// #cgo LDFLAGS: -lfield_trial_default -lfield_trial
// #cgo LDFLAGS: -laudio_coding_module -lwebrtc_vp9 -lwebrtc_vp8 -lvideo_coding_utility
// #cgo LDFLAGS: -lrent_a_codec -lred -lneteq -lisac_common -lilbc
// #cgo LDFLAGS: -lg722 -lcng -lwebrtc_opus -lwebrtc_common -lpcm16b -lopus -lg711
// #cgo LDFLAGS: -lvpx_new -lvpx_intrinsics_sse4_1 -lvpx_intrinsics_sse2
// #cgo LDFLAGS: -lvpx_intrinsics_mmx -lvpx_intrinsics_avx2
// #cgo LDFLAGS: -lvpx_intrinsics_avx -lvpx_intrinsics_ssse3
// #cgo LDFLAGS: -lvideo_capture_module -lvideo_capture -lcommon_video -lyuv -ljpeg_turbo
// #cgo LDFLAGS: -lprotobuf_lite -laudio_processing_sse2 -laudio_processing
// #cgo LDFLAGS: -lmetrics_default -lisac -lhistogram
// #cgo LDFLAGS: -lcommon_audio -laudioproc_debug_proto
// #cgo LDFLAGS: -laudio_encoder_interface -laudio_decoder_interface
// #cgo LDFLAGS: -lopenmax_dl -lcommon_audio_sse2
// #cgo LDFLAGS: -lsystem_wrappers -lrtc_base_approved
//
// #include "peerconn.h"
import "C"
import (
	"log"
	"unsafe"
)

//export Debug
func Debug(s *C.char) {
	log.Println("c++:", C.GoString(s))
}

//export callbackCandidate
func callbackCandidate(pc unsafe.Pointer, sdp, mid *C.char, line C.int) {
	conns[pc].Candidate <- candidateMsg{
		Sdp:  C.GoString(sdp),
		Mid:  C.GoString(mid),
		Line: int(line),
	}
}

//export callbackOffer
func callbackOffer(pc unsafe.Pointer, sdp *C.char) {
	log.Println("doing callback for 0x%x", pc)
	log.Println("details", conns, conns[pc], conns[pc].Offer)
	conns[pc].Offer <- C.GoString(sdp)
}

func init() {
	C.init()
}

type PeerConn struct {
	Candidate chan candidateMsg
	Offer     chan string
	unsafe.Pointer
}

// A map to make looking up the right channel easier for callbacks.
// The alternative would be to send method pointers as callbacks but
// that's messy with cgo and we need a container anyway.
var conns = make(map[unsafe.Pointer]PeerConn)

func Offer() PeerConn {
	// TODO put this stuff in a NewPeerConn func and make Offer synchronous.
	pc := PeerConn{
		Candidate: make(chan candidateMsg),
		Offer:     make(chan string),
		Pointer:   C.Offer(),
	}
	conns[pc.Pointer] = pc
	return pc
}

func (pc PeerConn) AddAnswer(sdp string) {
	C.AddAnswer(pc.Pointer, C.CString(sdp))
}

func (pc PeerConn) AddCandidate(sdp, mid string, line int) {
	C.AddCandidate(pc.Pointer, C.CString(sdp), C.CString(mid), C.int(line))
}

func (pc PeerConn) Del(sdp, mid string, line int) {
	C.AddCandidate(pc.Pointer, C.CString(sdp), C.CString(mid), C.int(line))
}
