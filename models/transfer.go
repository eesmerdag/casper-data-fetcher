package models

type Transfer struct {
	BlockHash   string `json:"block_hash"`
	BlockHeight uint64 `json:"block_height"`
	FromAccount string `json:"from"`
	ToAccount   string `json:"to"`
	Amount      string `json:"amount"`
	Gas         string `json:"gas"`
}
