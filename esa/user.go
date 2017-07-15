package esa

const (
	UserURL = "/v1/user"
)

type UserService struct {
	client *Client
}

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Icon       string `json:"icon"`
	Email      string `json:"email"`
}

func (s *UserService) Get() (User, error) {
	user := User{}

	res, err := s.client.get(UserURL, nil, &user)
	if err != nil { return user, err }
	res.Body.Close()

	return user, nil
}
