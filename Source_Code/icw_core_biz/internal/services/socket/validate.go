package socket

import (
	"context"
	"encoding/json"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/internal/services/socket/utils"
)

// ValidateSocketTicket 校验 WebSocket 连接票据
func (s *Service) ValidateSocketTicket(ctx context.Context, req *bizpb.ValidateSocketTicketRequest) (*bizpb.ValidateSocketTicketResponse, error) {
	resp := &bizpb.ValidateSocketTicketResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.validateSocketTicket(req, resp)
	})
	return resp, err
}

func (s *Service) validateSocketTicket(req *bizpb.ValidateSocketTicketRequest, _ *bizpb.ValidateSocketTicketResponse) error {
	ticket := strings.TrimSpace(req.Ticket)
	projectCode := strings.TrimSpace(req.ProjectCode)
	socketScope := strings.TrimSpace(req.SocketScope)
	if ticket == "" || projectCode == "" || socketScope == "" {
		return rpc_err.BadRequestDefault("socket ticket resource is required")
	}

	// 消费 WebSocket 连接票据上下文
	rawContext, err := s.Redis().ConsumeSocketTicket(s.Ctx(), utils.TicketHash(ticket))
	if err != nil {
		return err
	}
	if rawContext == "" {
		return rpc_err.BadRequestDefault("socket ticket is invalid")
	}

	// 解析 WebSocket 连接票据上下文
	var ticketContext utils.SocketTicketContext
	if err := json.Unmarshal([]byte(rawContext), &ticketContext); err != nil {
		return err
	}
	if ticketContext.ProjectCode != projectCode || ticketContext.SocketScope != socketScope {
		return rpc_err.BadRequestDefault("socket ticket is mismatched")
	}

	return nil
}
