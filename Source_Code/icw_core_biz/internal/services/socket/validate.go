package socket

import (
	"encoding/json"
	"strings"

	"icw_core_biz/internal/services/socket/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// ValidateSocketTicket 校验 WebSocket 连接票据
func (s *Service) ValidateSocketTicket(req *dto.ValidateSocketTicketRequest, resp *dto.ValidateSocketTicketResponse) error {
	return s.CallRPC("SocketService.ValidateSocketTicket", req, resp, func() error {
		return s.validateSocketTicket(req, resp)
	})
}

func (s *Service) validateSocketTicket(req *dto.ValidateSocketTicketRequest, resp *dto.ValidateSocketTicketResponse) error {
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
	var ticketContext dto.SocketTicketContext
	if err := json.Unmarshal([]byte(rawContext), &ticketContext); err != nil {
		return err
	}
	if ticketContext.ProjectCode != projectCode || ticketContext.SocketScope != socketScope {
		return rpc_err.BadRequestDefault("socket ticket is mismatched")
	}

	resp.UserId = ticketContext.UserId
	resp.ProjectId = ticketContext.ProjectId
	resp.ProjectCode = ticketContext.ProjectCode

	return nil
}
