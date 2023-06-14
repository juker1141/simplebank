package token

import (
	"testing"
	"time"

	"github.com/juker1141/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	_, err := NewPasetoMaker(util.RandomString(20))
	require.Error(t, err)

	maker1, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token1, err := maker1.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token1)

	token2, err := maker2.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token2)

	payload, err := maker1.VerifyToken(token2)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}