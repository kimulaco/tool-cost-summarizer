package costexplorer

type Result struct {
	StartDate   string            `json:"startDate"`
	EndDate     string            `json:"endDate"`
	TotalAmount float64           `json:"totalAmount"`
	Breakdown   []BreakdownResult `json:"breakdown"`
}

type BreakdownResult struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}
