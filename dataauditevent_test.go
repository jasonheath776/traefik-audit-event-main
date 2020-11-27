package traefik_audit_event

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTPp(t *testing.T) {
	tests := []struct {
		desc string
		req   string
		res  string
		expEvent   string
		expStatusCode int
	}{
		{
			desc:          "generate event payload",
			req:           "/jason",
			expEvent:      "{\"req\":\"/jason\",\"res\":\"\"}\"",
			expStatusCode: http.StatusOK,
		}, {
			desc:          "generate event payload",
			req:           "/jason2",
			expEvent:      "{\"req\":\"/jason2\",\"res\":\"\"}\"",
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			cfg := &AuditEvent{
				Req: test.req,
				Res: test.res,
				next:        next,
			}


			assert.Equal(t, GenerateEventPayload(cfg).String(), test.expEvent, "Both events should be equal.")

		})
	}
}
