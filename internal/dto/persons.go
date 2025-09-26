package dto

type CreatePersonRequest struct {
	Name    string  `json:"name" binding:"required"`
	Age     *int    `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

type UpdatePersonRequest struct {
	//Name    *string `json:"name,omitempty"`
	Age     *int    `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

type PersonResponse struct {
	ID      uint64  `json:"id"`
	Name    string  `json:"name"`
	Age     *int    `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}
