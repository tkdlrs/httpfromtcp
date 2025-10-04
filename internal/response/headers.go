package response

import (
	"fmt"

	"github.com/tkdlrs/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contenLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contenLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
