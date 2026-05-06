package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewCreateSocketTicketResponse(resp *bizpb.CreateSocketTicketResponse) *apipb.CreateSocketTicketResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateSocketTicketResponse{
		Ticket:    resp.Ticket,
		ExpiresIn: resp.ExpiresIn,
	}
}
