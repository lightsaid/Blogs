package service

import (
	"context"
	"net/http"
	"time"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/dbrepo"
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

func (srv *UserServer) Create(ctx context.Context, req forms.AddUserRequest) (int64, *dbrepo.DBError) {
	user := models.User{
		Email:    req.Email,
		UserName: req.UserName,
		Avatar:   req.Avatar,
	}

	err := user.SetHashedPassword(req.Password)
	if err != nil {
		return 0, dbrepo.NewDBError("密码格式不对", http.StatusBadRequest).AsError(err)
	}

	newID, err := srv.store.Insert(ctx, &user)
	if err != nil {
		return 0, dbrepo.CheckError(err)
	}

	return newID, nil
}

func (srv *UserServer) Get(ctx context.Context, id int64) (*models.User, *dbrepo.DBError) {
	user, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return user, nil
}

func (srv *UserServer) GetByEmail(ctx context.Context, email string) (*models.User, *dbrepo.DBError) {
	user, err := srv.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return user, nil
}

func (srv *UserServer) ActivateUser(ctx context.Context, id int64) *dbrepo.DBError {
	err := srv.store.Activate(ctx, id)
	if err != nil {
		return dbrepo.CheckError(err)
	}
	return nil
}

func (srv *UserServer) Update(ctx context.Context, req forms.UpdateUserRequest) *dbrepo.DBError {
	user, err := srv.store.Get(ctx, req.ID)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.UserName != "" {
		user.UserName = req.UserName
	}

	err = srv.store.Update(ctx, user)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	return nil
}

func (srv *UserServer) Login(ctx context.Context, req forms.LoginRequest) (*forms.LoginResponse, *dbrepo.DBError) {
	user, err := srv.store.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	// 检查密码
	ok := user.MatchesPassword(req.Password, user.Password)
	if !ok {
		return nil, dbrepo.NewDBError("密码不匹配", http.StatusBadRequest)
	}

	// 检查是否已激活
	if user.ActivatedAt == nil || *user.ActivatedAt == "" {
		return nil, dbrepo.NewDBError("请先去邮箱激活账号，再登录", http.StatusBadRequest)
	}

	// 创建Token
	accessPayload, aErr := token.NewPayload(user.ID, config.ParseDuration(config.AppConf.Token.TokenExpire, 15*time.Minute))

	// 刷新 Token refresh Token
	refreshPayload, rErr := token.NewPayload(user.ID, config.ParseDuration(config.AppConf.Token.RefreshExpire, 72*time.Hour))

	if aErr != nil || rErr != nil {
		return nil, dbrepo.ErrDBInternal.AsError(aErr).AsError(rErr)
	}

	// 访问 token access Token
	aToken, aErr := config.TokenMaker.GenToken(accessPayload)

	rToken, rErr := config.TokenMaker.GenToken(refreshPayload)

	if aErr != nil || rErr != nil {
		return nil, dbrepo.ErrDBInternal.AsError(aErr).AsError(rErr)
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
		return nil, dbrepo.CheckError(err)
	}

	// 返回 response
	res := forms.LoginResponse{
		ID:           user.ID,
		UserName:     user.UserName,
		Avatar:       user.Avatar,
		RefreshToken: rToken,
		AccessToken:  aToken,
	}

	return &res, nil
}
