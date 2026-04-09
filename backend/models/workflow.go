package models

type Workflow struct {
	BaseModel   `bson:",inline"`
	Name        string       `json:"name" bson:"name" validate:"required"`
	ProjectID   string       `json:"project_id" bson:"project_id"`
	Statuses    []string     `json:"statuses" bson:"statuses" validate:"required,min=1"`
	Transitions []Transition `json:"transitions" bson:"transitions"`
}

type Transition struct {
	From       string      `json:"from" bson:"from" validate:"required"`
	To         string      `json:"to" bson:"to" validate:"required"`
	Conditions []Condition `json:"conditions,omitempty" bson:"conditions,omitempty"`
	Actions    []Action    `json:"actions,omitempty" bson:"actions,omitempty"`
}

type Condition struct {
	Field    string `json:"field" bson:"field" validate:"required"`
	Operator string `json:"operator" bson:"operator" validate:"required"` // not_empty, equals, greater_than
	Value    string `json:"value,omitempty" bson:"value,omitempty"`
}

type Action struct {
	Type   string                 `json:"type" bson:"type" validate:"required"` // assign_reviewer, notify, set_field
	Params map[string]interface{} `json:"params,omitempty" bson:"params,omitempty"`
}
