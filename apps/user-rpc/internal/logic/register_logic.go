package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	studentID := normalizeStudentID(in.StudentId)
	password := strings.TrimSpace(in.Password)
	cardID := normalizeCardID(in.CardId)
	email := strings.TrimSpace(in.Email)

	if studentID == "" || password == "" || cardID == "" {
		return nil, errorsx.ErrParameterInvalid
	}
	if len(password) < 6 || len(password) > 20 {
		return nil, errorsx.ErrPasswordLength
	}

	existing, err := l.svcCtx.Query.User.WithContext(l.ctx).
		Where(l.svcCtx.Query.User.StudentID.Eq(studentID)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Errorf("检查用户是否已存在失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if existing != nil {
		if existing.Password != hashPassword(password) {
			return nil, errorsx.ErrWrongAccountOrPassword
		}
		return &pb.RegisterResponse{}, nil
	}

	student, err := l.svcCtx.Query.Student.WithContext(l.ctx).
		Where(l.svcCtx.Query.Student.StudentID.Eq(studentID)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorsx.ErrUserNotExist
	}
	if err != nil {
		l.Errorf("查询学生信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if !strings.EqualFold(student.IDCard, cardID) {
		return nil, errorsx.ErrUserNotExist
	}

	if err := l.svcCtx.Query.User.WithContext(l.ctx).Create(&model.User{
		StudentID: studentID,
		Password:  hashPassword(password),
		Email:     email,
	}); err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.RegisterResponse{}, nil
}
