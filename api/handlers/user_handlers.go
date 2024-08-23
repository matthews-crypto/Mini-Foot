package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/matthews-crypto/Mini-Foot/internal/user"
	"github.com/matthews-crypto/Mini-Foot/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("Erreur lors du décodage des informations d'inscription : %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Tentative d'inscription pour le téléphone : %s", u.Telephone)
	log.Printf("Mot de passe reçu dans le handler : %s", u.MotDePasse)

	if err := h.userService.Register(r.Context(), &u); err != nil {
		log.Printf("Erreur lors de l'inscription : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Inscription réussie pour l'utilisateur : %s", u.Telephone)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Utilisateur créé avec succès"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Telephone  string `json:"telephone"`
		MotDePasse string `json:"mot_de_passe"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erreur lors de la lecture du corps de la requête : %v", err)
		http.Error(w, "Erreur lors de la lecture de la requête", http.StatusBadRequest)
		return
	}
	log.Printf("Corps de la requête reçu : %s", string(body))

	if err := json.Unmarshal(body, &creds); err != nil {
		log.Printf("Erreur lors du décodage des informations de connexion : %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Tentative de connexion pour le téléphone : %s", creds.Telephone)
	log.Printf("Mot de passe reçu (longueur) : %d", len(creds.MotDePasse))

	u, err := h.userService.Authenticate(r.Context(), creds.Telephone, creds.MotDePasse)
	if err != nil {
		log.Printf("Échec de l'authentification : %v", err)
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(u.ID.Hex())
	if err != nil {
		log.Printf("Erreur lors de la génération du token : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Connexion réussie pour l'utilisateur : %s", u.ID.Hex())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = userID

	if err := h.userService.UpdateUser(r.Context(), &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profil mis à jour avec succès"})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Compte supprimé avec succès"})
}
