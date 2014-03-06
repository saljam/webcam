package main

// #cgo CXXFLAGS: -DV8_DEPRECATION_WARNINGS -DEXPAT_RELATIVE_PATH
// #cgo CXXFLAGS: -DFEATURE_ENABLE_VOICEMAIL -DGTEST_RELATIVE_PATH
// #cgo CXXFLAGS: -DJSONCPP_RELATIVE_PATH -DLOGGING=1 -DSRTP_RELATIVE_PATH
// #cgo CXXFLAGS: -DFEATURE_ENABLE_SSL -DFEATURE_ENABLE_PSTN -DHAVE_SRTP
// #cgo CXXFLAGS: -DHAVE_WEBRTC_VIDEO -DHAVE_WEBRTC_VOICE
// #cgo CXXFLAGS: -DUSE_WEBRTC_DEV_BRANCH -D_FILE_OFFSET_BITS=64
// #cgo CXXFLAGS: -DCHROMIUM_BUILD -DTOOLKIT_VIEWS=1 -DUI_COMPOSITOR_IMAGE_TRANSPORT
// #cgo CXXFLAGS: -DUSE_AURA=1 -DUSE_CAIRO=1 -DUSE_GLIB=1
// #cgo CXXFLAGS: -DUSE_DEFAULT_RENDER_THEME=1 -DUSE_LIBJPEG_TURBO=1 -DUSE_NSS=1
// #cgo CXXFLAGS: -DUSE_X11=1 -DUSE_CLIPBOARD_AURAX11=1 -DENABLE_ONE_CLICK_SIGNIN
// #cgo CXXFLAGS: -DUSE_XI2_MT=2 -DENABLE_REMOTING=1 -DENABLE_WEBRTC=1
// #cgo CXXFLAGS: -DENABLE_PEPPER_CDMS -DENABLE_CONFIGURATION_POLICY
// #cgo CXXFLAGS: -DENABLE_INPUT_SPEECH -DENABLE_NOTIFICATIONS -DUSE_UDEV
// #cgo CXXFLAGS: -DENABLE_EGLIMAGE=1 -DENABLE_TASK_MANAGER=1 -DENABLE_EXTENSIONS=1
// #cgo CXXFLAGS: -DENABLE_PLUGIN_INSTALLATION=1 -DENABLE_PLUGINS=1
// #cgo CXXFLAGS: -DENABLE_SESSION_SERVICE=1 -DENABLE_THEMES=1
// #cgo CXXFLAGS: -DENABLE_AUTOFILL_DIALOG=1 -DENABLE_BACKGROUND=1
// #cgo CXXFLAGS: -DENABLE_AUTOMATION=1 -DENABLE_GOOGLE_NOW=1 -DCLD_VERSION=2
// #cgo CXXFLAGS: -DENABLE_FULL_PRINTING=1 -DENABLE_PRINTING=1 -DENABLE_SPELLCHECK=1
// #cgo CXXFLAGS: -DENABLE_CAPTIVE_PORTAL_DETECTION=1 -DENABLE_APP_LIST=1
// #cgo CXXFLAGS: -DENABLE_SETTINGS_APP=1 -DENABLE_MANAGED_USERS=1 -DENABLE_MDNS=1
// #cgo CXXFLAGS: -DLIBPEERCONNECTION_LIB=1 -DLINUX -DHAVE_SCTP
// #cgo CXXFLAGS: -DHASH_NAMESPACE=__gnu_cxx -DPOSIX -DDISABLE_DYNAMIC_CAST
// #cgo CXXFLAGS: -D_REENTRANT -DSSL_USE_NSS -DHAVE_NSS_SSL_H -DSSL_USE_NSS_RNG
// #cgo CXXFLAGS: -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0
// 
// #cgo CXXFLAGS: -pthread -std=c++11 -Wno-narrowing -Wno-write-strings
// #cgo CXXFLAGS: -Iwebrtc/trunk
// #cgo CXXFLAGS: -Iwebrtc/trunk/third_party
// #cgo CXXFLAGS: -Iwebrtc/trunk/third_party/webrtc
// #cgo CXXFLAGS: -Iwebrtc/trunk/webrtc
// #cgo CXXFLAGS: -Iwebrtc/trunk/net/third_party/nss/ssl
// #cgo CXXFLAGS: -Iwebrtc/trunk/third_party/jsoncpp/overrides/include
// #cgo CXXFLAGS: -Iwebrtc/trunk/third_party/jsoncpp/source/include
// 
// #cgo LDFLAGS: -lstdc++ -lm -lnss3 -lnssutil3 -lX11 -lXext -lcrypto -lplc4 -lnspr4 -lexpat -ldl
// #cgo LDFLAGS: -Wl,--start-group
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/jsoncpp/libjsoncpp.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/talk/libjingle_peerconnection.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/net/third_party/nss/libcrssl.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/talk/libjingle.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/talk/libjingle_media.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/libyuv.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libjpeg_turbo/libjpeg_turbo.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/usrsctp/libusrsctplib.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libvideo_capture_module.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libwebrtc_utility.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudio_coding_module.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libCNG.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/common_audio/libcommon_audio.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/system_wrappers/source/libsystem_wrappers.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/common_audio/libcommon_audio_sse2.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libG711.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libG722.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libiLBC.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libiSAC.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libiSACFix.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libPCM16B.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libNetEq.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libwebrtc_opus.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/opus/libopus.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libacm2.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libNetEq4.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libmedia_file.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libwebrtc_video_coding.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libwebrtc_i420.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/common_video/libcommon_video.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/video_coding/utility/libvideo_coding_utility.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/video_coding/codecs/vp8/libwebrtc_vp8.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libvpx/libvpx.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libvpx/libvpx_asm_offsets_vp8.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libvpx/libvpx_intrinsics_mmx.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libvpx/libvpx_intrinsics_sse2.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libvpx/libvpx_intrinsics_ssse3.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libvideo_render_module.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/video_engine/libvideo_engine_core.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/librtp_rtcp.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libpaced_sender.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libremote_bitrate_estimator.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/remote_bitrate_estimator/librbe_components.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libbitrate_controller.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libvideo_processing.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libvideo_processing_sse2.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/voice_engine/libvoice_engine.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudio_conference_mixer.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudio_processing.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudioproc_debug_proto.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/protobuf/libprotobuf_lite.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudio_processing_sse2.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/webrtc/modules/libaudio_device.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/talk/libjingle_sound.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/talk/libjingle_p2p.a
// #cgo LDFLAGS: webrtc/trunk/out/Release/obj/third_party/libsrtp/libsrtp.a
// #cgo LDFLAGS: -Wl,--end-group
// #include "conductor.h"
import "C"
import "log"

//export Debug
func Debug(s *C.char) {
	log.Println("c++:", C.GoString(s))
}

//export RegisterCandidate
func RegisterCandidate(sdp, mid *C.char, line C.int) {
	// TODO multiplex this crap
	log.Println("go: channeling candidate");
	candidate <- candidateMsg{
		Sdp:  C.GoString(sdp),
		Mid:  C.GoString(mid),
		Line: int(line),
	}
}

//export RegisterOffer
func RegisterOffer(sdp *C.char) {
	log.Println("go: channeling offer")
	offer <- C.GoString(sdp)
}

func MakePeerConnection() {
	log.Println("go: initialising peer")
	C.InitPeerConn()
}

func Answer(sdp string) {
	log.Println("go: processing answer")
	C.Answer(C.CString(sdp))
}

func Candidate(sdp, mid string, line int) {
	C.Candidate(C.CString(sdp), C.CString(mid), C.int(line))
}
