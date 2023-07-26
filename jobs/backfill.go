package jobs

import (
	"cspr-fetcher/data"
)

type Backfill struct {
	db        data.DBConnectorI
	rpcClient data.RpcClientI
}

func NewBackfill(db data.DBConnectorI, rpcClient data.RpcClientI) *Backfill {
	return &Backfill{db: db, rpcClient: rpcClient}
}

func (b Backfill) FetchDataWorker(startLength, endLength uint64) error {
	ind := startLength
	for ind <= endLength {
		block, err := b.rpcClient.GetBlockByHeight(ind)
		if err != nil {
			return err
		}

		err = b.db.InsertNewBlock(convertBlock(block))
		if err != nil {
			return err
		}

		transfers, err := b.rpcClient.GetBlockTransfersByHeight(ind)
		if err != nil {
			return err
		}

		transferList := convertTransfersLocalObject(block, transfers)
		err = b.db.InsertTransfers(transferList)
		if err != nil {
			return err
		}

		ind++
	}

	return nil
}
