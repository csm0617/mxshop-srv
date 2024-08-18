package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

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
		Nickname: user.NiceName,
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
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(users)
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
