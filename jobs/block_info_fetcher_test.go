package jobs

import (
	"cspr-fetcher/data"
	"cspr-fetcher/models"
	"errors"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestFetchBlockInfo(t *testing.T) {
	tests := []struct {
		name               string
		backfillRequired   bool
		errorFromRpcClient bool
		errorExpected      bool
	}{
		{
			name:          "successful",
			errorExpected: false,
		},
		{
			name:             "backfill is required",
			backfillRequired: true,
			errorExpected:    true,
		},
		{
			name:               "error from rpc client",
			errorFromRpcClient: true,
			errorExpected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRpcClient := data.NewMockRpcClientI(mockCtrl)
			mockDbConnector := data.NewMockDBConnectorI(mockCtrl)

			if !tt.errorExpected {
				latestFetchedBlock := models.Block{
					Height:    uint64(100),
					Hash:      "hash-test-100",
					EraId:     uint64(1),
					Timestamp: time.Now().AddDate(0, 0, -1),
				}
				mockDbConnector.EXPECT().GetLatestFetchedBlock().Return(&latestFetchedBlock, nil)

				latestBlock := sdk.BlockResponse{
					Hash: "hash-test-105",
					Header: sdk.BlockHeader{
						Height:    105,
						EraID:     1,
						Timestamp: time.Now(),
					},
				}
				mockRpcClient.EXPECT().GetLatestBlock().Return(latestBlock, nil)

				for i := 101; i <= 105; i++ {
					blck := sdk.BlockResponse{
						Hash: fmt.Sprintf("hash-test-%s", strconv.Itoa(i)),
						Header: sdk.BlockHeader{
							Height:    105,
							EraID:     1,
							Timestamp: time.Now(),
						},
					}
					mockRpcClient.EXPECT().GetBlockByHeight(uint64(i)).Return(blck, nil)
					mockDbConnector.EXPECT().InsertNewBlock(convertBlock(blck)).Return(nil)

					transfers := []sdk.TransferResponse{
						{
							From:   "from-test",
							To:     "to-test",
							Amount: "amount-test",
							Gas:    "gas-test",
						},
					}
					mockRpcClient.EXPECT().GetBlockTransfersByHeight(uint64(i)).Return(transfers, nil)
					mockDbConnector.EXPECT().InsertTransfers(convertTransfersLocalObject(blck, transfers)).Return(nil)
				}
			} else {
				if tt.backfillRequired {
					mockDbConnector.EXPECT().GetLatestFetchedBlock().Return(nil, nil)
				}
				if tt.errorFromRpcClient {
					latestFetchedBlock := models.Block{
						Height:    uint64(100),
						Hash:      "hash-test-100",
						EraId:     uint64(1),
						Timestamp: time.Now().AddDate(0, 0, -1),
					}
					mockDbConnector.EXPECT().GetLatestFetchedBlock().Return(&latestFetchedBlock, nil)
					mockRpcClient.EXPECT().GetLatestBlock().Return(sdk.BlockResponse{}, errors.New("some problem at rpc client"))
				}
			}

			fetcher := NewBlockInfoFetcher(mockDbConnector, mockRpcClient)
			err := fetcher.FetchBlockInfo()

			if tt.errorExpected && tt.backfillRequired {
				assert.Error(t, err)
				assert.Equal(t, "no data fetched. Backfill job should be triggered", err.Error())
			} else if tt.errorExpected && tt.errorFromRpcClient {
				assert.Error(t, err)
				assert.Equal(t, "some problem at rpc client", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
