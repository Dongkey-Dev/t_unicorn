package mock

import (
	"fmt"
	"math/rand"
	"strconv"
	. "t_unicorn/models"
	"time"
)

type MockUser struct {
	UserID        int      `json:"userid"`
	UserName      string   `json:"username"`
	UserPswd      string   `json:"userpswd"`
	UserEmail     string   `json:"useremail"`
	UserCreatedOn string   `json:"usercreatedon"`
	UserInfo      UserInfo `json:"userunfo"`
}

func getRandomPhoneNumber() []string {
	s_i := strconv.Itoa(rand.Intn(99999999))
	return []string{s_i[:4], s_i[4:]}
}

func GetMockUser(num int) []MockUser {
	var MockUsers []MockUser
	for i := 1; i < num; i++ {
		phone := getRandomPhoneNumber()
		UserInfo := new(UserInfo)
		UserInfo.DOB = "1995-12-01"
		UserInfo.Name = fmt.Sprintf("DongkeyDev_%d", i)
		UserInfo.Gender = "male"
		UserInfo.Phone = fmt.Sprintf("010-%s-%s", phone[0], phone[1])
		User := MockUser{
			UserName:      fmt.Sprintf("test_%d", i),
			UserPswd:      "testtest",
			UserEmail:     fmt.Sprintf("test_%d@gmail.com", i),
			UserCreatedOn: time.Now().String(),
			UserInfo:      *UserInfo,
		}
		MockUsers = append(MockUsers, User)
	}
	return MockUsers
}
