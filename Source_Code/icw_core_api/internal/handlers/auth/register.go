package auth

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	bizDto "icw_core_biz/pkg/dto"
)

// Register 注册
// @router /auth/register [POST]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.RegisterRequest{
		Email:     req.Email,
		EmailCode: req.EmailCode,
		Password:  req.Password,
		Name:      req.Name,
	}
	rpcResp := &bizDto.RegisterResponse{}
	if err := h.CoreBizCall("AuthService.Register", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewRegisterResponse(rpcResp))
}
