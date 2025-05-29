package dto

type OnboardingDTO struct {
	Fullname    string `json:"fullname"`
	Bio         string `json:"bio"`
	NativeLng   string `json:"native_lng"`
	LearningLng string `json:"learning_lng"`
	Location    string `json:"location"`
	ProfilePic  string `json:"profile_pic"`
}
