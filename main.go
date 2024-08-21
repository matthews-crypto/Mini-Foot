package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Charger les variables d'environnement
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Erreur lors du chargement du fichier .env")
    }

    // Connexion à MongoDB
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Vérifier la connexion
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connecté à MongoDB!")

    // Initialiser le routeur
    r := mux.NewRouter()

    // Définir les routes (à compléter plus tard)
    r.HandleFunc("/", HomeHandler).Methods("GET")

    // Démarrer le serveur
    fmt.Println("Serveur démarré sur le port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Bienvenue sur l'API Mini-Foot!")
}