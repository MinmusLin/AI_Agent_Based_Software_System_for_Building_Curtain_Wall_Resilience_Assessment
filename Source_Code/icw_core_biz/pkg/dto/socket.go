package dto

type SocketTicketContext struct {
	UserId      uint64
	ProjectId   uint64
	ProjectCode string
	SocketScope string
	RequestId   string
}

type CreateSocketTicketRequest struct {
	Meta        *Meta
	UserId      uint64
	ProjectId   uint64
	ProjectCode string
	SocketScope string
	RequestId   string
}

type CreateSocketTicketResponse struct {
	Ticket    string
	ExpiresIn int64
}

type ValidateSocketTicketRequest struct {
	Meta        *Meta
	ProjectCode string
	SocketScope string
	Ticket      string
}

type ValidateSocketTicketResponse struct {
	UserId      uint64
	ProjectId   uint64
	ProjectCode string
}
