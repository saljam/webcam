build: webrtc
	python webrtc/build/gyp_webrtc
	ninja -C out/Release

webrtc: depot_tools
	mkdir webrtc
	export PATH=`pwd`/depot_tools:"$PATH" && cd webrtc && fetch --nohooks webrtc && gclient sync

depot_tools:
	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git

clean:
	rm -rf webrtc depot_tools
