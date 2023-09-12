package dbrepo_test

import (
	"context"
	"testing"

	"github.com/lightsaid/blogs/models"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {
	user := models.User{
		Email:    "lightsaid@qq.com",
		Password: "abc123",
		UserName: "lightsaid",
		Avatar:   "http://",
	}
	err := user.SetHashedPassword(user.Password)
	require.NoError(t, err)

	newID, err := testRepo.UserRepo.Insert(context.Background(), &user)
	require.NoError(t, err)
	require.Greater(t, newID, int64(0))

	qUser, err := testRepo.UserRepo.Get(context.Background(), newID)
	require.NoError(t, err)
	ok := qUser.MatchesPassword("abc123", qUser.Password)
	require.True(t, ok)

	require.Equal(t, user.Email, qUser.Email)
	require.Equal(t, user.UserName, qUser.UserName)
	require.Equal(t, user.Avatar, qUser.Avatar)

	// ....

	// TODO: 写详情测试

}
