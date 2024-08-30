package costexplorer

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/kimulaco/tool-cost-summarizer/services/aws/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCostExplorerClient struct {
	mock.Mock
}

func (m *mockCostExplorerClient) GetCostAndUsage(
	ctx context.Context,
	params *costexplorer.GetCostAndUsageInput,
	optFns ...func(*costexplorer.Options),
) (*costexplorer.GetCostAndUsageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*costexplorer.GetCostAndUsageOutput), args.Error(1)
}

func TestNewClient(t *testing.T) {
	cfg, _ := config.NewConfig("testAccessKey", "testSecretKey")
	client := NewClient(cfg)
	assert.NotNil(t, client)
}

func TestGetCostSummary_Success(t *testing.T) {
	mockClient := new(mockCostExplorerClient)
	client := &Client{Client: mockClient}

	startDate := "2023-01-01"
	endDate := "2023-01-31"
	params := CostParams{
		StartDate:   startDate,
		EndDate:     endDate,
		Granularity: costexplorertypes.GranularityMonthly,
	}

	mockOutput := &costexplorer.GetCostAndUsageOutput{
		ResultsByTime: []costexplorertypes.ResultByTime{
			{
				TimePeriod: &costexplorertypes.DateInterval{
					Start: &startDate,
					End:   &endDate,
				},
				Groups: []costexplorertypes.Group{
					{
						Keys: []string{"Lambda"},
						Metrics: map[string]costexplorertypes.MetricValue{
							"UnblendedCost": {
								Amount: stringPtr("10.1"),
								Unit:   stringPtr("USD"),
							},
						},
					},
					{
						Keys: []string{"CloudFront"},
						Metrics: map[string]costexplorertypes.MetricValue{
							"UnblendedCost": {
								Amount: stringPtr("0"),
								Unit:   stringPtr("USD"),
							},
						},
					},
					{
						Keys: []string{"S3"},
						Metrics: map[string]costexplorertypes.MetricValue{
							"UnblendedCost": {
								Amount: stringPtr("0.5"),
								Unit:   stringPtr("USD"),
							},
						},
					},
				},
			},
		},
	}

	mockClient.On("GetCostAndUsage", mock.Anything, mock.Anything, mock.Anything).Return(mockOutput, nil)

	results, err := client.GetCostSummary(params)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, startDate, results[0].StartDate)
	assert.Equal(t, endDate, results[0].EndDate)
	assert.InDelta(t, 10.6, results[0].TotalAmount, 0.001)
	assert.Len(t, results[0].Breakdown, 2)

	assert.Equal(t, "Lambda", results[0].Breakdown[0].Name)
	assert.InDelta(t, 10.1, results[0].Breakdown[0].Amount, 0.001)
	assert.Equal(t, "USD", results[0].Breakdown[0].Unit)

	assert.Equal(t, "S3", results[0].Breakdown[1].Name)
	assert.InDelta(t, 0.5, results[0].Breakdown[1].Amount, 0.001)
	assert.Equal(t, "USD", results[0].Breakdown[1].Unit)

	mockClient.AssertExpectations(t)
}

func TestGetCostSummary_GetCostAndUsageError(t *testing.T) {
	mockClient := new(mockCostExplorerClient)
	client := &Client{Client: mockClient}

	startDate := "2023-01-01"
	endDate := "2023-01-31"
	params := CostParams{
		StartDate:   startDate,
		EndDate:     endDate,
		Granularity: costexplorertypes.GranularityMonthly,
	}

	mockClient.On("GetCostAndUsage", mock.Anything, mock.Anything, mock.Anything).Return(&costexplorer.GetCostAndUsageOutput{}, errors.New("mock error"))

	results, err := client.GetCostSummary(params)

	assert.Equal(t, "mock error", err.Error())
	assert.Len(t, results, 0)

	mockClient.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
