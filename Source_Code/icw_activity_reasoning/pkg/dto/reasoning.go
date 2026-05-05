package dto

type ArtifactUploadPlan struct {
	Type             string `json:"type"`
	PresignUploadURL string `json:"presign_upload_url"`
	ContentType      string `json:"content_type"`
}

type StartRequest struct {
	TaskCode      string                `json:"task_code"`
	ImageUuid     string                `json:"image_uuid"`
	PresignGetURL string                `json:"presign_get_url"`
	Artifacts     []*ArtifactUploadPlan `json:"artifacts"`
}

type StartResponse struct {
	Accepted bool `json:"accepted"`
}
