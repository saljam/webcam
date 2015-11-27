#include <owr/owr.h>
#include <owr/owr_types.h>
#include <owr/owr_session.h>
#include <owr/owr_transport_agent.h>
#include <owr/owr_media_session.h>
#include <owr/owr_local_media_source.h>

#include "_cgo_export.h"

// The example code in openwebrtc/tests/test_client.c was very useful for writing this.

OwrTransportAgent *transport_agent;

void init() {
	owr_init(NULL);
	// Should this be per-session?
	transport_agent = owr_transport_agent_new(FALSE);
	owr_transport_agent_add_helper_server(transport_agent, OWR_HELPER_SERVER_TYPE_STUN,
		"stun.services.mozilla.com", 3478, NULL, NULL);
	owr_run_in_background();
}

// void got_local_sources(GList *sources, OwrMediaSession *session) {}

static gchar *candidate_types[] = { "host", "srflx", "relay", NULL };
static gchar *tcp_types[] = { "", "active", "passive", "so", NULL };

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
		"ufrag", &ufrag,
		"password", &password,
		"type", &type,
		"component-type", &component,
		"foundation", &foundation,
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

OwrSession* new_session() {
	OwrMediaSession *session = owr_media_session_new(0);
	
	// Not using this since we're using auto for now.
	//owr_get_capture_sources(OWR_MEDIA_TYPE_AUDIO | OWR_MEDIA_TYPE_VIDEO,
	//	(OwrCaptureSourcesCallback)got_local_sources, session);

	owr_media_session_set_send_source(session, NULL);
	g_signal_connect(session, "on-new-candidate",
		G_CALLBACK(got_candidate), NULL);
	g_signal_connect(session, "on-candidate-gathering-done",
		G_CALLBACK(candidate_gathering_done_go), NULL);

	owr_transport_agent_add_session(transport_agent, OWR_SESSION(session));
	
	return OWR_SESSION(session);
}

void del_session(OwrSession* session, OwrTransportAgent *transport) {
	// transport, in theory, owns session and will do the correct thing.
	g_object_unref(transport);
}
