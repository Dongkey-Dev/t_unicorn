package models

type User struct {
	UserID        int    `json:"userid"`
	UserName      string `json:"username"`
	UserEmail     string `json:"useremail"`
	UserCreatedOn string `json:"usercreatedon"`
}

type UserInfo struct {
	UserID int    `json:"userid"`
	Name   string `json:"name"`
	DOB    string `json:"dob"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
}

type UserSequenceID struct {
	nextval int `json:"userid"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

type JsonResponseSequenceID struct {
	Type    string           `json:"type"`
	Data    []UserSequenceID `json:"data"`
	Message string           `json:"message"`
}
