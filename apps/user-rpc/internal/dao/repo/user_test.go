package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestUserRepo(t *testing.T) (*UserRepo, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`
		CREATE TABLE user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id TEXT NOT NULL,
			password TEXT NOT NULL,
			phone_num TEXT NOT NULL DEFAULT '',
			user_type TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			device_id TEXT NOT NULL DEFAULT '',
			yxy_uid TEXT NOT NULL DEFAULT '',
			zf_password TEXT NOT NULL DEFAULT '',
			oauth_password TEXT NOT NULL DEFAULT '',
			create_time DATETIME,
			update_time DATETIME
		)
	`).Error)

	return NewUserRepo(query.Use(db)), db
}

func TestUserRepoGetUserByIdMapsNotFound(t *testing.T) {
	repo, _ := newTestUserRepo(t)

	_, err := repo.GetUserById(context.Background(), 123)

	require.ErrorIs(t, err, errorsx.ErrUserNotExist)
}

func TestUserRepoExistsByStudentID(t *testing.T) {
	repo, db := newTestUserRepo(t)
	ctx := context.Background()
	require.NoError(t, db.WithContext(ctx).Create(&model.User{StudentID: "ABC123", Password: "secret"}).Error)

	exists, err := repo.ExistsByStudentID(ctx, "ABC123")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = repo.ExistsByStudentID(ctx, "MISSING")
	require.NoError(t, err)
	require.False(t, exists)
}

func TestUserRepoUpdatePasswordByID(t *testing.T) {
	repo, db := newTestUserRepo(t)
	ctx := context.Background()
	user := &model.User{StudentID: "ABC123", Password: "old"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	require.NoError(t, repo.UpdatePasswordByID(ctx, user.ID, "new"))

	var stored model.User
	require.NoError(t, db.WithContext(ctx).First(&stored, user.ID).Error)
	require.Equal(t, "new", stored.Password)
}

func TestUserRepoUpdateYxyBindByID(t *testing.T) {
	repo, db := newTestUserRepo(t)
	ctx := context.Background()
	user := &model.User{StudentID: "ABC123", Password: "secret"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	require.NoError(t, repo.UpdateYxyBindByID(ctx, user.ID, "device-1", "uid-1"))

	var stored model.User
	require.NoError(t, db.WithContext(ctx).First(&stored, user.ID).Error)
	require.Equal(t, "device-1", stored.DeviceID)
	require.Equal(t, "uid-1", stored.YxyUID)
}

func TestUserRepoUpdateZfBindByID(t *testing.T) {
	repo, db := newTestUserRepo(t)
	ctx := context.Background()
	user := &model.User{StudentID: "ABC123", Password: "secret"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	require.NoError(t, repo.UpdateZfBindByID(ctx, user.ID, "zf-ciphertext"))

	var stored model.User
	require.NoError(t, db.WithContext(ctx).First(&stored, user.ID).Error)
	require.Equal(t, "zf-ciphertext", stored.ZfPassword)
}

func TestUserRepoUpdateOauthBindByID(t *testing.T) {
	repo, db := newTestUserRepo(t)
	ctx := context.Background()
	user := &model.User{StudentID: "ABC123", Password: "secret"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	require.NoError(t, repo.UpdateOauthBindByID(ctx, user.ID, "oauth-ciphertext"))

	var stored model.User
	require.NoError(t, db.WithContext(ctx).First(&stored, user.ID).Error)
	require.Equal(t, "oauth-ciphertext", stored.OauthPassword)
}
