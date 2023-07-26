package jobs

import (
	"cspr-fetcher/data"
	"cspr-fetcher/models"
	"errors"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
)

type BlockInfoFetcher struct {
	db        data.DBConnectorI
	rpcClient data.RpcClientI
}

func NewBlockInfoFetcher(db data.DBConnectorI, rpcClient data.RpcClientI) *BlockInfoFetcher {
	return &BlockInfoFetcher{db: db, rpcClient: rpcClient}
}

func (b BlockInfoFetcher) FetchBlockInfo() error {
	lastFetchedBlock, err := b.db.GetLatestFetchedBlock()
	if err != nil {
		return err
	}

	if lastFetchedBlock == nil {
		return errors.New("no data fetched. Backfill job should be triggered")
	}

	latestBlock, err := b.rpcClient.GetLatestBlock()
	if err != nil {
		return err
	}

	if uint64(latestBlock.Header.Height)-lastFetchedBlock.Height == 0 {
		// no data to fetch
		return nil
	}

	height := lastFetchedBlock.Height + 1
	for height <= uint64(latestBlock.Header.Height) {
		block, err := b.rpcClient.GetBlockByHeight(height)
		if err != nil {
			return err
		}

		err = b.db.InsertNewBlock(convertBlock(block))
		if err != nil {
			return err
		}

		transfers, err := b.rpcClient.GetBlockTransfersByHeight(height)
		if err != nil {
			return err
		}

		transferList := convertTransfersLocalObject(block, transfers)
		err = b.db.InsertTransfers(transferList)
		if err != nil {
			return err
		}

		height++
	}

	return nil
}

func convertTransfersLocalObject(block sdk.BlockResponse, transfers []sdk.TransferResponse) []models.Transfer {
	var list []models.Transfer
	for _, tr := range transfers {
		list = append(list, models.Transfer{
			BlockHash:   block.Hash,
			BlockHeight: uint64(block.Header.Height),
			FromAccount: tr.From,
			ToAccount:   tr.To,
			Amount:      tr.Amount,
			Gas:         tr.Gas,
		})
	}
	return list
}

func convertBlock(block sdk.BlockResponse) models.Block {
	return models.Block{
		Hash:      block.Hash,
		Height:    uint64(block.Header.Height),
		EraId:     uint64(block.Header.EraID),
		Timestamp: block.Header.Timestamp,
	}
}
