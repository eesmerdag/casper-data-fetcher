package data

import (
	"cspr-fetcher/models"
	"database/sql"
	"fmt"
	"strings"
)

type DBConnectorI interface {
	GetAllBlocks(limit, offset uint64) ([]*models.Block, error)
	GetLatestFetchedBlock() (*models.Block, error)
	GetBlockByHeight(height uint64) (*models.Block, error)
	InsertNewBlock(bl models.Block) error
	InsertTransfers(transfers []models.Transfer) error
}

type DBConnector struct {
	db *sql.DB
}

func NewDBConnector(db *sql.DB) *DBConnector {
	return &DBConnector{db: db}
}

func (c DBConnector) InsertNewBlock(block models.Block) error {
	_, e := c.db.Exec(`insert into "blocks" ("hash",	"height", "era_id", "timestamp") values ($1, $2, $3, $4)`, block.Hash, block.Height, block.EraId, block.Timestamp)
	return e
}

func (c DBConnector) GetAllBlocks(limit, offset uint64) ([]*models.Block, error) {
	rows, err := c.db.Query(`SELECT "hash", "height", "era_id", "timestamp" FROM "blocks" order by height limit $1 offset $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	blocks := make([]*models.Block, 0)

	for rows.Next() {

		block := new(models.Block)
		if err := rows.Scan(&block.Height, &block.EraId, &block.Timestamp); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (c DBConnector) GetLatestFetchedBlock() (*models.Block, error) {
	block := new(models.Block)
	row := c.db.QueryRow(`SELECT "hash", "height", "era_id", "timestamp" FROM "blocks" order by height desc limit 1`)

	err := row.Scan(&block.Height, &block.EraId, &block.Timestamp)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c DBConnector) GetBlockByHeight(height uint64) (*models.Block, error) {
	block := new(models.Block)
	row := c.db.QueryRow(`SELECT "hash", "height", "era_id", "timestamp" FROM "blocks" where height = $1 order by height desc limit 1`, height)

	err := row.Scan(&block.Height, &block.EraId, &block.Timestamp)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c DBConnector) InsertTransfers(transfers []models.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}

	values := make([]string, 0, len(transfers))
	args := make([]interface{}, 0, len(transfers)*6)
	for i, tr := range transfers {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		args = append(args, tr.BlockHash)
		args = append(args, tr.BlockHeight)
		args = append(args, tr.FromAccount)
		args = append(args, tr.ToAccount)
		args = append(args, tr.Amount)
		args = append(args, tr.Gas)
	}

	stmt := fmt.Sprintf("insert into transfers (block_hash, block_height, from_account, to_account, amount, gas) values %s",
		strings.Join(values, ","))
	_, err := c.db.Exec(stmt, args...)
	return err
}
