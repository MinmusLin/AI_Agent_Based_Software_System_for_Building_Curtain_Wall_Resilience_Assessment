package socket

import (
	"encoding/json"
	"strings"
	"time"

	"icw_core_biz/internal/services/socket/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// CreateSocketTicket 创建 WebSocket 连接票据
func (s *Service) CreateSocketTicket(req *dto.CreateSocketTicketRequest, resp *dto.CreateSocketTicketResponse) error {
	return s.CallRPC(req, resp, func() error {
		return s.createSocketTicket(req, resp)
	})
}

func (s *Service) createSocketTicket(req *dto.CreateSocketTicketRequest, resp *dto.CreateSocketTicketResponse) error {
	if req.UserId == 0 || req.ProjectId == 0 {
		return rpc_err.BadRequestDefault("socket ticket resource is required")
	}
	projectCode := strings.TrimSpace(req.ProjectCode)
	socketScope := strings.TrimSpace(req.SocketScope)
	if projectCode == "" || socketScope == "" {
		return rpc_err.BadRequestDefault("socket ticket resource is required")
	}

	// 生成 WebSocket 连接票据
	ticket, err := utils.NewTicket()
	if err != nil {
		return err
	}

	// 生成 WebSocket 连接票据上下文
	expiresIn := int64(s.Config().SocketTicketTTL.Seconds())
	ticketContext := &dto.SocketTicketContext{
		UserId:      req.UserId,
		ProjectId:   req.ProjectId,
		ProjectCode: projectCode,
		SocketScope: socketScope,
		RequestId:   strings.TrimSpace(req.RequestId),
		CreateAt:    time.Now().Format("2006-01-02 15:04:05"),
	}
	contextBytes, err := json.Marshal(ticketContext)
	if err != nil {
		return err
	}

	// 保存 WebSocket 连接票据上下文
	if err := s.Redis().SaveSocketTicket(s.Ctx(), utils.TicketHash(ticket), string(contextBytes), time.Duration(expiresIn)*time.Second); err != nil {
		return err
	}

	resp.Ticket = ticket
	resp.ExpiresIn = expiresIn

	return nil
}
