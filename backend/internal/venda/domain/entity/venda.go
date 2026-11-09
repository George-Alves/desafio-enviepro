package entity

import "time"

// Venda é uma venda fechada com um cliente.
type Venda struct {
	ID        int64
	ClienteID int64
	Valor     float64
	Data      time.Time
	Produto   string
}
