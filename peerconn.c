#include <owr/owr.h>
#include <owr/owr_types.h>
#include <owr/owr_session.h>
#include <owr/owr_media_session.h>
#include <owr/owr_local_media_source.h>

OwrMediaSession* new_session() {
	OwrMediaSession *session = owr_media_session_new(0);
	GObject *src = g_object_new(owr_local_media_source_get_type(), NULL);

	GValue val = G_VALUE_INIT;
	g_value_init(&val, G_TYPE_INT);
	g_value_set_int(&val, -1);

	g_object_set_property(src, "device-index", &val);
	owr_media_session_set_send_source(session, (OwrMediaSource*)src);
	//session.on-new-candidate = cbCandidate;
	return session;
}
