package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/matthews-crypto/Mini-Foot/pkg/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Début de la vérification du token")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Aucun en-tête d'autorisation trouvé")
			http.Error(w, "Autorisation manquante", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			log.Println("AuthMiddleware: Format d'autorisation invalide")
			http.Error(w, "Format d'autorisation invalide", http.StatusUnauthorized)
			return
		}

		log.Println("AuthMiddleware: Validation du token")
		claims, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			log.Printf("AuthMiddleware: Erreur de validation du token: %v", err)
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}

		log.Printf("AuthMiddleware: Token valide pour l'utilisateur ID: %s", claims.UserID)
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
