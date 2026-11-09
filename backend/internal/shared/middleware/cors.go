package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS libera o front-end (origens permitidas) a chamar a API pelo navegador.
// Sem isso o navegador bloqueia as respostas do back-end.
func CORS(origensPermitidas ...string) gin.HandlerFunc {
	permitidas := make(map[string]bool, len(origensPermitidas))
	for _, o := range origensPermitidas {
		permitidas[o] = true
	}
	return func(c *gin.Context) {
		origem := c.GetHeader("Origin")
		if permitidas[origem] {
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
