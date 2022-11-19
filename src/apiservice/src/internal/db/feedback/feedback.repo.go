package feedback

import "context"

type FeedbackRepo interface {
	SelectFeedback(context.Context, *Search) (*Feedback, error)
	InsertFeedback(context.Context, *Feedback) (*Feedback, error)
	UpdateFeedback(context.Context, *Search, *Feedback) error
	SelectRatingByBussiness(context.Context, *Search) (*Feedback, error)
	ListRatingByBussiness(context.Context, *Search) ([]*Feedback, error)
	ListFeedbacksByBusiness(context.Context, *Search) ([]*Feedback, error)
	TotalFeedbacksByBusiness(context.Context, *Search) (*int64, error)
}
