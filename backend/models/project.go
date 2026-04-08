package models

type Project struct {
	BaseModel    `bson:",inline"`
	Name         string        `json:"name" bson:"name" validate:"required,min=1,max=100"`
	Key          string        `json:"key" bson:"key" validate:"required,uppercase,min=2,max=10"`
	Description  string        `json:"description" bson:"description"`
	WorkflowID   string        `json:"workflow_id" bson:"workflow_id"`
	CustomFields []CustomField `json:"custom_fields" bson:"custom_fields"`
	Members      []string      `json:"members" bson:"members"`
}

type CreateProjectRequest struct {
	Name         string        `json:"name" validate:"required,min=1,max=100"`
	Key          string        `json:"key" validate:"required,uppercase,min=2,max=10"`
	Description  string        `json:"description"`
	CustomFields []CustomField `json:"custom_fields"`
}

type UpdateProjectRequest struct {
	Name         *string        `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description  *string        `json:"description,omitempty"`
	CustomFields *[]CustomField `json:"custom_fields,omitempty"`
}
