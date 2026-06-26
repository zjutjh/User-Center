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

func newTestStudentRepo(t *testing.T) (*StudentRepo, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`
		CREATE TABLE student (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id TEXT NOT NULL,
			id_card TEXT,
			create_time DATETIME,
			update_time DATETIME
		)
	`).Error)

	return NewStudentRepo(query.Use(db)), db
}

func TestStudentRepoGetStudentByStudentIDMapsNotFound(t *testing.T) {
	repo, _ := newTestStudentRepo(t)

	_, err := repo.GetStudentByStudentID(context.Background(), "missing")

	require.ErrorIs(t, err, errorsx.ErrUserNotExist)
}

func TestStudentRepoGetStudentByStudentID(t *testing.T) {
	repo, db := newTestStudentRepo(t)
	ctx := context.Background()
	require.NoError(t, db.WithContext(ctx).Create(&model.Student{StudentID: "ABC123", IDCard: "ID123"}).Error)

	student, err := repo.GetStudentByStudentID(ctx, "ABC123")

	require.NoError(t, err)
	require.Equal(t, "ID123", student.IDCard)
}
