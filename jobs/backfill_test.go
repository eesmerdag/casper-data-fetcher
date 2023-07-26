package jobs

import (
	"cspr-fetcher/data"
	"errors"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestFetchDataWorker(t *testing.T) {
	tests := []struct {
		name          string
		start         uint64
		end           uint64
		errorExpected bool
	}{
		{
			name:          "successful",
			start:         4,
			end:           6,
			errorExpected: false,
		},
		{
			name:          "fails",
			start:         4,
			end:           6,
			errorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRpcClient := data.NewMockRpcClientI(mockCtrl)
			mockDbConnector := data.NewMockDBConnectorI(mockCtrl)

			if !tt.errorExpected {
				for i := tt.start; i <= tt.end; i++ {
					blck := sdk.BlockResponse{
						Hash: fmt.Sprintf("hash-test-%s", strconv.FormatUint(i, 10)),
						Header: sdk.BlockHeader{
							Height:    int(i),
							EraID:     0,
							Timestamp: time.Now().AddDate(0, 0, -1),
						},
					}
					mockRpcClient.EXPECT().GetBlockByHeight(i).Return(blck, nil)
					mockDbConnector.EXPECT().InsertNewBlock(convertBlock(blck)).Return(nil)

					transfers := []sdk.TransferResponse{
						{
							From:   "from-test",
							To:     "to-test",
							Amount: "amount-test",
							Gas:    "gas-test",
						},
					}
					mockRpcClient.EXPECT().GetBlockTransfersByHeight(i).Return(transfers, nil)
					mockDbConnector.EXPECT().InsertTransfers(convertTransfersLocalObject(blck, transfers)).Return(nil)
				}
			} else {
				mockRpcClient.EXPECT().GetBlockByHeight(tt.start).Return(sdk.BlockResponse{}, errors.New("something is wrong"))
			}

			fetcher := NewBackfill(mockDbConnector, mockRpcClient)
			err := fetcher.FetchDataWorker(tt.start, tt.end)

			if tt.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
