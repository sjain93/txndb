package user_test

import (
	"fmt"
	"testing"

	"github.com/sjain93/userservice/api/common"
	"github.com/sjain93/userservice/api/user"
	"github.com/sjain93/userservice/config"
	"github.com/stretchr/testify/assert"
)

// Testing Tables: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests

/*
	Test setup using in memory store
	Benefit: persistence is tied to test func context
	Disadvantage: does not test against actual DB

	Production grade services should be tested against an instance of the DB.
*/

func TestCreateUser(t *testing.T) {
	memStore := config.GetInMemoryStore()
	uR, err := user.NewUserRepository(nil, memStore)
	if err != nil {
		t.Fatal(err)
	}

	uS := user.NewUserService(uR)

	//Setup existing users
	existingUsers := getTestUsers(2)
	for _, u := range existingUsers {
		err := uR.Create(&u)
		if err != nil {
			t.Fatal(err)
		}
	}

	testcases := map[string]struct {
		user          user.User
		expectedError error
	}{
		"Happy Path - Success": {
			user: user.User{
				Username: "Test Case 1",
				Email:    "testcase1@gmail.com",
			},
			expectedError: nil,
		},
		"Email in use - failure": {
			user: user.User{
				Username: "Test Case 2",
				Email:    "existinguser1@gmail.com",
			},
			expectedError: user.ErrSvcUserExists,
		},
		"Username in use - failure": {
			user: user.User{
				Username: "Existing User 2",
				Email:    "testcase3@gmail.com",
			},
			expectedError: user.ErrSvcUserExists,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			_, err := uS.CreateUser(tc.user)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func getTestUsers(num int) []user.User {
	generatedUsers := []user.User{}

	for i := 1; i <= num; i++ {
		uName := fmt.Sprintf("Existing User %v", i)
		email := fmt.Sprintf("existinguser%v@gmail.com", i)
		u := user.User{
			ID:       common.GetMD5HashWithSum(uName + email),
			Username: uName,
			Email:    email,
		}
		generatedUsers = append(generatedUsers, u)
	}

	return generatedUsers
}
