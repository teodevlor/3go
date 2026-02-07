package common

const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusAccepted            = 202
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

var httpMessage = map[int]string{
	StatusOK:                  "Success",
	StatusCreated:             "Created",
	StatusAccepted:            "Accepted",
	StatusNoContent:           "No Content",
	StatusBadRequest:          "Bad Request",
	StatusUnprocessableEntity: "Unprocessable Entity",
	StatusTooManyRequests:     "Too Many Requests",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
	StatusBadGateway:          "Bad Gateway",
	StatusServiceUnavailable:  "Service Unavailable",
	StatusGatewayTimeout:      "Gateway Timeout",
}

func HTTPMessage(code int) string {
	if msg, ok := httpMessage[code]; ok {
		return msg
	}
	return ""
}
