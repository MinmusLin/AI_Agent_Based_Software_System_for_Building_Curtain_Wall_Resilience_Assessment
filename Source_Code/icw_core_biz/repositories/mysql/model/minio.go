package model

// ProjectImageObjectReference 项目图像对象引用
type ProjectImageObjectReference struct {
	ProjectId uint64
	ImageUuid string
}

// ProjectDetectionObjectReference 项目检测产物对象引用
type ProjectDetectionObjectReference struct {
	ProjectId uint64
	ImageUuid string
	TaskCode  string
}

// ProjectReportObjectReference 项目报告对象引用
type ProjectReportObjectReference struct {
	ProjectId uint64
}
