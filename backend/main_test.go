package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func novoServidorDeTeste(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := abrirBanco(filepath.Join(t.TempDir(), "teste.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if err := semear(db, time.Now()); err != nil {
		t.Fatal(err)
	}
	return montarRoteador(db)
}

func get(t *testing.T, r *gin.Engine, caminho string, destino any) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, caminho, nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET %s: status %d, corpo %s", caminho, w.Code, w.Body)
	}
	if err := json.Unmarshal(w.Body.Bytes(), destino); err != nil {
		t.Fatalf("GET %s: JSON inválido: %v", caminho, err)
	}
	return w
}

func TestAPI_ClientesVemDoBancoComAlertas(t *testing.T) {
	var clientes []struct {
		Nome           string `json:"nome"`
		DiasSemContato int    `json:"dias_sem_contato"`
		Alerta         bool   `json:"alerta"`
	}
	get(t, novoServidorDeTeste(t), "/clientes", &clientes)

	if len(clientes) != 20 {
		t.Fatalf("esperava 20 clientes, vieram %d", len(clientes))
	}
	comAlerta := 0
	for _, c := range clientes {
		if c.Alerta != (c.DiasSemContato > 15) {
			t.Errorf("%s: alerta=%v com %d dias", c.Nome, c.Alerta, c.DiasSemContato)
		}
		if c.Alerta {
			comAlerta++
		}
	}
	if comAlerta == 0 || comAlerta == len(clientes) {
		t.Errorf("o seed deveria ter clientes com e sem alerta, mas %d de %d têm", comAlerta, len(clientes))
	}
}

func TestAPI_FunilTemTodasAsEtapasESomaAs30Oportunidades(t *testing.T) {
	var etapas []struct {
		Etapa      string  `json:"etapa"`
		Quantidade int     `json:"quantidade"`
		Valor      float64 `json:"valor"`
	}
	get(t, novoServidorDeTeste(t), "/funil", &etapas)

	if len(etapas) != 5 || etapas[0].Etapa != "Lead" || etapas[4].Etapa != "Fechado" {
		t.Fatalf("etapas inesperadas: %+v", etapas)
	}
	total := 0
	for _, e := range etapas {
		total += e.Quantidade
	}
	if total != 30 {
		t.Errorf("esperava 30 oportunidades no funil, vieram %d", total)
	}
}

func TestAPI_VendasEResumoDoMes(t *testing.T) {
	r := novoServidorDeTeste(t)

	var vendas []struct {
		Data  string  `json:"data"`
		Valor float64 `json:"valor"`
	}
	get(t, r, "/vendas", &vendas)
	if len(vendas) != 25 {
		t.Fatalf("esperava 25 vendas, vieram %d", len(vendas))
	}

	mes := time.Now().Format("2006-01")
	somaMes, qtdMes := 0.0, 0
	for _, v := range vendas {
		if len(v.Data) >= 7 && v.Data[:7] == mes {
			somaMes += v.Valor
			qtdMes++
		}
	}

	var resumo struct {
		Total       float64 `json:"total"`
		Quantidade  int     `json:"quantidade"`
		TicketMedio float64 `json:"ticket_medio"`
	}
	get(t, r, "/vendas/resumo", &resumo)

	if resumo.Quantidade != qtdMes || resumo.Total != somaMes || qtdMes == 0 {
		t.Errorf("resumo %+v não bate com as vendas do mês (qtd %d, soma %.2f)", resumo, qtdMes, somaMes)
	}
}

func TestAPI_CORSLiberaOFrontELogoNaoLiberaOutrasOrigens(t *testing.T) {
	r := novoServidorDeTeste(t)
	pedir := func(origem string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/funil", nil)
		req.Header.Set("Origin", origem)
		r.ServeHTTP(w, req)
		return w
	}
	for _, origem := range []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000"} {
		if got := pedir(origem).Header().Get("Access-Control-Allow-Origin"); got != origem {
			t.Errorf("%s deveria ser liberada, cabeçalho veio %q", origem, got)
		}
	}
	for _, origem := range []string{"http://site-qualquer.com", "http://localhost.evil.com", "https://localhost:3000"} {
		if got := pedir(origem).Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("%s não deveria ser liberada, cabeçalho veio %q", origem, got)
		}
	}
}

func TestSemear_PodeRodarDeNovoSemDuplicar(t *testing.T) {
	db, err := abrirBanco(filepath.Join(t.TempDir(), "outro.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for i := 0; i < 2; i++ {
		if err := semear(db, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM clientes").Scan(&n); err != nil || n != 20 {
		t.Fatalf("depois de semear duas vezes esperava 20 clientes, veio %d (erro %v)", n, err)
	}
}
