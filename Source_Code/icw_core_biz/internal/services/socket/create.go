package socket

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
	"icw_common/rpc/error"
	"icw_core_biz/internal/services/socket/utils"
)

// CreateSocketTicket 创建 WebSocket 连接票据
func (s *Service) CreateSocketTicket(ctx context.Context, req *bizpb.CreateSocketTicketRequest) (*bizpb.CreateSocketTicketResponse, error) {
	resp := &bizpb.CreateSocketTicketResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.createSocketTicket(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) createSocketTicket(ctx context.Context, req *bizpb.CreateSocketTicketRequest, resp *bizpb.CreateSocketTicketResponse) error {
	if req.UserId == 0 || req.ProjectId == 0 {
		return rpc_error.BadRequestDefault("socket ticket resource is required")
	}
	projectCode := strings.TrimSpace(req.ProjectCode)
	socketScope := strings.TrimSpace(req.SocketScope)
	if projectCode == "" || socketScope == "" {
		return rpc_error.BadRequestDefault("socket ticket resource is required")
	}

	// 生成 WebSocket 连接票据
	ticket, err := utils.NewTicket()
	if err != nil {
		return err
	}

	// 生成 WebSocket 连接票据上下文
	expiresIn := int64(s.Config().SocketTicketTTL.Seconds())
	ticketContext := &utils.SocketTicketContext{
		UserId:      req.UserId,
		ProjectId:   req.ProjectId,
		ProjectCode: projectCode,
		SocketScope: socketScope,
		RequestId:   rpc.RequestIdFromIncomingContext(ctx),
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
