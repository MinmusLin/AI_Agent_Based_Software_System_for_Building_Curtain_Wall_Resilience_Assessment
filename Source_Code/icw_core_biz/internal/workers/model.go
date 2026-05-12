package workers

type projectReportSource struct {
	GroupCount int                         `json:"group_count,omitempty"`
	ImageCount int                         `json:"image_count,omitempty"`
	Project    projectReportSourceProject  `json:"project,omitempty"`
	Groups     []*projectReportSourceGroup `json:"groups,omitempty"`
}

type projectReportSourceProject struct {
	ProjectName         string `json:"project_name,omitempty"`
	BuildingName        string `json:"building_name,omitempty"`
	BuildingLocation    string `json:"building_location,omitempty"`
	BuiltYear           uint32 `json:"built_year,omitempty"`
	BuildingDescription string `json:"building_description,omitempty"`
	KnownIssues         string `json:"known_issues,omitempty"`
	AssessmentGoal      string `json:"assessment_goal,omitempty"`
}

type projectReportSourceGroup struct {
	GroupName string                      `json:"group_name,omitempty"`
	Images    []*projectReportSourceImage `json:"images,omitempty"`
}

type projectReportSourceImage struct {
	FileName  string                        `json:"file_name,omitempty"`
	Detection *projectReportSourceDetection `json:"detection,omitempty"`
	Review    *projectReportSourceReview    `json:"review,omitempty"`
}

type projectReportSourceDetection struct {
	Corrosion *projectReportSourceCorrosion `json:"corrosion,omitempty"`
	Crack     *projectReportSourceCrack     `json:"crack,omitempty"`
	Stain     *projectReportSourceStain     `json:"stain,omitempty"`
	Flatness  *projectReportSourceFlatness  `json:"flatness,omitempty"`
	Spalling  *projectReportSourceSpalling  `json:"spalling,omitempty"`
	Summary   *projectReportSourceSummary   `json:"summary,omitempty"`
}

type projectReportSourceReview struct {
	Verdict string `json:"verdict,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type projectReportSourceCorrosion struct {
	HasCorrosion      bool                                  `json:"has_corrosion,omitempty"`
	CorrosionCount    uint32                                `json:"corrosion_count,omitempty"`
	MaxConfidence     float64                               `json:"max_confidence,omitempty"`
	AverageConfidence float64                               `json:"average_confidence,omitempty"`
	CorrosionPixels   uint64                                `json:"corrosion_pixels,omitempty"`
	CorrosionRatio    float64                               `json:"corrosion_ratio,omitempty"`
	Regions           []*projectReportSourceCorrosionRegion `json:"regions,omitempty"`
}

type projectReportSourceCorrosionRegion struct {
	Confidence float64 `json:"confidence,omitempty"`
	MaskPixels uint64  `json:"mask_pixels,omitempty"`
	MaskRatio  float64 `json:"mask_ratio,omitempty"`
}

type projectReportSourceCrack struct {
	HasCrack    bool                              `json:"has_crack,omitempty"`
	CrackCount  uint32                            `json:"crack_count,omitempty"`
	CrackPixels uint64                            `json:"crack_pixels,omitempty"`
	CrackRatio  float64                           `json:"crack_ratio,omitempty"`
	Regions     []*projectReportSourceCrackRegion `json:"regions,omitempty"`
}

type projectReportSourceCrackRegion struct {
	MaskPixels uint64  `json:"mask_pixels,omitempty"`
	MaskRatio  float64 `json:"mask_ratio,omitempty"`
}

type projectReportSourceStain struct {
	HasStain          bool                              `json:"has_stain,omitempty"`
	StainCount        uint32                            `json:"stain_count,omitempty"`
	AverageStainRatio float64                           `json:"average_stain_ratio,omitempty"`
	MaxStainRatio     float64                           `json:"max_stain_ratio,omitempty"`
	Regions           []*projectReportSourceStainRegion `json:"regions,omitempty"`
}

type projectReportSourceStainRegion struct {
	Confidence   float64 `json:"confidence,omitempty"`
	RegionWidth  uint32  `json:"region_width,omitempty"`
	RegionHeight uint32  `json:"region_height,omitempty"`
	StainPixels  uint64  `json:"stain_pixels,omitempty"`
	StainRatio   float64 `json:"stain_ratio,omitempty"`
}

type projectReportSourceFlatness struct {
	Result      string                               `json:"result,omitempty"`
	UnevenCount uint32                               `json:"uneven_count,omitempty"`
	Regions     []*projectReportSourceFlatnessRegion `json:"regions,omitempty"`
}

type projectReportSourceFlatnessRegion struct {
	EdgeUnevenDetected      bool    `json:"edge_uneven_detected,omitempty"`
	LineUnevenDetected      bool    `json:"line_uneven_detected,omitempty"`
	GradientUnevenDetected  bool    `json:"gradient_uneven_detected,omitempty"`
	FrequencyUnevenDetected bool    `json:"frequency_uneven_detected,omitempty"`
	EdgeCount               uint64  `json:"edge_count,omitempty"`
	LaplacianVariance       float64 `json:"laplacian_variance,omitempty"`
	LineCount               uint64  `json:"line_count,omitempty"`
	AngleStd                float64 `json:"angle_std,omitempty"`
	GradientMean            float64 `json:"gradient_mean,omitempty"`
	GradientStd             float64 `json:"gradient_std,omitempty"`
	FrequencyMin            float64 `json:"frequency_min,omitempty"`
	FrequencyMax            float64 `json:"frequency_max,omitempty"`
}

type projectReportSourceSpalling struct {
	HasSpalling bool    `json:"has_spalling,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
}

type projectReportSourceSummary struct {
	Result string `json:"result,omitempty"`
}
