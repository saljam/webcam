package main

// Using the Video4Linux docs and examples at http://linuxtv.org/downloads/v4l-dvb-apis/.

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"syscall"
	"unsafe"

	"launchpad.net/gommap"
)

// #cgo pkg-config: vpx
// #include <stdlib.h>
// #include <errno.h>
// #include <string.h>
// #include <sys/select.h>
// #include <linux/videodev2.h>
//
// #define VPX_CODEC_DISABLE_COMPAT 1
// #include "vpx/vpx_encoder.h"
// #include "vpx/vp8cx.h"
//
// struct vpxframe {
// 	void *   buf;
// 	size_t   sz;
// 	vpx_codec_pts_t   pts;
// 	unsigned long   duration;
// 	vpx_codec_frame_flags_t   flags;
// 	int   partition_id;
// };
//
// int selectfd(int fd)
// {
// 	struct timeval tout = { 2, 0 };
// 	fd_set fds;
//
// 	FD_ZERO(&fds);
// 	FD_SET(fd, &fds);
//
// 	return select(fd + 1, &fds, NULL, NULL, &tout);
// }
import "C"

func ioctl(fd int, request, argp uintptr) syscall.Errno {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	return errno
}

const (
	width  = 640
	height = 480
)

func streamVid(w io.Writer) {
	// Initialize encoder
	enc, err := initVP8()
	if err != nil {
		log.Println("disabling video: can't initialize encoder:", err)
		return
	}

	var img C.vpx_image_t
	C.vpx_img_alloc(&img, C.VPX_IMG_FMT_YV12, width, height, 1)

	// Initialize V4L
	fd, buffers, err := initV4L("/dev/video0")
	if err != nil {
		log.Println("disabling video: can't initialize video device:", err)
		return
	}
	defer syscall.Close(fd)

	writeIVFFileHeader(w, enc, 0)
	nframe := C.vpx_codec_pts_t(1)
	for {
		errno := syscall.Errno(C.selectfd(C.int(fd)))
		if errno.Temporary() {
			continue
		}
		if errno < 0 {
			log.Fatalf("select error: %v", errno)
		}

		// Replace these with ReadFrame/<buf>.Release.
		// Maybe no Release, deque prev frame on accessing next?
		buf := C.struct_v4l2_buffer{
			_type:  C.V4L2_BUF_TYPE_VIDEO_CAPTURE,
			memory: C.V4L2_MEMORY_MMAP,
		}
		errno = ioctl(fd, C.VIDIOC_DQBUF, uintptr(unsafe.Pointer(&buf)))
		if errno.Temporary() {
			continue
		}
		if errno != 0 {
			log.Fatalf("couldn't dequeue buffer: %v", err)
		}

		// Do something
		size := img.w * img.h * 2 // YV12
		pixels := (*[1 << 30]byte)(unsafe.Pointer(img.planes[0]))[:size:size]
		copy(pixels, buffers[buf.index][:buf.bytesused])

		encerr := C.vpx_codec_encode(enc, &img, nframe, 1, 0, C.VPX_DL_REALTIME)
		if encerr != 0 {
			log.Fatalf("failed to encode a frame: %s", C.GoString(C.vpx_codec_error_detail(enc)))
		}

		iter := C.vpx_codec_iter_t(nil)
		for pkt := C.vpx_codec_get_cx_data(enc, &iter); pkt != nil; pkt = C.vpx_codec_get_cx_data(enc, &iter) {
			frame := (*C.struct_vpxframe)(unsafe.Pointer(&pkt.data))
			cx := (*[1 << 30]byte)(unsafe.Pointer(frame.buf))[:frame.sz:frame.sz]
			switch pkt.kind {
			case C.VPX_CODEC_CX_FRAME_PKT:
				writeIVFFrameHeader(w, frame)
				w.Write(cx)
			default:
			}
			if pkt.kind == C.VPX_CODEC_CX_FRAME_PKT && (frame.flags & C.VPX_FRAME_IS_KEY != 0) {
				fmt.Printf("K")
			} else {
				fmt.Printf(".")
			}
		}

		nframe++

		errno = ioctl(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&buf)))
		if errno != 0 {
			log.Fatalf("couldn't queue buffer %d: %v", buf.index, err)
		}
	}
}

const fourcc = 0x30385056

func writeIVFFileHeader(w io.Writer, enc *C.vpx_codec_ctx_t, nframe int) {
	cfg := (*C.vpx_codec_enc_cfg_t)(unsafe.Pointer(&enc.config))

	w.Write([]byte{'D', 'K', 'I', 'F'})
	binary.Write(w, binary.LittleEndian, uint16(0))                  // Version
	binary.Write(w, binary.LittleEndian, uint16(32))                 // Header size
	binary.Write(w, binary.LittleEndian, uint32(fourcc))             // Fourcc
	binary.Write(w, binary.LittleEndian, uint16(cfg.g_w))            // Width
	binary.Write(w, binary.LittleEndian, uint16(cfg.g_h))            // Height
	binary.Write(w, binary.LittleEndian, uint32(cfg.g_timebase.den)) // Rate
	binary.Write(w, binary.LittleEndian, uint32(cfg.g_timebase.num)) // Scale
	binary.Write(w, binary.LittleEndian, uint32(nframe))             // Length
	binary.Write(w, binary.LittleEndian, uint32(0))                  // Unused
}

func writeIVFFrameHeader(w io.Writer, frame *C.struct_vpxframe) {
	binary.Write(w, binary.LittleEndian, uint32(frame.sz))
	binary.Write(w, binary.LittleEndian, uint32(frame.pts&0xFFFFFFFF))
	binary.Write(w, binary.LittleEndian, uint32(frame.pts>>32))
}

// This crap should also go into its own vp8 interface package at some point.
func initVP8() (*C.vpx_codec_ctx_t, error) {
	var cfg C.vpx_codec_enc_cfg_t
	errno := C.vpx_codec_enc_config_default(C.vpx_codec_vp8_cx(), &cfg, 0)
	if errno != 0 {
		return nil, fmt.Errorf(C.GoString(C.vpx_codec_err_to_string(errno)))
	}

	cfg.rc_target_bitrate = width * height * cfg.rc_target_bitrate / cfg.g_w / cfg.g_h
	cfg.g_w = width
	cfg.g_h = height

	var enc C.vpx_codec_ctx_t
	errno = C.vpx_codec_enc_init_ver(&enc, C.vpx_codec_vp8_cx(), &cfg, 0, C.VPX_ENCODER_ABI_VERSION)
	if errno != 0 {
		return nil, fmt.Errorf(C.GoString(C.vpx_codec_err_to_string(errno)))
	}

	return &enc, nil
}

// This crap should go into its own v4l interface package at some point.
func initV4L(path string) (fd int, buffers []gommap.MMap, err error) {
	fd, err = syscall.Open(path, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
	if err != nil {
		return fd, buffers, err
	}

	// Check camera capabilities.
	var cap C.struct_v4l2_capability
	errno := ioctl(fd, C.VIDIOC_QUERYCAP, uintptr(unsafe.Pointer(&cap)))
	if errno != 0 {
		return fd, buffers, errno
	}
	if (cap.capabilities & C.V4L2_CAP_VIDEO_CAPTURE) == 0 {
		return fd, buffers, fmt.Errorf("%s does not support video capture", path)
	}
	if (cap.capabilities & C.V4L2_CAP_STREAMING) == 0 {
		return fd, buffers, fmt.Errorf("%s does not support streaming io", path)
	}

	// TODO allow options here
	vfmt := C.struct_v4l2_format{
		_type: C.V4L2_BUF_TYPE_VIDEO_CAPTURE,
	}
	pix := (*C.struct_v4l2_pix_format)(unsafe.Pointer(&vfmt.fmt))
	pix.width = width
	pix.height = height
	pix.pixelformat = C.V4L2_PIX_FMT_YVU420
	pix.field = C.V4L2_FIELD_INTERLACED
	errno = ioctl(fd, C.VIDIOC_S_FMT, uintptr(unsafe.Pointer(&vfmt)))
	if errno != 0 {
		return fd, buffers, errno
	}

	// Sort out buffers.
	req := C.struct_v4l2_requestbuffers{
		count:  4,
		_type:  C.V4L2_BUF_TYPE_VIDEO_CAPTURE,
		memory: C.V4L2_MEMORY_MMAP,
	}

	errno = ioctl(fd, C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		return fd, buffers, errno
	}
	if req.count < 2 {
		return fd, buffers, fmt.Errorf("not enough buffers: %v", err)
	}

	buffers = make([]gommap.MMap, req.count)

	for i := 0; i < len(buffers); i++ {
		buf := C.struct_v4l2_buffer{
			_type:  C.V4L2_BUF_TYPE_VIDEO_CAPTURE,
			memory: C.V4L2_MEMORY_MMAP,
			index:  C.__u32(i),
		}
		offset := (*C.__u32)(unsafe.Pointer(&buf.m))

		errno = ioctl(fd, C.VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(&buf)))
		if errno != 0 {
			return fd, buffers, errno
		}
		buffers[i], err = gommap.MapRegion(uintptr(fd),
			int64(*offset), int64(buf.length),
			gommap.PROT_READ|gommap.PROT_WRITE,
			gommap.MAP_SHARED)
		if err != nil {
			return fd, buffers, err
		}
	}

	// Start capture
	for i := 0; i < len(buffers); i++ {
		buf := C.struct_v4l2_buffer{
			_type:  C.V4L2_BUF_TYPE_VIDEO_CAPTURE,
			memory: C.V4L2_MEMORY_MMAP,
			index:  C.__u32(i),
		}
		errno = ioctl(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&buf)))
		if errno != 0 {
			return fd, buffers, errno
		}
	}

	typ := C.V4L2_BUF_TYPE_VIDEO_CAPTURE
	errno = ioctl(fd, C.VIDIOC_STREAMON, uintptr(unsafe.Pointer(&typ)))
	if errno != 0 {
		return fd, buffers, errno
	}

	return
}
