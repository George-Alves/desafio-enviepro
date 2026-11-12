package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	usuarioHandler "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/handler"
	usuarioRepo "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/repository"
	usuarioUC "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/app/usecase"
	"github.com/gin-gonic/gin"
)

// Estes testes passam pela API de verdade: rotas, handler, casos de uso, repositório e um banco SQLite temporário.

func novoRoteadorDeUsuarios(t *testing.T) (*gin.Engine, *sql.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := abrirBanco(filepath.Join(t.TempDir(), "api-usuarios.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	repo := usuarioRepo.NewSQLiteUsuarioRepository(db)
	h := usuarioHandler.NewUsuarioHandler(
		usuarioUC.NewCriarUsuario(repo),
		usuarioUC.NewListarUsuarios(repo),
		usuarioUC.NewBuscarUsuario(repo),
		usuarioUC.NewAtualizarUsuario(repo),
		usuarioUC.NewApagarUsuario(repo),
	)
	r := gin.New()
	h.Registrar(r)
	return r, db
}

func jsonDe(campos map[string]any) string {
	b, _ := json.Marshal(campos)
	return string(b)
}

func chamarAPI(t *testing.T, r *gin.Engine, metodo, caminho, corpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, caminho, strings.NewReader(corpo))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func mensagemDeErro(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var corpo map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &corpo); err != nil || corpo["erro"] == "" {
		t.Fatalf("esperava {\"erro\": \"...\"}, veio %q", w.Body.String())
	}
	return corpo["erro"]
}

func novoUsuarioJSON(nome, email, senha, papel string) string {
	return jsonDe(map[string]any{"nome": nome, "email": email, "senha": senha, "papel": papel})
}

func criarNaAPI(t *testing.T, r *gin.Engine, nome, email string) int64 {
	t.Helper()
	w := chamarAPI(t, r, http.MethodPost, "/usuarios", novoUsuarioJSON(nome, email, "senha-forte-1", "vendedor"))
	if w.Code != http.StatusCreated {
		t.Fatalf("criar %s: status %d, corpo %s", email, w.Code, w.Body)
	}
	var u struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(w.Body.Bytes(), &u)
	return u.ID
}

func TestAPIUsuarios_CriarDevolve201SemHashEGravaBcrypt(t *testing.T) {
	r, db := novoRoteadorDeUsuarios(t)

	w := chamarAPI(t, r, http.MethodPost, "/usuarios", novoUsuarioJSON("Ana", "Ana@Teste.com", "senha-forte-1", "admin"))
	if w.Code != http.StatusCreated {
		t.Fatalf("esperava 201, veio %d: %s", w.Code, w.Body)
	}

	var u map[string]any
	json.Unmarshal(w.Body.Bytes(), &u)
	if u["nome"] != "Ana" || u["email"] != "ana@teste.com" || u["papel"] != "admin" || u["id"] == nil || u["criado_em"] == nil {
		t.Errorf("resposta inesperada: %v", u)
	}

	// critério: nenhuma resposta JSON contém senha_hash (nem a senha, nem o hash)
	corpo := w.Body.String()
	for _, proibido := range []string{"senha", "hash", "$2a$", "senha-forte-1"} {
		if strings.Contains(corpo, proibido) {
			t.Errorf("a resposta não pode conter %q: %s", proibido, corpo)
		}
	}

	// critério: no banco a coluna senha_hash começa com $2a$, nunca a senha
	var hash string
	if err := db.QueryRow("SELECT senha_hash FROM usuarios WHERE email = 'ana@teste.com'").Scan(&hash); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$2a$") || strings.Contains(hash, "senha-forte-1") {
		t.Errorf("senha_hash deveria começar com $2a$ e não conter a senha, veio %q", hash)
	}
}

func TestAPIUsuarios_ErrosDeEntradaSaoClaros(t *testing.T) {
	r, _ := novoRoteadorDeUsuarios(t)

	casos := []struct {
		nome, corpo, contem string
		status              int
	}{
		{"senha de 5 caracteres", novoUsuarioJSON("Ana", "ana@teste.com", "12345", "admin"), "mínimo 8 caracteres", 400},
		{"e-mail inválido", novoUsuarioJSON("Ana", "ana-sem-arroba", "senha-forte-1", "admin"), "e-mail válido", 400},
		{"papel inexistente", novoUsuarioJSON("Ana", "ana@teste.com", "senha-forte-1", "gerente"), "'admin' ou 'vendedor'", 400},
		{"nome vazio", novoUsuarioJSON("", "ana@teste.com", "senha-forte-1", "admin"), "nome é obrigatório", 400},
		{"JSON quebrado", `{"nome": "Ana",`, "JSON válido", 400},
		{"corpo vazio", ``, "JSON válido", 400},
	}
	for _, c := range casos {
		w := chamarAPI(t, r, http.MethodPost, "/usuarios", c.corpo)
		if w.Code != c.status {
			t.Errorf("%s: esperava %d, veio %d (%s)", c.nome, c.status, w.Code, w.Body)
			continue
		}
		if msg := mensagemDeErro(t, w); !strings.Contains(msg, c.contem) {
			t.Errorf("%s: a mensagem deveria conter %q, veio %q", c.nome, c.contem, msg)
		}
	}

	// vários campos errados de uma vez: todas as mensagens vêm juntas
	w := chamarAPI(t, r, http.MethodPost, "/usuarios", novoUsuarioJSON("", "x", "1", "z"))
	msg := mensagemDeErro(t, w)
	for _, parte := range []string{"nome", "e-mail", "senha", "papel"} {
		if !strings.Contains(msg, parte) {
			t.Errorf("a mensagem deveria citar %q: %q", parte, msg)
		}
	}
}

func TestAPIUsuarios_EmailRepetidoDa409(t *testing.T) {
	r, _ := novoRoteadorDeUsuarios(t)
	criarNaAPI(t, r, "Ana", "ana@teste.com")

	for _, email := range []string{"ana@teste.com", "ANA@Teste.COM"} {
		w := chamarAPI(t, r, http.MethodPost, "/usuarios", novoUsuarioJSON("Outra", email, "senha-forte-1", "vendedor"))
		if w.Code != http.StatusConflict {
			t.Errorf("%s: esperava 409, veio %d", email, w.Code)
		}
		if msg := mensagemDeErro(t, w); !strings.Contains(msg, "já está cadastrado") {
			t.Errorf("mensagem inesperada: %q", msg)
		}
	}
	w := chamarAPI(t, r, http.MethodGet, "/usuarios", "")
	var lista []map[string]any
	json.Unmarshal(w.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Errorf("só o primeiro usuário deveria existir, há %d", len(lista))
	}
}

func TestAPIUsuarios_ListarEBuscar_NuncaMostramHash(t *testing.T) {
	r, _ := novoRoteadorDeUsuarios(t)

	w := chamarAPI(t, r, http.MethodGet, "/usuarios", "")
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Errorf("lista vazia deveria ser [] com 200, veio %d %q", w.Code, w.Body)
	}

	id := criarNaAPI(t, r, "Ana", "ana@teste.com")
	criarNaAPI(t, r, "Bia", "bia@teste.com")

	w = chamarAPI(t, r, http.MethodGet, "/usuarios", "")
	var lista []map[string]any
	json.Unmarshal(w.Body.Bytes(), &lista)
	if w.Code != 200 || len(lista) != 2 || lista[0]["nome"] != "Ana" {
		t.Errorf("listar: %d %s", w.Code, w.Body)
	}

	w2 := chamarAPI(t, r, http.MethodGet, "/usuarios/"+itoa(id), "")
	if w2.Code != 200 || !strings.Contains(w2.Body.String(), `"email":"ana@teste.com"`) {
		t.Errorf("buscar: %d %s", w2.Code, w2.Body)
	}
	for _, corpo := range []string{w.Body.String(), w2.Body.String()} {
		if strings.Contains(corpo, "hash") || strings.Contains(corpo, "$2a$") || strings.Contains(corpo, "senha") {
			t.Errorf("a resposta não pode conter hash nem senha: %s", corpo)
		}
	}
}

func itoa(n int64) string { b, _ := json.Marshal(n); return string(b) }

func TestAPIUsuarios_IdInvalidoEInexistente(t *testing.T) {
	r, _ := novoRoteadorDeUsuarios(t)

	for _, metodo := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		w := chamarAPI(t, r, metodo, "/usuarios/abc", novoUsuarioJSON("A", "a@teste.com", "senha-forte-1", "admin"))
		if w.Code != 400 {
			t.Errorf("%s /usuarios/abc: esperava 400, veio %d", metodo, w.Code)
		}
		w = chamarAPI(t, r, metodo, "/usuarios/9999", novoUsuarioJSON("A", "a@teste.com", "senha-forte-1", "admin"))
		if w.Code != 404 {
			t.Errorf("%s /usuarios/9999: esperava 404, veio %d", metodo, w.Code)
		}
	}
}

func TestAPIUsuarios_AtualizarComESemSenha(t *testing.T) {
	r, db := novoRoteadorDeUsuarios(t)
	id := criarNaAPI(t, r, "Ana", "ana@teste.com")
	hashDe := func() string {
		var h string
		db.QueryRow("SELECT senha_hash FROM usuarios WHERE id = ?", id).Scan(&h)
		return h
	}
	antes := hashDe()
	caminho := "/usuarios/" + itoa(id)

	// sem senha: troca o nome e o papel, o hash continua o mesmo
	w := chamarAPI(t, r, http.MethodPut, caminho, jsonDe(map[string]any{"nome": "Ana Souza", "email": "ana@teste.com", "papel": "admin"}))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"nome":"Ana Souza"`) || !strings.Contains(w.Body.String(), `"papel":"admin"`) {
		t.Fatalf("PUT sem senha: %d %s", w.Code, w.Body)
	}
	if hashDe() != antes {
		t.Error("sem senha no PUT, o hash não pode mudar")
	}

	// com senha: o hash muda
	w = chamarAPI(t, r, http.MethodPut, caminho, novoUsuarioJSON("Ana Souza", "ana@teste.com", "nova-senha-9", "admin"))
	if w.Code != 200 || hashDe() == antes {
		t.Errorf("PUT com senha: %d, e o hash deveria ter mudado", w.Code)
	}
	if strings.Contains(w.Body.String(), "hash") || strings.Contains(w.Body.String(), "$2a$") {
		t.Errorf("a resposta do PUT não pode conter hash: %s", w.Body)
	}

	// senha curta e e-mail de outro usuário
	if w := chamarAPI(t, r, http.MethodPut, caminho, novoUsuarioJSON("Ana", "ana@teste.com", "123", "admin")); w.Code != 400 {
		t.Errorf("senha curta no PUT: esperava 400, veio %d", w.Code)
	}
	criarNaAPI(t, r, "Bia", "bia@teste.com")
	if w := chamarAPI(t, r, http.MethodPut, caminho, jsonDe(map[string]any{"nome": "Ana", "email": "BIA@teste.com", "papel": "admin"})); w.Code != 409 {
		t.Errorf("e-mail de outro no PUT: esperava 409, veio %d", w.Code)
	}
}

func TestAPIUsuarios_ApagarDa204EDepoisSome(t *testing.T) {
	r, _ := novoRoteadorDeUsuarios(t)
	id := criarNaAPI(t, r, "Ana", "ana@teste.com")
	caminho := "/usuarios/" + itoa(id)

	w := chamarAPI(t, r, http.MethodDelete, caminho, "")
	if w.Code != http.StatusNoContent || w.Body.Len() != 0 {
		t.Errorf("DELETE: esperava 204 sem corpo, veio %d %q", w.Code, w.Body)
	}
	if w := chamarAPI(t, r, http.MethodGet, caminho, ""); w.Code != 404 {
		t.Errorf("depois de apagar, GET deveria dar 404, veio %d", w.Code)
	}
}

func TestAPIUsuarios_CamposExtrasDoJSONSaoIgnorados(t *testing.T) {
	r, db := novoRoteadorDeUsuarios(t)

	corpo := jsonDe(map[string]any{
		"nome": "Ana", "email": "ana@teste.com", "senha": "senha-forte-1", "papel": "vendedor",
		"id": 999, "senha_hash": "$2a$10$hash-injetado-pelo-cliente", "criado_em": "2000-01-01T00:00:00Z",
	})
	w := chamarAPI(t, r, http.MethodPost, "/usuarios", corpo)
	if w.Code != 201 {
		t.Fatalf("esperava 201, veio %d: %s", w.Code, w.Body)
	}

	var id int64
	var hash, criadoEm string
	db.QueryRow("SELECT id, senha_hash, criado_em FROM usuarios WHERE email = 'ana@teste.com'").Scan(&id, &hash, &criadoEm)
	if id == 999 || hash == "$2a$10$hash-injetado-pelo-cliente" || strings.HasPrefix(criadoEm, "2000") {
		t.Errorf("o cliente não pode escolher id, hash nem data: id=%d hash=%q criado_em=%q", id, hash, criadoEm)
	}
}
