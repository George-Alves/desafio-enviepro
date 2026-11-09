package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	funil "github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
)

const formatoData = "2006-01-02"

var (
	nomesClientes = []string{
		"Ana Souza", "Bruno Lima", "Carla Mendes", "Diego Santos", "Eduarda Rocha",
		"Fábio Alves", "Gabriela Costa", "Henrique Dias", "Isabela Martins", "João Pereira",
		"Karina Ribeiro", "Lucas Barbosa", "Marina Carvalho", "Nelson Teixeira", "Olívia Nunes",
		"Paulo Araújo", "Renata Gomes", "Sérgio Freitas", "Tatiana Moreira", "Vinícius Pinto",
	}
	empresas = []string{
		"Alfa Logística", "Beta Têxtil", "Gama Alimentos", "Delta Móveis", "Ômega Tech",
		"Sigma Farma", "Lambda Pet", "Zeta Cosméticos", "Kappa Esportes", "Épsilon Casa",
	}
	responsaveis    = []string{"Marcos", "Júlia", "Rafael", "Camila"}
	statusPossiveis = []string{"ativo", "ativo", "em_acompanhamento", "inativo"}
	produtos        = []string{"Plano Básico", "Plano Pro", "Plano Premium", "Consultoria", "Integração"}
)

// semear apaga os dados atuais e grava dados fictícios: 20 clientes, 30 oportunidades e 25 vendas.
// As datas são relativas a "agora", então sempre há vendas no mês e clientes sem contato recente.
// A mesma data de entrada gera sempre os mesmos dados.
func semear(db *sql.DB, agora time.Time) error {
	rng := rand.New(rand.NewSource(42))
	hoje := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, time.UTC)
	diasAtras := func(n int) string { return hoje.AddDate(0, 0, -n).Format(formatoData) }

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tabela := range []string{"vendas", "oportunidades", "clientes"} {
		if _, err := tx.Exec("DELETE FROM " + tabela); err != nil {
			return err
		}
	}
	if _, err := tx.Exec("DELETE FROM sqlite_sequence"); err != nil {
		return err
	}

	ids := make([]int64, 0, len(nomesClientes))
	for i, nome := range nomesClientes {
		// A cada 3 clientes, um está há mais de 15 dias sem contato.
		dias := rng.Intn(15)
		if i%3 == 0 {
			dias = 16 + rng.Intn(30)
		}
		res, err := tx.Exec(`
			INSERT INTO clientes (nome, empresa, email, telefone, status, responsavel, ultimo_contato)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			nome, empresas[i%len(empresas)],
			fmt.Sprintf("cliente%02d@exemplo.com.br", i+1),
			fmt.Sprintf("(11) 9%04d-%04d", 1000+i, 2000+rng.Intn(7000)),
			statusPossiveis[i%len(statusPossiveis)],
			responsaveis[rng.Intn(len(responsaveis))],
			diasAtras(dias))
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for i := 0; i < 30; i++ {
		_, err := tx.Exec(`
			INSERT INTO oportunidades (cliente_id, titulo, etapa_funil, valor, data_criacao)
			VALUES (?, ?, ?, ?, ?)`,
			ids[rng.Intn(len(ids))],
			fmt.Sprintf("%s — proposta %d", produtos[rng.Intn(len(produtos))], i+1),
			funil.Etapas[i%len(funil.Etapas)],
			float64(30+rng.Intn(570))*100,
			diasAtras(rng.Intn(90)))
		if err != nil {
			return err
		}
	}

	inicioDoMes := time.Date(hoje.Year(), hoje.Month(), 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 25; i++ {
		var data time.Time
		if i < 12 {
			// 12 vendas no mês atual, do dia 1 até hoje.
			data = inicioDoMes.AddDate(0, 0, rng.Intn(hoje.Day()))
		} else {
			// As outras, nos meses anteriores.
			data = inicioDoMes.AddDate(0, 0, -(1 + rng.Intn(90)))
		}
		_, err := tx.Exec(`
			INSERT INTO vendas (cliente_id, valor, data, produto) VALUES (?, ?, ?, ?)`,
			ids[rng.Intn(len(ids))],
			float64(16+rng.Intn(280))*50,
			data.Format(formatoData),
			produtos[rng.Intn(len(produtos))])
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
