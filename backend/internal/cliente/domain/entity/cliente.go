package entity

import "time"

// Cliente é uma pessoa ou empresa acompanhada pelo CRM.
type Cliente struct {
	ID            int64
	Nome          string
	Empresa       string
	Email         string
	Telefone      string
	Status        string
	Responsavel   string
	UltimoContato time.Time
}
