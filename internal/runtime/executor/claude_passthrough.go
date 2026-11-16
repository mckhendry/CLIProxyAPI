package executor

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// claudePassthroughSentinel, set as a claude-api-key entry's api-key, makes the
// proxy forward the caller's own Authorization bearer upstream (FORK).
const claudePassthroughSentinel = "passthrough-client-credential"

// claudeIncomingBearer returns the caller's bearer token from the incoming
// request headers, or from the gin request when none were handed in.
func claudeIncomingBearer(r *http.Request, incoming http.Header) string {
	h := incoming
	if h == nil && r != nil {
		if g, ok := r.Context().Value("gin").(*gin.Context); ok && g != nil && g.Request != nil {
			h = g.Request.Header
		}
	}
	if h == nil {
		return ""
	}
	v := strings.TrimSpace(h.Get("Authorization"))
	if len(v) > 7 && strings.EqualFold(v[:7], "bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return strings.TrimSpace(h.Get("X-Api-Key"))
}
