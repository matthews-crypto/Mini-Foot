package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nom             string             `bson:"nom" json:"nom"`
	Prenom          string             `bson:"prenom" json:"prenom"`
	Telephone       string             `bson:"telephone" json:"telephone"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	MotDePasse      string             `bson:"mot_de_passe" json:"mot_de_passe,omitempty"`
	TypeUtilisateur string             `bson:"type_utilisateur" json:"type_utilisateur"`
	ImageProfil     string             `bson:"image_profil,omitempty" json:"image_profil,omitempty"`
	DateCreation    time.Time          `bson:"date_creation" json:"date_creation"`
}
