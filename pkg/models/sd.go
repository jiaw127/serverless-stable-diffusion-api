package models

type Txt2ImgResult struct {
	Images     []string               `json:"images"`
	Parameters map[string]interface{} `json:"parameters"`
	Info       string                 `json:"info"`
}

type ProgressResult struct {
	CurrentImage string  `json:"current_image"`
	EtaRelative  float64 `json:"eta_relative"`
	Progress     float64 `json:"progress"`
	State        State   `json:"state"`
}

type State struct {
	Interrupted   bool   `json:"interrupted"`
	Job           string `json:"job"`
	JobCount      int    `json:"job_count"`
	JobNo         int    `json:"job_no"`
	JobTimestamp  string `json:"job_timestamp"`
	SamplingStep  int    `json:"sampling_step"`
	SamplingSteps int    `json:"sampling_steps"`
	Skipped       bool   `json:"skipped"`
}
