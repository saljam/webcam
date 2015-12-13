#include <stdlib.h>
#include <string.h>

#include <owr/owr.h>
#include <owr/owr_types.h>
#include <owr/owr_session.h>
#include <owr/owr_transport_agent.h>
#include <owr/owr_media_session.h>
#include <owr/owr_local_media_source.h>
#include <owr/owr_local.h>
#include <owr/owr_video_payload.h>

#include "_cgo_export.h"

// The example code in openwebrtc/tests/test_client.c was very useful for writing this.

OwrTransportAgent *transport_agent;
GList *local_sources;

void got_local_sources(GList *sources, gpointer user_data) {
	local_sources = g_list_copy(sources);
	// Should this be per-session?
	transport_agent = owr_transport_agent_new(FALSE);
	owr_transport_agent_add_helper_server(transport_agent, OWR_HELPER_SERVER_TYPE_STUN,
		"stun.services.mozilla.com", 3478, NULL, NULL);
	got_local_sources_go();
}

void init() {
	owr_init(NULL);
	owr_run_in_background();
	owr_get_capture_sources(OWR_MEDIA_TYPE_AUDIO | OWR_MEDIA_TYPE_VIDEO,
		(OwrCaptureSourcesCallback)got_local_sources, NULL);
}

void add_candidate(OwrSession *session, char *ufrag, char *password, int type, int component, char *foundation, uint priority, int transport_type, char *address, uint port) {
	OwrCandidate *candidate = owr_candidate_new(type, component);
	g_object_set(candidate,
		"foundation", foundation,
		"ufrag", ufrag,
		"password", password,
		"transport-type", transport_type,
		"address", address,
		"port", port,
		"priority", priority,
		NULL);
	owr_session_add_remote_candidate(session, candidate);
}

void got_candidate(GObject *session, OwrCandidate *candidate, gpointer user_data) {
	gchar *ufrag, *password;
	int type;
	int component;
	gchar *foundation;
	guint priority;
	int transport_type;
	gchar *address, *base_address;
	guint port, base_port;

	g_object_get(candidate,
		"foundation", &foundation,
		"ufrag", &ufrag,
		"password", &password,
		"type", &type,
		"component-type", &component,
		"priority", &priority,
		"transport-type", &transport_type,
		"address", &address, "port", &port,
		"base-address", &base_address, "base-port", &base_port,
		NULL);

	got_candidate_go(OWR_SESSION(session),
		ufrag, password, type, component, foundation, priority,
		transport_type, port, base_port, address, base_address);

	g_free(ufrag);
	g_free(password);
	g_free(foundation);
	g_free(address);
	g_free(base_address);
}

void got_dtls_certificate(GObject *media_session, GParamSpec *pspec, gpointer user_data)
{
	guint i;
	gchar *pem, *line;
	guchar *der, *tmp;
	gchar **lines;
	gint state = 0;
	guint save = 0;
	gsize der_length = 0;
	GChecksum *checksum;
	guint8 *digest;
	gsize digest_length;
	GString *fingerprint;

	g_object_get(media_session, "dtls-certificate", &pem, NULL);
	der = tmp = g_new0(guchar, (strlen(pem) / 4) * 3 + 3);
	lines = g_strsplit(pem, "\n", 0);
	for (i = 0, line = lines[i]; line; line = lines[++i]) {
		if (line[0] && !g_str_has_prefix(line, "-----"))
			tmp += g_base64_decode_step(line, strlen(line), tmp, &state, &save);
	}
	der_length = tmp - der;
	checksum = g_checksum_new(G_CHECKSUM_SHA256);
	digest_length = g_checksum_type_get_length(G_CHECKSUM_SHA256);
	digest = g_new(guint8, digest_length);
	g_checksum_update(checksum, der, der_length);
	g_checksum_get_digest(checksum, digest, &digest_length);
	fingerprint = g_string_new(NULL);
	for (i = 0; i < digest_length; i++) {
		if (i)
			g_string_append(fingerprint, ":");
		g_string_append_printf(fingerprint, "%02X", digest[i]);
	}
	gchar *fprint = g_string_free(fingerprint, FALSE);
	got_dtls_certificate_go(OWR_SESSION(media_session), fprint);

	g_free(fprint);
	g_free(digest);
	g_checksum_free(checksum);
	g_free(der);
	g_strfreev(lines);
}

OwrSession* new_session() {
	Debug("new session");
	OwrMediaSession *session = owr_media_session_new(TRUE);
	g_signal_connect(session, "on-new-candidate",
		G_CALLBACK(got_candidate), NULL);
	g_signal_connect(session, "on-candidate-gathering-done",
		G_CALLBACK(candidate_gathering_done_go), NULL);
	g_signal_connect(session, "notify::dtls-certificate",
		G_CALLBACK(got_dtls_certificate), NULL);

	if (local_sources != NULL && local_sources->data != NULL) {
		Debug("setting local source");
		owr_media_session_set_send_source(session, OWR_MEDIA_SOURCE(local_sources->data));
	}
	// TODO make the specs (payload:120, encoding VP8, clock: 9000) mutable.
	OwrPayload *payload = owr_video_payload_new(OWR_CODEC_TYPE_VP8, 120, 90000, TRUE, TRUE);
	owr_media_session_set_send_payload(session, payload);

	// Add session to the transport.
	owr_transport_agent_add_session(transport_agent, OWR_SESSION(session));

	return OWR_SESSION(session);
}

void del_session(OwrSession* session) {
	// Free session.
}
