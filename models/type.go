package models

type Gate struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Duration     int      `json:"duration"`
	Dependencies []string `json:"dependencies"`
}

type SchedulingRequest struct {
	Graph       []Gate   `json:"graph"`
	Resources   []string `json:"resources"`
	MaxLatency  int      `json:"maxLatency"`
	MaxResource int      `json:"maxResource"`
}

type SchedulingResult struct {
	Schedule      map[string]int         `json:"schedule"`
	TotalLatency  int                    `json:"totalLatency"`
	ResourceUsage map[int]map[string]int `json:"resourceUsage"`
}
