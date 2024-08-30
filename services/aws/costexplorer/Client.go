package costexplorer

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/kimulaco/tool-cost-summarizer/services/aws/config"
)

type Client struct {
	Client ClientType
}

type ClientType interface {
	GetCostAndUsage(
		ctx context.Context,
		params *costexplorer.GetCostAndUsageInput,
		optFns ...func(*costexplorer.Options),
	) (*costexplorer.GetCostAndUsageOutput, error)
}

func NewClient(config config.Config) *Client {
	client := costexplorer.NewFromConfig(config.Config)

	return &Client{Client: client}
}

func (c *Client) GetCostSummary(params CostParams) ([]Result, error) {
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorertypes.DateInterval{
			Start: aws.String(params.StartDate),
			End:   aws.String(params.EndDate),
		},
		Metrics:     []string{"UnblendedCost"},
		Granularity: params.Granularity,
		GroupBy: []costexplorertypes.GroupDefinition{
			{
				Type: costexplorertypes.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	}

	cost, err := c.Client.GetCostAndUsage(context.TODO(), input)
	if err != nil {
		return []Result{}, err
	}

	results := []Result{}
	for _, resultByTime := range cost.ResultsByTime {
		result := Result{
			StartDate:   *resultByTime.TimePeriod.Start,
			EndDate:     *resultByTime.TimePeriod.End,
			TotalAmount: 0.0,
			Breakdown:   []BreakdownResult{},
		}

		for _, group := range resultByTime.Groups {
			amount, err := strconv.ParseFloat(*group.Metrics["UnblendedCost"].Amount, 64)
			if err != nil {
				log.Printf("failed to parse amount: %v", err)
				continue
			}
			if amount <= 0.0 {
				continue
			}

			breakdown := BreakdownResult{
				Name:   group.Keys[0],
				Amount: amount,
				Unit:   *group.Metrics["UnblendedCost"].Unit,
			}
			result.Breakdown = append(result.Breakdown, breakdown)
			result.TotalAmount += amount
		}

		results = append(results, result)
	}

	return results, nil
}
