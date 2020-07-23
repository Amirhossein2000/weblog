package schema

type LoginRequest struct {
	Email    string
	Password string
}

type TimelineRequest struct {
	UserId uint
	Limit  int
	Offset int
}
