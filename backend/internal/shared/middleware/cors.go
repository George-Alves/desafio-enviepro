package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// OrigemLocal diz se a origem é uma página aberta neste mesmo computador (localhost ou 127.0.0.1,
// em qualquer porta). O Next.js troca de porta sozinho quando a 3000 está ocupada.
func OrigemLocal(origem string) bool {
	u, err := url.Parse(origem)
	if err != nil || u.Scheme != "http" {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

// CORS libera o front-end a chamar a API pelo navegador, para as origens aceitas por permitida.
// Sem isso o navegador bloqueia as respostas do back-end.
func CORS(permitida func(origem string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origem := c.GetHeader("Origin")
		if origem != "" && permitida(origem) {
			c.Header("Access-Control-Allow-Origin", origem)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
