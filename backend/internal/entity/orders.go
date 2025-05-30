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

func DefaultOrder(
	clientID primitive.ObjectID,
	title, description string,
	completionTime int64,
	cost int,
	createdAt, updatedAt time.Time,
) OrderExt {
	return OrderExt{
		Order: Order{
			ClientID:       clientID,
			Title:          title,
			Description:    description,
			Cost:           cost,
			CompletionTime: completionTime,
			Active:         true,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		},
		Responses: []Response{},
		Statuses: []Status{
			{
				Type:      StatusTypeBeginning,
				CreatedAt: createdAt,
			},
		},
	}
}
