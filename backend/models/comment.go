package models

type Comment struct {
	BaseModel `bson:",inline"`
	IssueID   string   `json:"issue_id" bson:"issue_id" validate:"required"`
	UserID    string   `json:"user_id" bson:"user_id" validate:"required"`
	Content   string   `json:"content" bson:"content" validate:"required,min=1"`
	Mentions  []string `json:"mentions" bson:"mentions"`
	ParentID  string   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
}

type CommentWithUser struct {
	Comment `bson:",inline"`
	User    *UserRef `json:"user" bson:"-"`
}

type CreateCommentRequest struct {
	Content  string `json:"content" validate:"required,min=1"`
	ParentID string `json:"parent_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}
