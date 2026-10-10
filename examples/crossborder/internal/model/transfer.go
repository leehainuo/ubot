package model

import "time"

type TransferStatus string

const (
	TransferValidated TransferStatus = "validated"
	TransferCreated   TransferStatus = "created"
)

type Transfer struct {
	ID            string
	SKU           string
	FromWarehouse string
	ToWarehouse   string
	Quantity      int
	Status        TransferStatus
	CreatedAt     time.Time
}

type IdempotencyRecord struct {
	Fingerprint [32]byte
	Transfer    Transfer
}
