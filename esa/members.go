package esa

import "net/url"

const (
// MembersURL esa API のメンバーのベ-スURL
	MembersURL = "/v1/teams"
)

// MembersService API docs: https://docs.esa.io/posts/102#6-0-0
type MembersService struct {
	client *Client
}

// Member メンバー情報
type Member struct {
	Email      string `json:"email"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// MembersResponse メンバー情報のレスポンス
type MembersResponse struct {
	Members []Member `json:"members"`
	NextPage   interface{} `json:"next_page"`
	PrevPage   interface{} `json:"prev_page"`
	TotalCount int         `json:"total_count"`
}

// GetTeamMembers チ-ム名を指定してメンバー情報を取得する
func (s *MembersService) Get(teamName string) ([]Member, error) {
	membersURL	:= MembersURL+ "/" + teamName + "/members"
	members		:= []Member{}
	page		:= "1"

	for {
		query := url.Values{}

		query.Add("page", page)
		query.Add("per_page", "100")

		var membersRes MembersResponse

		res, err := s.client.get(membersURL, query, &membersRes)
		if err != nil {
			return members, err
		}
		res.Body.Close()

		members = append(members, membersRes.Members...)

		if next_page, ok := membersRes.NextPage.(string); ok {
			page = next_page
		} else {
			break
		}
	}

	return members, nil
}
