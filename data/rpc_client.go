package data

import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

type RpcClientI interface {
	GetLatestBlock() (sdk.BlockResponse, error)
	GetBlockByHeight(height uint64) (sdk.BlockResponse, error)
	GetBlockTransfersByHeight(height uint64) ([]sdk.TransferResponse, error)
}
