package user

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, user *User) error {
	// Vérifier si l'utilisateur existe déjà
	existingUser, _ := s.repo.GetByTelephone(ctx, user.Telephone)
	if existingUser != nil {
		return errors.New("un utilisateur avec ce numéro de téléphone existe déjà")
	}

	log.Printf("Mot de passe reçu pour l'inscription : %s", user.MotDePasse)

	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erreur lors du hachage du mot de passe : %v", err)
		return err
	}
	user.MotDePasse = string(hashedPassword)

	log.Printf("Mot de passe haché pour l'inscription : %s", user.MotDePasse)

	// Créer l'utilisateur
	err = s.repo.Create(ctx, user)
	if err != nil {
		log.Printf("Erreur lors de la création de l'utilisateur : %v", err)
		return err
	}

	log.Printf("Utilisateur créé avec succès : %s", user.Telephone)
	return nil
}

func (s *Service) Authenticate(ctx context.Context, telephone, password string) (*User, error) {
	user, err := s.repo.GetByTelephone(ctx, telephone)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur : %v", err)
		return nil, errors.New("utilisateur non trouvé")
	}

	log.Printf("Mot de passe haché stocké : %s", user.MotDePasse)
	log.Printf("Mot de passe fourni : %s", password)

	err = bcrypt.CompareHashAndPassword([]byte(user.MotDePasse), []byte(password))
	if err != nil {
		log.Printf("Erreur lors de la comparaison des mots de passe : %v", err)
		return nil, errors.New("mot de passe incorrect")
	}

	log.Printf("Authentification réussie pour l'utilisateur : %s", user.Telephone)
	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id primitive.ObjectID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
