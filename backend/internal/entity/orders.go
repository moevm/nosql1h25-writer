package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ClientID       primitive.ObjectID `bson:"clientId"`
	Title          string             `bson:"title"`
	Description    string             `bson:"description"`
	CompletionTime int64              `bson:"completionTime"`
	Cost           int                `bson:"cost,omitempty"`
	Active         bool               `bson:"active"`
	FreelancerID   primitive.ObjectID `bson:"freelancerId,omitempty"`
	Budget         int                `bson:"budget,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type OrderExt struct {
	Order     `bson:",inline"`
	Responses []Response `bson:"responses"`
	Statuses  []Status   `bson:"statuses"`
}

type Response struct {
	FreelancerName string             `bson:"freelancerName"`
	FreelancerID   primitive.ObjectID `bson:"freelancerId"`
	ChatID         primitive.ObjectID `bson:"chatId"`
	CoverLetter    string             `bson:"coverLetter"`
	Active         bool               `bson:"active"`
	CreatedAt      time.Time          `bson:"createdAt"`
}

type Status struct {
	Type      StatusType `bson:"type"`
	Content   string     `bson:"content,omitempty"`
	CreatedAt time.Time  `bson:"createdAt"`
}

type StatusType string

const (
	StatusTypeBeginning   StatusType = "beginning"
	StatusTypeNegotiation StatusType = "negotiation"
	StatusTypeBudgeting   StatusType = "budgeting"
	StatusTypeWork        StatusType = "work"
	StatusTypeReviews     StatusType = "reviews"
	StatusTypeFinished    StatusType = "finished"
	StatusTypeDispute     StatusType = "dispute"
)
