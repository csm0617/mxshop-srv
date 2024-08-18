package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func ModelToResponse(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Nickname: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

type UserService struct {
}

// 把方法绑定到接口上
func (s *UserService) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users) //必须传&地址进去，否则报错，panic: reflect: reflect.Value.SetLen using unaddressable value
	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

func (s *UserService) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//没有这个用户
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	//查询报错
	if result.Error != nil {
		return nil, result.Error
	}
	//找到用户
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	//没有这个用户
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	//查询报错
	if result.Error != nil {
		return nil, result.Error
	}
	//找到用户
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreatUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	var user model.User
	user.Mobile = req.Mobile
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//用户已存在
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//用户不存在则新建
	user.NickName = req.Nickname
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodePwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePwd)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRsp := ModelToResponse(user)
	userInfoRsp.Password = ""
	return &userInfoRsp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	//此时已经找到用户了，user也被赋值了，只需要修改能修改的
	user.NickName = req.NickName
	user.Gender = req.Gender
	//把uint64转int64后再转成time
	birthday := time.Unix(int64(req.Birthday), 0)
	user.Birthday = &birthday
	//更新用户
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserService) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	//加密规则
	options := &password.Options{16, 100, 32, sha512.New}
	//req.EncryptedPassword是加密后的密码
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	//req.Password 原始密码
	verify := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)

	return &proto.CheckResponse{Success: verify}, nil
}
