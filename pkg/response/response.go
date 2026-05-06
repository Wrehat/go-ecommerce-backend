package response

import "github.com/gin-gonic/gin"

// Nampan standar untuk semua API di sistem kita
type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Fungsi global yang bisa dipanggil oleh Handler mana pun
func JSON(c *gin.Context, status int, message string, data any) {
	c.JSON(status, APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
