package dto

type ArtifactUploadPlan struct {
	Name             string `json:"artifact_name"`
	PresignUploadURL string `json:"presign_upload_url"`
	ContentType      string `json:"content_type"`
}

type StartRequest struct {
	TaskUuid      string                `json:"task_uuid"`
	TaskCode      string                `json:"task_code"`
	ImageUuid     string                `json:"image_uuid"`
	PresignGetURL string                `json:"presign_get_url"`
	Artifacts     []*ArtifactUploadPlan `json:"artifacts"`
}

type StartResponse struct{}
