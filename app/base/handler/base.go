package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var statusCode = 200

type ErrorResponse struct {
	ResponseCode        string `json:"response_code"`
	ResponseDescription string `json:"response_description"`
	ResponseData        string `json:"response_data"`
}

// SetStatusCode -- Set status code
func SetStatusCode(statCode int) int {
	statusCode := statCode
	return statusCode
}

// RespondJSON -- set response to json format
func RespondJSON(c *gin.Context, data interface{}, request interface{}) {
	c.JSON(statusCode, data)
	return
}

// RespondSuccess ..
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"response_code":        "00",
		"response_description": "Success",
		"response_data":        data,
	})
	return
}

func FailedAPIValidation(message string, c *gin.Context) {

	var res = ErrorResponse{
		"03",
		"Failed Validation BackEnd",
		message,
	}

	c.JSON(400, res)
}

func FailedConnectionBackend() {

}

// FailedResponseBackend ..
func FailedResponseBackend(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"response_code":        "03",
		"response_description": "Failed Validation BackEnd",
		"response_data":        err.Error()})
	return
}

func GeneralError(message string, c *gin.Context) {
	var res = ErrorResponse{
		"03",
		"Failed Validation BackEnd",
		message,
	}

	c.JSON(500, res)
}

// RespondCreated -- Set response for create process
func RespondCreated(c *gin.Context, message string, data interface{}, request interface{}) {
	if message == "" {
		message = "Resource Created"
	}
	statusCode := SetStatusCode(201)
	c.JSON(statusCode, data)
	return
}

// RespondUpdated -- Set response for update process
func RespondUpdated(c *gin.Context, message string) {
	if message == "" {
		message = "Resource Updated"
	}
	statusCode := SetStatusCode(200)
	c.JSON(statusCode, gin.H{"message": message})
	return
}

// RespondDeleted -- Set response for delete process
func RespondDeleted(c *gin.Context, message string) {
	if message == "" {
		message = "Resource Deleted"
	}
	statusCode := SetStatusCode(200)
	c.JSON(statusCode, gin.H{"message": message})
	return
}

// RespondError -- Set response for error
func RespondError(c *gin.Context, message interface{}, statusCode int, request interface{}) {
	data := gin.H{"status": statusCode, "message": message}
	c.JSON(statusCode, data)
	return
}

// RespondFailValidation -- Set response for fail validation
/* func RespondFailValidationNew(c *gin.Context, code string, status string, errorMessage interface{}, request interface{}) {
	statusCode := SetStatusCode(422)
	data := gin.H{"code": code, "status": status, "error": errorMessage}
	c.JSON(statusCode, data)
	return
} */

// RespondFailValidation ..
func RespondFailValidation(c *gin.Context, message interface{}, request interface{}) {
	RespondError(c, message, 422, request)
	return
}

// RespondUnauthorized -- Set response not authorized
func RespondUnauthorized(c *gin.Context, message string) {

	c.JSON(http.StatusUnauthorized, gin.H{
		"response_code":        "401",
		"response_description": "Unauthorized",
		"response_data":        message,
	})
	c.Abort()
	return
}

// RespondNotFound -- Set response not found
func RespondNotFound(c *gin.Context, code string, status string, message string, request interface{}) {
	if message == "" {
		message = "Resource Not Found"
	}
	statusCode := SetStatusCode(404)
	data := gin.H{"code": code, "status": status, "message": message}
	c.JSON(statusCode, data)
	return
}

// RespondMethodNotAllowed -- Set response method not allowed
func RespondMethodNotAllowed(c *gin.Context, message string) {
	if message == "" {
		message = "Method Not Allowed"
	}
	statusCode := SetStatusCode(405)
	c.JSON(statusCode, gin.H{"errors": message})
	return
}

func RespondForbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"response_code":        "403",
		"response_description": "Unauthorized",
		"reponse_data":         "forbidden",
	})
	return
}

// RespondError
/* func RespondErrorAPINew(c *gin.Context, code string, status string, message string, request interface{}) {
	statusCode := SetStatusCode(400)
	data := gin.H{"code": code, "status": status, "message": message}
	c.JSON(statusCode, data)
	return
} */

// RespondErrorAPI ..
func RespondErrorAPI(c *gin.Context, message string, request interface{}) {
	statusCode := SetStatusCode(400)
	data := gin.H{"status": statusCode, "message": message}
	c.JSON(statusCode, data)
	return
}
