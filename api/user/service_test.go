package user_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sjain93/userservice/api/common"
	"github.com/sjain93/userservice/api/user"
	"github.com/sjain93/userservice/config"
	"github.com/sjain93/userservice/migrations"
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
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	config.ConnectDatabase()
	migrations.AutoMigrate(config.DB)

	uR, err := user.NewUserRepository(config.DB)
	if err != nil {
		t.Fatal(err)
	}

	uS := user.NewUserService(uR)

	cleanup := []user.User{}

	// Setup existing users
	existingUsers := getTestUsers(2)
	for _, u := range existingUsers {
		err := uR.Create(&u)
		if err != nil {
			t.Fatal(err)
		}
		cleanup = append(cleanup, u)
	}

	testcases := map[string]struct {
		user          user.User
		expectedError error
	}{
		"Happy Path - Success": {
			user: user.User{
				Username:      "Test Case 1",
				AccountNumber: "1234567",
				Email:         "testcase1@gmail.com",
			},
			expectedError: nil,
		},
		"Email in use - failure": {
			user: user.User{
				Username:      "Test Case 2",
				AccountNumber: "1234568",
				Email:         "existinguser1@gmail.com",
			},
			expectedError: user.ErrSvcUserExists,
		},
		"Username in use - failure": {
			user: user.User{
				Username:      "Existing User 2",
				AccountNumber: "1234569",
				Email:         "testcase3@gmail.com",
			},
			expectedError: user.ErrSvcUserExists,
		},
		"Account number exists": {
			user: user.User{
				Username:      "Existing Account Number",
				AccountNumber: existingUsers[0].AccountNumber,
				Email:         "testcase4@gmail.com",
			},
			expectedError: user.ErrSvcUserExists,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			user, err := uS.CreateUser(tc.user)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				cleanup = append(cleanup, user)
			}
		})
	}

	for _, u := range cleanup {
		if err := uR.DeleteUser(&u); err != nil {
			t.Fatal(err)
		}
	}
}

func getTestUsers(num int) []user.User {
	generatedUsers := []user.User{}

	for i := 1; i <= num; i++ {
		uName := fmt.Sprintf("Existing User %v", i)
		email := fmt.Sprintf("existinguser%v@gmail.com", i)
		u := user.User{
			ID:            common.GetMD5HashWithSum(uName + email),
			Username:      uName,
			Email:         email,
			AccountNumber: rangeIn(1000, 9999),
		}
		generatedUsers = append(generatedUsers, u)
	}

	return generatedUsers
}

func rangeIn(low, hi int) string {
	num := low + rand.Intn(hi-low)
	return strconv.Itoa(num)
}
