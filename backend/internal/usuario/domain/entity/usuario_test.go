package entity

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestUsuario_SenhaHashNuncaVaiParaOJSON(t *testing.T) {
	b, err := json.Marshal(Usuario{ID: 1, Nome: "Ana", Email: "ana@teste.com", SenhaHash: "$2a$10$segredo", Papel: PapelAdmin})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "segredo") || strings.Contains(strings.ToLower(string(b)), "senha") {
		t.Errorf("o hash não pode aparecer no JSON da entidade: %s", b)
	}
}

func TestNormalizarEmail_TiraEspacosEPoeEmMinusculas(t *testing.T) {
	if got := NormalizarEmail("  Ana@Teste.COM \n"); got != "ana@teste.com" {
		t.Errorf("veio %q", got)
	}
}

func TestValidarEmail(t *testing.T) {
	validos := []string{"ana@teste.com", "a.b+c@sub.dominio.com.br", "  ANA@Teste.com "}
	for _, e := range validos {
		if err := ValidarEmail(e); err != nil {
			t.Errorf("%q deveria ser válido, deu erro: %v", e, err)
		}
	}

	invalidos := []string{
		"", "   ", "ana", "ana@", "@teste.com", "ana@teste", "ana@.com", "ana@teste.",
		"Ana <ana@teste.com>", "ana teste@x.com", strings.Repeat("a", 250) + "@x.com",
	}
	for _, e := range invalidos {
		if err := ValidarEmail(e); err == nil {
			t.Errorf("%q deveria ser inválido", e)
		}
	}
}

func TestValidarSenha(t *testing.T) {
	// 5 caracteres: o caso do critério de aceite, com mensagem clara
	err := ValidarSenha("12345")
	if err == nil || !strings.Contains(err.Error(), "mínimo 8") {
		t.Errorf("senha de 5 caracteres deveria falar do mínimo de 8, veio %v", err)
	}

	if err := ValidarSenha("12345678"); err != nil {
		t.Errorf("8 caracteres deveria valer: %v", err)
	}
	if err := ValidarSenha(strings.Repeat("a", 72)); err != nil {
		t.Errorf("72 bytes deveria valer: %v", err)
	}

	err = ValidarSenha(strings.Repeat("a", 73))
	if err == nil || !strings.Contains(err.Error(), "72") {
		t.Errorf("73 bytes deveria falar do limite de 72, veio %v", err)
	}

	// o mínimo conta caracteres, e não bytes: 8 letras com acento (16 bytes) valem; 7 não
	if err := ValidarSenha("éééééééé"); err != nil {
		t.Errorf("8 letras acentuadas deveriam valer: %v", err)
	}
	if err := ValidarSenha("ééééééé"); err == nil {
		t.Error("7 letras acentuadas não deveriam valer")
	}
}

func TestValidarPapel(t *testing.T) {
	for _, p := range []string{PapelAdmin, PapelVendedor} {
		if err := ValidarPapel(p); err != nil {
			t.Errorf("%q deveria valer: %v", p, err)
		}
	}
	for _, p := range []string{"", "gerente", "Admin", "ADMIN", " admin"} {
		if err := ValidarPapel(p); err == nil {
			t.Errorf("%q não deveria valer", p)
		}
	}
}

func TestValidarNome(t *testing.T) {
	if err := ValidarNome("Ana"); err != nil {
		t.Errorf("Ana deveria valer: %v", err)
	}
	for _, n := range []string{"", "   ", strings.Repeat("a", 101)} {
		if err := ValidarNome(n); err == nil {
			t.Errorf("%q não deveria valer", n)
		}
	}
}

func TestValidarNovoUsuario_JuntaTodosOsErros(t *testing.T) {
	if err := ValidarNovoUsuario("Ana", "ana@teste.com", "senha-forte", PapelVendedor); err != nil {
		t.Fatalf("dados bons não deveriam dar erro: %v", err)
	}

	err := ValidarNovoUsuario("", "isso-nao-e-email", "123", "chefe")
	var erros ErrosValidacao
	if !errors.As(err, &erros) {
		t.Fatalf("esperava ErrosValidacao, veio %T", err)
	}
	campos := map[string]bool{}
	for _, e := range erros {
		campos[e.Campo] = true
	}
	for _, c := range []string{"nome", "email", "senha", "papel"} {
		if !campos[c] {
			t.Errorf("faltou o erro do campo %q: %v", c, erros)
		}
	}
}

func TestValidarAtualizacao_SenhaOpcional(t *testing.T) {
	if err := ValidarAtualizacao("Ana", "ana@teste.com", "", PapelAdmin); err != nil {
		t.Errorf("senha vazia deveria significar 'manter a atual': %v", err)
	}
	if err := ValidarAtualizacao("Ana", "ana@teste.com", "novasenha1", PapelAdmin); err != nil {
		t.Errorf("senha nova válida deveria passar: %v", err)
	}
	err := ValidarAtualizacao("Ana", "ana@teste.com", "123", PapelAdmin)
	var erros ErrosValidacao
	if !errors.As(err, &erros) || len(erros) != 1 || erros[0].Campo != "senha" {
		t.Errorf("senha curta deveria dar só o erro de senha, veio %v", err)
	}
}
