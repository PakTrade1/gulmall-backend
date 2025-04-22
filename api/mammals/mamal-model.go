package mammals

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the database
type User struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IsEmailVerified bool               `json:"isEmailVerified" bson:"isEmailVerified"`
	DisplayName     string             `json:"displayName" bson:"displayName"`
	Email           string             `json:"email" bson:"email"`
	PhotoURL        string             `json:"photoUrl" bson:"photoUrl"`
	PublicID        int64              `json:"publicId" bson:"publicId"`
	Credit          int64              `json:"credit" bson:"credit"`
	AdsRemaining    int                `json:"adsRemaining" bson:"adsRemaining"`
	CreationDate    string             `json:"creationDate" bson:"creationDate"`
	ServerDate      primitive.DateTime `json:"serverDate" bson:"serverDate"`
	LastSignedIn    string             `json:"lastSignedIn" bson:"lastSignedIn"`
	PlanID          primitive.ObjectID `json:"planId,omitempty" bson:"planId,omitempty"`
	ProviderInfo    string             `json:"providerInfo" bson:"providerInfo"`
	BusinessPhone   string             `json:"businessPhone" bson:"businessPhone"`
	AccountStatus   bool               `json:"accountStatus" bson:"accountStatus"`
	IsBusiness      bool               `json:"isBusiness" bson:"isBusiness"`
	PrimaryPhone    string             `json:"primaryPhone" bson:"primaryPhone"`
	UID             string             `json:"uid" bson:"uid"`
	Platform        string             `json:"platform" bson:"platform"`
	Ip              string             `json:"ip" bson:"ip"`
}
