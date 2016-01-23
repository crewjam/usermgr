package usermgr

import . "gopkg.in/check.v1"

var _ = Suite(&TestUser{})

type TestUser struct {
}

func (s *TestUser) TestGroup(c *C) {
	var testCases = []struct {
		User   User
		Groups []string
		Result bool
	}{
		{User: User{Groups: []string{"a", "b", "c"}}, Groups: []string{"a"}, Result: true},
		{User: User{Groups: []string{"a", "b", "c"}}, Groups: []string{"c"}, Result: true},
		{User: User{Groups: []string{"a", "b", "c"}}, Groups: []string{"d"}, Result: false},
		{User: User{Groups: []string{}}, Groups: []string{"a"}, Result: false},
		{User: User{Groups: []string{"a"}}, Groups: []string{}, Result: false},
		{User: User{Groups: []string{"a"}}, Groups: nil, Result: false},
		{User: User{Groups: nil}, Groups: []string{"a"}, Result: false},
	}

	for _, testCase := range testCases {
		result := testCase.User.InAnyGroup(testCase.Groups)
		c.Assert(result, Equals, testCase.Result)
	}
}
