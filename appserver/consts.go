package appserver

import (
	"net/http"
)

const minPathLength = 3
const maxPathLength = 50

const operationalErrorCode = http.StatusServiceUnavailable // Throw 503 instead of 500
