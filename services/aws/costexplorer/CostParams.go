package costexplorer

import (
	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type CostParams struct {
	StartDate   string
	EndDate     string
	Granularity costexplorertypes.Granularity
}
