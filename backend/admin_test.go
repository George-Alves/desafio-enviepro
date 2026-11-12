package main

import (
	"context"
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func novoBancoDeTeste(t *testing.T) *sql.DB {
	t.Helper()
	db, err := abrirBanco(filepath.Join(t.TempDir(), "admin.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func hashDoAdmin(t *testing.T, db *sql.DB) (email, hash, papel string) {
	t.Helper()
	if err := db.QueryRow("SELECT email, senha_hash, papel FROM usuarios WHERE papel = 'admin'").Scan(&email, &hash, &papel); err != nil {
		t.Fatalf("não achei o admin no banco: %v", err)
	}
	return
}

func TestGarantirAdmin_SemSenhaConfigurada_SorteiaSenhaEGravaSoOHash(t *testing.T) {
	db := novoBancoDeTeste(t)

	criado, gerada, err := garantirAdmin(context.Background(), db, "")
	if err != nil || !criado {
		t.Fatalf("esperava criar o admin, veio criado=%v erro=%v", criado, err)
	}
	if len(gerada) < 12 {
		t.Errorf("a senha sorteada deveria ser longa, veio %q", gerada)
	}

	email, hash, papel := hashDoAdmin(t, db)
	if email != "admin@enviepro.com" || papel != "admin" {
		t.Errorf("admin inesperado: %s / %s", email, papel)
	}
	if !strings.HasPrefix(hash, "$2a$") || strings.Contains(hash, gerada) {
		t.Errorf("o banco deveria ter só o hash $2a$, veio %q", hash)
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(gerada)) != nil {
		t.Error("o hash deveria conferir com a senha sorteada")
	}
}

func TestGarantirAdmin_ComSenhaConfigurada_NaoSorteiaNada(t *testing.T) {
	db := novoBancoDeTeste(t)

	criado, gerada, err := garantirAdmin(context.Background(), db, "minha-senha-forte")
	if err != nil || !criado || gerada != "" {
		t.Fatalf("criado=%v gerada=%q erro=%v", criado, gerada, err)
	}
	_, hash, _ := hashDoAdmin(t, db)
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte("minha-senha-forte")) != nil {
		t.Error("o hash deveria conferir com a senha configurada")
	}
}

func TestGarantirAdmin_NaoDuplicaNemTrocaASenha(t *testing.T) {
	db := novoBancoDeTeste(t)
	garantirAdmin(context.Background(), db, "primeira-senha-1")
	_, hashAntes, _ := hashDoAdmin(t, db)

	criado, gerada, err := garantirAdmin(context.Background(), db, "outra-senha-qualquer")
	if err != nil || criado || gerada != "" {
		t.Fatalf("segunda chamada não deveria criar nada: criado=%v gerada=%q erro=%v", criado, gerada, err)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&n)
	if _, hashDepois, _ := hashDoAdmin(t, db); n != 1 || hashDepois != hashAntes {
		t.Errorf("deveria continuar 1 usuário com o mesmo hash (n=%d)", n)
	}
}

func TestGarantirAdmin_SenhaConfiguradaCurtaDaErroClaro(t *testing.T) {
	db := novoBancoDeTeste(t)

	_, _, err := garantirAdmin(context.Background(), db, "123")
	if err == nil || !strings.Contains(err.Error(), "ADMIN_SENHA") || !strings.Contains(err.Error(), "mínimo 8") {
		t.Fatalf("esperava erro citando ADMIN_SENHA e o mínimo de 8, veio %v", err)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&n)
	if n != 0 {
		t.Errorf("nada deveria ter sido gravado, há %d usuários", n)
	}
}

func TestGarantirAdmin_SoOsVendedoresNaoImpedemOAdmin(t *testing.T) {
	db := novoBancoDeTeste(t)
	if _, err := db.Exec("INSERT INTO usuarios (nome, email, senha_hash, papel) VALUES ('Vera', 'vera@teste.com', '$2a$10$x', 'vendedor')"); err != nil {
		t.Fatal(err)
	}
	if criado, _, err := garantirAdmin(context.Background(), db, "senha-do-admin-1"); err != nil || !criado {
		t.Fatalf("sem nenhum admin, deveria criar: criado=%v erro=%v", criado, err)
	}
}

func TestGerarSenhaAleatoria_DiferenteACadaVezEValida(t *testing.T) {
	a, err1 := gerarSenhaAleatoria()
	b, err2 := gerarSenhaAleatoria()
	if err1 != nil || err2 != nil || a == b {
		t.Fatalf("duas senhas sorteadas deveriam ser diferentes: %q %q (%v %v)", a, b, err1, err2)
	}
	if len(a) < 12 || len(a) > 72 {
		t.Errorf("tamanho fora do esperado: %d", len(a))
	}
}

func TestSemear_NaoReaproveitaIdsDeUsuariosApagados(t *testing.T) {
	db := novoBancoDeTeste(t)
	for _, e := range []string{"a@teste.com", "b@teste.com"} {
		if _, err := db.Exec("INSERT INTO usuarios (nome, email, senha_hash, papel) VALUES ('X', ?, '$2a$10$x', 'vendedor')", e); err != nil {
			t.Fatal(err)
		}
	}
	db.Exec("DELETE FROM usuarios WHERE email = 'b@teste.com'") // o último id (2) some

	if err := semear(db, time.Now()); err != nil {
		t.Fatal(err)
	}
	res, err := db.Exec("INSERT INTO usuarios (nome, email, senha_hash, papel) VALUES ('Novo', 'novo@teste.com', '$2a$10$y', 'vendedor')")
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := res.LastInsertId(); id != 3 {
		t.Errorf("o novo usuário deveria ter id 3 (o 2 já foi de outra pessoa), veio %d", id)
	}
}

func TestSemear_NaoApagaUsuarios(t *testing.T) {
	db := novoBancoDeTeste(t)
	garantirAdmin(context.Background(), db, "senha-do-admin-1")
	if err := semear(db, time.Now()); err != nil {
		t.Fatal(err)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&n)
	if n != 1 {
		t.Errorf("o seed não deveria apagar usuários, restaram %d", n)
	}
}

func TestRoteador_LigaAsRotasDeUsuariosSemQuebrarAsAntigas(t *testing.T) {
	db := novoBancoDeTeste(t)
	if err := semear(db, time.Now()); err != nil {
		t.Fatal(err)
	}
	r := montarRoteador(db)

	w := chamarAPI(t, r, http.MethodPost, "/usuarios", novoUsuarioJSON("Ana", "ana@teste.com", "senha-forte-1", "vendedor"))
	if w.Code != http.StatusCreated {
		t.Errorf("POST /usuarios no roteador real: esperava 201, veio %d %s", w.Code, w.Body)
	}
	if w := chamarAPI(t, r, http.MethodGet, "/usuarios", ""); w.Code != 200 {
		t.Errorf("GET /usuarios: %d", w.Code)
	}
	for _, rota := range []string{"/clientes", "/funil", "/vendas/resumo", "/vendas"} {
		if w := chamarAPI(t, r, http.MethodGet, rota, ""); w.Code != 200 {
			t.Errorf("GET %s deveria continuar funcionando: %d", rota, w.Code)
		}
	}
}
