/*
 * Arc - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package utils

import (
	"github.com/evilsocket/arc/arcd/log"
	"github.com/gin-gonic/gin"
)

// Used when sending a HTTP status response different than 200
// swagger:response errorResponse
type ErrorResponse struct {
	// The error description.
	// in: body
	// Error code.
	Code int `json:"code"`
	// Error message.
	// in: body
	Message string `json:"message"`
}

func jError(level int, c *gin.Context, code int, message string) {
	log.Api(level, c, "[%d] %s", code, message)
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
	c.Abort()
}

// http://www.restapitutorial.com/httpstatuscodes.html
func NotFound(c *gin.Context, opt_msg ...string) {
	msg := "Not found."
	if len(opt_msg) > 0 {
		msg = opt_msg[0]
	}
	jError(log.WARNING, c, 404, msg)
}

func BadRequest(c *gin.Context, opt_msg ...string) {
	msg := "Bad request."
	if len(opt_msg) > 0 {
		msg = opt_msg[0]
	}
	jError(log.WARNING, c, 400, msg)
}

func Forbidden(c *gin.Context, opt_msg ...string) {
	msg := "Forbidden"
	if len(opt_msg) > 0 {
		msg = opt_msg[0]
	}
	jError(log.WARNING, c, 403, msg)
}

func ServerError(c *gin.Context, err error) {
	jError(log.ERROR, c, 500, "Error while performing operation: "+err.Error())
}
