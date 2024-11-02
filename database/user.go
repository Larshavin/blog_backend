package database

import (
	ent "blog/ent"
	"blog/ent/token"
	user "blog/ent/user"
	"context"
)

func SaveUser(u *ent.User) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	_, err := client.User.Create().
		SetEmail(u.Email).
		SetHashedPassword(u.HashedPassword).
		SetAddress(u.Address).
		SetName(u.Name).
		SetSex(u.Sex).
		SetPhone(u.Phone).
		SetAge(u.Age).
		SetRole(u.Role).
		SetAvatar(u.Avatar).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*ent.User, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	user, err := client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func SaveToken(email, accessToken, refreshToken string) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	user, err := client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return err
	}

	// 만료 시간 얻기
	// accessTokenExpiresAt := time.Now().Add(time.Hour * 24)

	_, err = client.Token.Create().
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		SetUser(user).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTokenByRefreshToken(refreshToken string) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	_, err := client.Token.Delete().
		Where(token.RefreshTokenEQ(refreshToken)).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTokenByAccessToken(accessToken string) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	_, err := client.Token.Delete().
		Where(token.AccessTokenEQ(accessToken)).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
