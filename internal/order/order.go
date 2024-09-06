package order

type Metadata struct {
	Constellation string `json:"constellation"`
}

type Rules struct {
	CloudCoveragePct   int `json:"cloud_coverage_pct"`
	AISResolutionMaxCm int `json:"is_resolution_max_cm"`
	AISResolutionMinCm int `json:"ais_resolution_min_cm`
	EOResolutionMaxCm  int `json:"eo_resolution_max_cm"`
	EOResolutionMinCm  int `json:"eo_resolution_min_cm`
	HSIResolutionMaxCm int `json:"hsi_resolution_max_cm"`
	HSIResolutionMinCm int `json:"hsi_resolution_min_cm`
	RFResolutionMaxCm  int `json:"rf_resolution_max_cm"`
	RFResolutionMinCm  int `json:"rf_resolution_min_cm`
	SARResolutionMaxCm int `json:"sar_resolution_max_cm"`
	SARResolutionMixCm int `json:"sar_resolution_min_cm"`
}
