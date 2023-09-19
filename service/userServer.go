package service

import (
	"context"
	"time"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/models"
	"github.com/lightsaid/blogs/routers/forms"
	"github.com/lightsaid/blogs/token"
)

type UserServer struct {
	store   dbrepo.UserRepo
	session dbrepo.SessionRepo
}

func NewUserServer(store dbrepo.UserRepo, session dbrepo.SessionRepo) *UserServer {
	return &UserServer{
		store:   store,
		session: session,
	}
}

func (srv *UserServer) Create(ctx context.Context, req forms.AddUserRequest) (int64, *errs.AppError) {
	user := models.User{
		Email:    req.Email,
		UserName: req.UserName,
		Avatar:   req.Avatar,
	}

	err := user.SetHashedPassword(req.Password)
	if err != nil {
		return 0, errs.ErrBadRequest.AsException(err)
	}

	newID, err := srv.store.Insert(ctx, &user)
	if err != nil {
		return 0, errs.HandleSQLError(err)
	}

	return newID, nil
}

func (srv *UserServer) Get(ctx context.Context, id int64) (*models.User, *errs.AppError) {
	user, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return user, nil
}

func (srv *UserServer) GetByEmail(ctx context.Context, email string) (*models.User, *errs.AppError) {
	user, err := srv.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return user, nil
}

func (srv *UserServer) ActivateUser(ctx context.Context, id int64) *errs.AppError {
	err := srv.store.Activate(ctx, id)
	if err != nil {
		return errs.HandleSQLError(err)
	}
	return nil
}

func (srv *UserServer) Update(ctx context.Context, req forms.UpdateUserRequest) *errs.AppError {
	user, err := srv.store.Get(ctx, req.ID)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.UserName != "" {
		user.UserName = req.UserName
	}

	err = srv.store.Update(ctx, user)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	return nil
}

func (srv *UserServer) Login(ctx context.Context, req forms.LoginRequest) (*forms.LoginResponse, *errs.AppError) {
	user, err := srv.store.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	// 检查密码
	ok := user.MatchesPassword(req.Password, user.Password)
	if !ok {
		return nil, errs.ErrBadRequest.AsException(err, "密码不匹配")
	}

	// 检查是否已激活
	if user.ActivatedAt == nil || *user.ActivatedAt == "" {
		return nil, errs.ErrBadRequest.AsException(err, "请先去邮箱激活账号，再登录")
	}

	// 创建Token
	accessPayload, aErr := token.NewPayload(user.ID, config.ParseDuration(config.AppConf.Token.TokenExpire, 15*time.Minute))

	// 刷新 Token refresh Token
	refreshPayload, rErr := token.NewPayload(user.ID, config.ParseDuration(config.AppConf.Token.RefreshExpire, 72*time.Hour))

	if aErr != nil || rErr != nil {
		return nil, errs.ErrInternalServer.AsException(aErr).AsException(rErr)
	}

	// 访问 token access Token
	aToken, aErr := config.TokenMaker.GenToken(accessPayload)

	rToken, rErr := config.TokenMaker.GenToken(refreshPayload)

	if aErr != nil || rErr != nil {
		return nil, errs.ErrInternalServer.AsException(aErr).AsException(rErr)
	}

	// 创建 session
	sess := models.Session{
		UserID:       user.ID,
		RefreshToken: rToken,
		ClientIP:     req.ClientIP,
		ExpiredAt:    refreshPayload.ExpiredAt.Format(srvTimeLayout),
	}
	_, err = srv.session.Insert(ctx, &sess)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	// 返回 response
	res := forms.LoginResponse{
		ID:           user.ID,
		UserName:     user.UserName,
		Avatar:       user.Avatar,
		RefreshToken: rToken,
		AccessToken:  aToken,
		LoginAt:      time.Now().Format(srvTimeLayout),
	}

	return &res, nil
}
