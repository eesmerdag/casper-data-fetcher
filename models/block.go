package models

import (
	"time"
)

type Block struct {
	Hash      string    `json:"hash"`
	Height    uint64    `json:"height"`
	EraId     uint64    `json:"era_id"`
	Timestamp time.Time `json:"timestamp"`
}
