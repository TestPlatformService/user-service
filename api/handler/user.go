package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "user/genproto/user"
)

// Register godoc
// @Summary Register user
// @Description create new users
// @Tags user
// @Param info body user.RegisterReq true "User info"
// @Success 200 {object} user.RegisterRes
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error"
// @Router /user/register [post]
func (h Handler) Register(c *gin.Context) {
	h.Log.Info("Register is starting")
	req := pb.RegisterRequest{}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.User.Register(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.Log.Info("Register ended")
	c.JSON(http.StatusOK, res)
}

