package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Register .
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
	if err := h.rpc.Call("AuthService.Register", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("Call icw.core.biz AuthService.Register failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteRPCError(c, err)
		return
	}

	response.OK(c, dto.NewRegisterResponse(rpcResp))
}
