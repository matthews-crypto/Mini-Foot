package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/matthews-crypto/Mini-Foot/api/handlers"
	"github.com/matthews-crypto/Mini-Foot/internal/user"
	"github.com/matthews-crypto/Mini-Foot/pkg/auth"
	"github.com/matthews-crypto/Mini-Foot/pkg/middleware"
)

func main() {
	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("Avertissement: Fichier .env non trouvé. Utilisation des variables d'environnement système.")
	}

	// Vérifier et initialiser la clé JWT
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.Fatal("La variable d'environnement JWT_SECRET_KEY n'est pas définie")
	}
	auth.InitJWTKey(jwtKey)

	// Connexion à MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Vérifier la connexion
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Connecté à MongoDB!")

	// Initialiser les dépendances
	db := client.Database("Mini-Foot")
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Configurer le routeur
	r := mux.NewRouter()

	// Servir les fichiers statiques
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Routes pour les pages HTML
	r.HandleFunc("/", serveTemplate("index.html"))
	r.HandleFunc("/register", serveTemplate("register.html"))
	r.HandleFunc("/login", serveTemplate("login.html"))
	r.HandleFunc("/profile", middleware.AuthMiddleware(serveTemplate("profile.html")))

	// Routes API pour les utilisateurs
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/profile", middleware.AuthMiddleware(userHandler.GetProfile)).Methods("GET")
	r.HandleFunc("/api/profile", middleware.AuthMiddleware(userHandler.UpdateProfile)).Methods("PUT")
	r.HandleFunc("/api/account", middleware.AuthMiddleware(userHandler.DeleteAccount)).Methods("DELETE")

	// Démarrer le serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serveur démarré sur le port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func serveTemplate(tmpl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{
			filepath.Join("web", "templates", "layout.html"),
			filepath.Join("web", "templates", tmpl),
		}

		data := map[string]interface{}{
			"Title": strings.TrimSuffix(tmpl, ".html"),
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Printf("Erreur lors du parsing du template %s : %v", tmpl, err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		if err := ts.ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("Erreur lors de l'exécution du template %s : %v", tmpl, err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}
}
