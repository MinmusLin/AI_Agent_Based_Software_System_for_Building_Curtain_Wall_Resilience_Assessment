package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
	"icw_core_api/utils"
)

// Register 注册
// @router /auth/register [POST]
func (h *Handler) Register(c *gin.Context) {
	var req apipb.RegisterRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizpb.RegisterRequest{
		Email:     req.Email,
		EmailCode: req.EmailCode,
		Password:  req.Password,
		Name:      req.Name,
	}
	rpcResp := &bizpb.RegisterResponse{}
	if err := auth.Register(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, utils.NewRegisterResponse(rpcResp))
}
