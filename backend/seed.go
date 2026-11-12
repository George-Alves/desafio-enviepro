package main

import (
	"context"
	cripto "crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	funil "github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
	usuarioRepo "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/repository"
	usuarioUC "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/app/usecase"
	usuario "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	usuarioPorts "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
)

const formatoData = "2006-01-02"

const (
	adminNome  = "Administrador"
	adminEmail = "admin@enviepro.com"
)

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
	// Só zera a numeração das tabelas fictícias. A de usuarios fica: se o último usuário foi
	// apagado, o id dele não pode ser reaproveitado por outra pessoa.
	if _, err := tx.Exec("DELETE FROM sqlite_sequence WHERE name IN ('vendas', 'oportunidades', 'clientes')"); err != nil {
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

// gerarSenhaAleatoria sorteia uma senha forte de 16 caracteres com o gerador seguro do
// sistema (crypto/rand), e não com o math/rand, que é previsível.
func gerarSenhaAleatoria() (string, error) {
	b := make([]byte, 12)
	if _, err := cripto.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// garantirAdmin cria o administrador inicial, mas só se ainda não existir nenhum.
// Usa o mesmo caso de uso da API, então a validação e o bcrypt valem aqui também.
// Se senhaConfigurada estiver vazia, sorteia uma senha e a devolve em senhaGerada.
func garantirAdmin(ctx context.Context, db *sql.DB, senhaConfigurada string) (criado bool, senhaGerada string, err error) {
	repo := usuarioRepo.NewSQLiteUsuarioRepository(db)

	existentes, err := usuarioUC.NewListarUsuarios(repo).Executar(ctx)
	if err != nil {
		return false, "", err
	}
	for _, u := range existentes {
		if u.Papel == usuario.PapelAdmin {
			return false, "", nil
		}
	}

	senha := senhaConfigurada
	if senha == "" {
		if senha, err = gerarSenhaAleatoria(); err != nil {
			return false, "", err
		}
		senhaGerada = senha
	}

	_, err = usuarioUC.NewCriarUsuario(repo).Executar(ctx, usuarioUC.DadosNovoUsuario{
		Nome: adminNome, Email: adminEmail, Senha: senha, Papel: usuario.PapelAdmin,
	})
	var invalido usuario.ErrosValidacao
	switch {
	case err == nil:
		return true, senhaGerada, nil
	case senhaConfigurada != "" && errors.As(err, &invalido):
		return false, "", fmt.Errorf("ADMIN_SENHA inválida: %w", err)
	case errors.Is(err, usuarioPorts.ErrEmailJaCadastrado):
		return false, "", fmt.Errorf("já existe um usuário com o e-mail %s que não é administrador", adminEmail)
	}
	return false, "", err
}

// prepararAdmin lê ADMIN_SENHA, garante o administrador inicial e conta no log o que aconteceu.
func prepararAdmin(db *sql.DB) error {
	configurada := os.Getenv("ADMIN_SENHA")
	criado, gerada, err := garantirAdmin(context.Background(), db, configurada)
	if err != nil {
		return err
	}
	switch {
	case criado && gerada != "":
		log.Printf("Administrador criado: %s | senha sorteada: %s", adminEmail, gerada)
		log.Println("ANOTE A SENHA AGORA: ela não será mostrada de novo. Para escolher a sua, defina ADMIN_SENHA antes de iniciar.")
	case criado:
		log.Printf("Administrador criado: %s (senha definida por ADMIN_SENHA)", adminEmail)
	case configurada != "":
		log.Println("ADMIN_SENHA ignorada: já existe um administrador. A senha dele não é alterada.")
	}
	return nil
}
