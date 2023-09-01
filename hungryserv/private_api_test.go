package hungryserv

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/matrix-org/complement/internal/b"
	"github.com/matrix-org/complement/internal/client"
	"github.com/matrix-org/complement/internal/match"
	"github.com/matrix-org/complement/internal/must"
)

// Test all private APIs require auth
func TestPrivateAPI(t *testing.T) {
	deployment := Deploy(t, b.BlueprintCleanHS)
	defer deployment.Destroy(t)

	endpoints := []struct {
		Path   string
		Method string
	}{
		{Path: "/_matrix/client/v3/capabilities", Method: http.MethodGet},

		// Devices endpoints
		{Path: "/_matrix/client/v3/sendToDevice/{type}/{txn}", Method: http.MethodPut},

		{Path: "/_matrix/client/v3/account/whoami", Method: http.MethodGet},

		// Account data endpoints
		{Path: "/_matrix/client/v3/user/{userID}/account_data/{type}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/user/{userID}/account_data/{type}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/user/{userID}/rooms/{roomID}/account_data/{type}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/user/{userID}/rooms/{roomID}/account_data/{type}", Method: http.MethodPut},
		{Path: "/_matrix/client/unstable/com.beeper.inbox/user/{userID}/rooms/{roomID}/inbox_state", Method: http.MethodPut},

		// Profile endpoints
		{Path: "/_matrix/client/v3/profile/{userID}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/profile/{userID}/{field}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/profile/{userID}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/profile/{userID}", Method: http.MethodPatch},
		{Path: "/_matrix/client/v3/profile/{userID}/{field}", Method: http.MethodPut},

		// Keys endpoints
		{Path: "/_matrix/client/v3/keys/claim", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/keys/query", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/keys/upload", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/keys/changes", Method: http.MethodGet},

		// Presence endpoints
		{Path: "/_matrix/client/v3/presence/{userID}/status", Method: http.MethodPut},

		// User endpoints
		{Path: "/_matrix/client/v3/user/{userID}/filter", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/user/{userID}/filter/{filterID}", Method: http.MethodGet},

		// Sync endpoints
		{Path: "/_matrix/client/v3/sync", Method: http.MethodGet},

		// Room endpoints
		{Path: "/_matrix/client/v3/createRoom", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/join/{roomIDOrAlias}", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/knock/{roomIDOrAlias}", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/joined_rooms", Method: http.MethodGet},

		{Path: "/_matrix/client/v3/rooms/{roomID}/invite", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/join", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/leave", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/forget", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/kick", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/ban", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/unban", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/joined_members", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/messages", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/members", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/read_markers", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/receipt/{type}/{eventID}", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/redact/{eventID}/{txn}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/rooms/{roomID}/send/{type}/{txn}", Method: http.MethodPut},
		{Path: "/_matrix/client/unstable/org.matrix.msc2716/rooms/{roomID}/batch_send", Method: http.MethodPost},
		{Path: "/_matrix/client/unstable/com.beeper.chatmerging/rooms/{roomID}/split", Method: http.MethodPost},
		{Path: "/_matrix/client/unstable/com.beeper.chatmerging/merge", Method: http.MethodPost},
		{Path: "/_matrix/client/unstable/com.beeper.yeet/rooms/{roomID}/delete", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}/", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}/{stateKey}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}/", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/rooms/{roomID}/state/{type}/{stateKey}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/rooms/{roomID}/event/{eventID}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/context/{eventID}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/rooms/{roomID}/typing/{userID}", Method: http.MethodPut},

		{Path: "/_matrix/client/v3/user/{userID}/rooms/{roomID}/tags", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/user/{userID}/rooms/{roomID}/tags/{tag}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/user/{userID}/rooms/{roomID}/tags/{tag}", Method: http.MethodDelete},

		{Path: "/_matrix/client/v3/directory/room/{roomAlias}", Method: http.MethodGet},
		{Path: "/_matrix/client/v3/directory/room/{roomAlias}", Method: http.MethodPut},
		{Path: "/_matrix/client/v3/directory/room/{roomAlias}", Method: http.MethodDelete},
		{Path: "/_matrix/client/v3/rooms/{roomID}/aliases", Method: http.MethodGet},

		{Path: "/_matrix/client/v3/logout", Method: http.MethodPost},
		{Path: "/_matrix/client/v3/logout/all", Method: http.MethodPost},

		{Path: "/_matrix/client/v3/thirdparty/protocols", Method: http.MethodGet},

		// Private media endpoints
		{Path: "/_matrix/media/v3/config", Method: http.MethodGet},
		{Path: "/_matrix/media/v3/upload", Method: http.MethodPost},
		{Path: "/_matrix/media/v3/preview_url", Method: http.MethodGet},

		// Private MSC2246 media endpoints
		{Path: "/_matrix/media/unstable/fi.mau.msc2246/create", Method: http.MethodPost},
		{Path: "/_matrix/media/unstable/fi.mau.msc2246/upload/{serverName}/{mediaID}", Method: http.MethodPut},
		{Path: "/_matrix/media/unstable/fi.mau.msc2246/upload/{serverName}/{mediaID}/complete", Method: http.MethodPost},

		// Roomserv Synapse change notify endpoints
		{Path: "/_matrix/hungryserv/unstable/refresh_devices", Method: http.MethodPost},
		{Path: "/_matrix/hungryserv/unstable/devices", Method: http.MethodPost},
		{Path: "/_matrix/hungryserv/unstable/pushers", Method: http.MethodPost},
		{Path: "/_matrix/hungryserv/unstable/push_rules", Method: http.MethodPost},
		{Path: "/_matrix/hungryserv/unstable/invalidate_user_account_data", Method: http.MethodPost},

		// Appservice websockets
		{Path: "/_matrix/client/unstable/fi.mau.as_sync", Method: http.MethodGet},

		// Asmux custom registration endpoint
		{Path: "/_matrix/asmux/mxauth/appservice/{username}/{bridge}", Method: http.MethodPut},
		{Path: "/_matrix/asmux/mxauth/appservice/{username}/{bridge}", Method: http.MethodGet},
		{Path: "/_matrix/asmux/mxauth/appservice/{username}/{bridge}", Method: http.MethodDelete},

		// Appservice exec endpoints
		{Path: "/_matrix/asmux/mxauth/appservice/{owner}/{prefix}/exec/{command}", Method: http.MethodPost},
	}

	anon := deployment.Client(t, "hs1", "")
	unkn := deployment.Client(t, "hs1", "")
	unkn.AccessToken = "invalid"

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("%s %s must require auth", endpoint.Method, endpoint.Path), func(t *testing.T) {
			path := strings.Split(endpoint.Path, "/")[1:]
			var resp *http.Response
			if endpoint.Method == http.MethodPost || endpoint.Method == http.MethodPut || endpoint.Method == http.MethodPatch {
				resp = anon.DoFunc(t, endpoint.Method, path, client.WithJSONBody(t, map[string]interface{}{}))
			} else {
				resp = anon.DoFunc(t, endpoint.Method, path)
			}

			must.MatchResponse(t, resp, match.HTTPResponse{
				StatusCode: 401,
			})
		})

		t.Run(fmt.Sprintf("%s %s must not work with wrong access token", endpoint.Method, endpoint.Path), func(t *testing.T) {
			path := strings.Split(endpoint.Path, "/")[1:]
			var resp *http.Response
			if endpoint.Method == http.MethodPost || endpoint.Method == http.MethodPut || endpoint.Method == http.MethodPatch {
				resp = unkn.DoFunc(t, endpoint.Method, path, client.WithJSONBody(t, map[string]interface{}{}))
			} else {
				resp = unkn.DoFunc(t, endpoint.Method, path)
			}

			must.MatchResponse(t, resp, match.HTTPResponse{
				StatusCode: 401,
			})
		})
	}
}
