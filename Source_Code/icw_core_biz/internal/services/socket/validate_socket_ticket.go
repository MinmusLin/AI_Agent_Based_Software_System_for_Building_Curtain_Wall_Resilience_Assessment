package socket

import (
	"context"
	"encoding/json"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/socket/utils"
)

// ValidateSocketTicket 校验 WebSocket 连接票据
func (s *Service) ValidateSocketTicket(ctx context.Context, req *bizpb.ValidateSocketTicketRequest) (*bizpb.ValidateSocketTicketResponse, error) {
	resp := &bizpb.ValidateSocketTicketResponse{}
	err := s.CallRPC(req, func() error {
		return s.validateSocketTicket(req, resp)
	})
	return resp, err
}

func (s *Service) validateSocketTicket(req *bizpb.ValidateSocketTicketRequest, _ *bizpb.ValidateSocketTicketResponse) error {
	ticket := strings.TrimSpace(req.Ticket)
	projectCode := strings.TrimSpace(req.ProjectCode)
	if ticket == "" || projectCode == "" || req.Scope == commonpb.SocketScope_Unknown {
		return rpc_error.BadRequestDefault("socket ticket resource is required")
	}

	// 校验是否是有效的 WebSocket 连接范围
	if !utils.ValidateSocketScope(req.Scope) {
		return rpc_error.BadRequestDefault("socket scope is invalid")
	}

	// 消费 WebSocket 连接票据上下文
	rawContext, err := s.Redis().ConsumeSocketTicket(s.Ctx(), utils.TicketHash(ticket))
	if err != nil {
		return err
	}
	if rawContext == "" {
		return rpc_error.BadRequestDefault("socket ticket is invalid")
	}

	// 解析 WebSocket 连接票据上下文
	var ticketContext utils.SocketTicketContext
	if err := json.Unmarshal([]byte(rawContext), &ticketContext); err != nil {
		return err
	}
	if ticketContext.ProjectCode != projectCode || ticketContext.SocketScope != req.Scope {
		return rpc_error.BadRequestDefault("socket ticket is mismatched")
	}

	return nil
}
