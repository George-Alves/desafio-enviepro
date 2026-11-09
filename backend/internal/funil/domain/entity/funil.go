package entity

import "time"

// Etapas do funil de vendas, na ordem em que o negócio avança.
var Etapas = []string{"Lead", "Contato", "Proposta", "Negociação", "Fechado"}

// Oportunidade é uma negociação em andamento com um cliente.
type Oportunidade struct {
	ID          int64
	ClienteID   int64
	Titulo      string
	Etapa       string
	Valor       float64
	DataCriacao time.Time
}

// TotalEtapa é o agregado de uma etapa do funil.
type TotalEtapa struct {
	Etapa      string
	Quantidade int
	Valor      float64
}
