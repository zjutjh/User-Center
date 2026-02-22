package repo

import (
	"context"

	"github.com/zjutjh/User-Center/comm"
	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/dao/query"
	"github.com/zjutjh/mygo/ndb"
)

type StudentRepo struct {
	query *query.Query
}

func NewStudentRepo(tx ...*query.Query) *StudentRepo {
	var q *query.Query
	if len(tx) > 0 {
		q = tx[0]
	} else {
		q = query.Use(ndb.Pick())
	}
	return &StudentRepo{
		query: q,
	}
}

func (r *StudentRepo) GetStudentByStudentID(ctx context.Context, studentID string) (*model.Student, error) {
	q := r.query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil {
		return nil, comm.UserNotExist
	}
	return student, nil
}

func (r *StudentRepo) CheckStudentByStudentIDAndCardId(ctx context.Context, studentID string, cardId string) (*model.Student, error) {
	q := r.query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil || student == nil || student.IDCard != cardId {
		return nil, comm.UserNotExist
	}
	return student, nil
}
