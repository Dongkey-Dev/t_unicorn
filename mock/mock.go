package mock

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func getRandomPhoneNumber() string {
	s_i := strconv.Itoa(rand.Intn(89999999) + 10000000)
	fmt.Printf(s_i + " \n")
	return s_i
}

func strToList(str string) []string {
	return []string{str}
}

func mapToRequest(r *http.Request, M map[string]string) *http.Request {
	r.Form = map[string][]string{}
	for k, v := range M {
		r.Form[k] = strToList(v)
	}
	return r
}

func GetMockUser(r *http.Request, num int) []*http.Request {
	twoGender := []string{"male", "female"}
	var mockUsers []*http.Request
	for i := 1; i < num; i++ {
		r_c := r.Clone(r.Context())
		var statementMap map[string]string
		statementMap = map[string]string{}
		statementMap["phone"] = getRandomPhoneNumber()
		statementMap["dob"] = "1995-12-01"
		statementMap["name"] = fmt.Sprintf("DongkeyDev_%d", i)
		statementMap["gender"] = twoGender[rand.Intn(len(twoGender))]
		statementMap["username"] = fmt.Sprintf("test_%d", i)
		statementMap["userpswd"] = "testtest"
		statementMap["email"] = fmt.Sprintf("test_%d@gmail.com", i)
		userRequest := mapToRequest(r_c, statementMap)
		mockUsers = append(mockUsers, userRequest)
	}
	return mockUsers
}
