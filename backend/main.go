package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	clienteHandler "github.com/George-Alves/desafio-enviepro/backend/internal/cliente/adapter/handler"
	clienteRepo "github.com/George-Alves/desafio-enviepro/backend/internal/cliente/adapter/repository"
	clienteUC "github.com/George-Alves/desafio-enviepro/backend/internal/cliente/app/usecase"
	funilHandler "github.com/George-Alves/desafio-enviepro/backend/internal/funil/adapter/handler"
	funilRepo "github.com/George-Alves/desafio-enviepro/backend/internal/funil/adapter/repository"
	funilUC "github.com/George-Alves/desafio-enviepro/backend/internal/funil/app/usecase"
	"github.com/George-Alves/desafio-enviepro/backend/internal/shared/middleware"
	usuarioHandler "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/handler"
	usuarioRepo "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/repository"
	usuarioUC "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/app/usecase"
	vendaHandler "github.com/George-Alves/desafio-enviepro/backend/internal/venda/adapter/handler"
	vendaRepo "github.com/George-Alves/desafio-enviepro/backend/internal/venda/adapter/repository"
	vendaUC "github.com/George-Alves/desafio-enviepro/backend/internal/venda/app/usecase"
	"github.com/gin-gonic/gin"
)

// Uso:
//
//	go run .        sobe a API (cria e popula o banco se ele ainda estiver vazio)
//	go run . seed   apaga os dados fictícios e grava novos (os usuários não são apagados)
//
// Nas duas formas, se ainda não existir um administrador, ele é criado com a senha de
// ADMIN_SENHA; sem essa variável, uma senha é sorteada e mostrada uma única vez no log.
func main() {
	db, err := abrirBanco(valorOuPadrao("DB_PATH", "crm.db"))
	if err != nil {
		log.Fatalf("erro ao abrir o banco: %v", err)
	}
	defer db.Close()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := semear(db, time.Now()); err != nil {
			log.Fatalf("erro ao popular o banco: %v", err)
		}
		log.Println("banco populado com dados fictícios")
		if err := prepararAdmin(db); err != nil {
			log.Fatalf("erro ao criar o administrador: %v", err)
		}
		return
	}

	if vazio, err := bancoVazio(db); err != nil {
		log.Fatalf("erro ao consultar o banco: %v", err)
	} else if vazio {
		if err := semear(db, time.Now()); err != nil {
			log.Fatalf("erro ao popular o banco: %v", err)
		}
		log.Println("banco estava vazio: populado com dados fictícios")
	}
	if err := prepararAdmin(db); err != nil {
		log.Fatalf("erro ao criar o administrador: %v", err)
	}

	porta := valorOuPadrao("PORT", "8080")
	log.Printf("API rodando em http://localhost:%s", porta)
	if err := montarRoteador(db).Run(":" + porta); err != nil {
		log.Fatal(err)
	}
}

// montarRoteador liga as peças: repositório -> caso de uso -> handler.
// É o único lugar que conhece todas as implementações.
func montarRoteador(db *sql.DB) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.Use(middleware.CORS(middleware.OrigemLocal))

	clientes := clienteRepo.NewSQLiteClienteRepository(db)
	funil := funilRepo.NewSQLiteFunilRepository(db)
	vendas := vendaRepo.NewSQLiteVendaRepository(db)
	usuarios := usuarioRepo.NewSQLiteUsuarioRepository(db)

	usuarioHandler.NewUsuarioHandler(
		usuarioUC.NewCriarUsuario(usuarios),
		usuarioUC.NewListarUsuarios(usuarios),
		usuarioUC.NewBuscarUsuario(usuarios),
		usuarioUC.NewAtualizarUsuario(usuarios),
		usuarioUC.NewApagarUsuario(usuarios),
	).Registrar(r)
	clienteHandler.NewClienteHandler(clienteUC.NewListarClientes(clientes)).Registrar(r)
	funilHandler.NewFunilHandler(funilUC.NewObterFunil(funil)).Registrar(r)
	vendaHandler.NewVendaHandler(
		vendaUC.NewResumoVendasDoMes(vendas),
		vendaUC.NewListarVendas(vendas),
	).Registrar(r)

	return r
}

func valorOuPadrao(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}
