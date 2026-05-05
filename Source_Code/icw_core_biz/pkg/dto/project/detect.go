package project

import (
	"icw_core_biz/pkg/dto"
)

type ReasoningArtifact struct {
	Type     string
	Uploaded bool
	Sha256   string
}

type ReportClassificationResultRequest struct {
	Meta         *dto.Meta
	TaskUuid     string
	ImageUuid    string
	Status       string
	TaskCodes    []string
	ErrorMessage string
}

type ReportClassificationResultResponse struct{}

type ReportReasoningResultRequest struct {
	Meta         *dto.Meta
	TaskCode     string
	TaskUuid     string
	ImageUuid    string
	Status       string
	ResultJSON   string
	Artifacts    []*ReasoningArtifact
	ErrorMessage string
}

type ReportReasoningResultResponse struct{}

type ReportSummaryResultRequest struct {
	Meta         *dto.Meta
	TaskUuid     string
	ImageUuid    string
	Status       string
	ResultJSON   string
	ErrorMessage string
}

type ReportSummaryResultResponse struct{}
