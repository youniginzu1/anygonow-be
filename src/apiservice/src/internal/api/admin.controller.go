package api

import (
	"strconv"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var adminSet = wire.NewSet(wire.Struct(new(AdminController), "*"), wire.Struct(new(AdminService), "*"))

type AdminController struct {
	S AdminService
}

func (s AdminController) HandleUserBanPost(g *gin.Context) {
	req := pb.AdminBanUserPostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.BanUser(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleBusinessBanPost(g *gin.Context) {
	req := pb.AdminBusinessBanPostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.BanBusiness(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleBusinessDeletePost(g *gin.Context) {
	req := pb.AdminBusinessDeletePostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.DeleteBusiness(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleBusinessesGet(g *gin.Context) {
	req := pb.AdminBusinessesGetRequest{
		XUserId: g.GetString("userId"),
		Mail:    g.Query("mail"),
		Phone:   g.Query("phone"),
		Limit:   g.DefaultQuery("limit", "5"),
		Offset:  g.DefaultQuery("offset", "0"),
	}
	res, err := s.S.ListBusinesses(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleUsersGet(g *gin.Context) {
	req := pb.AdminUsersGetRequest{
		XUserId: g.GetString("userId"),
		Mail:    g.DefaultQuery("mail", ""),
		Phone:   g.DefaultQuery("phone", ""),
		Limit:   g.DefaultQuery("limit", "5"),
		Offset:  g.DefaultQuery("offset", "0"),
	}
	res, err := s.S.ListUsers(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleUserUnbanPost(g *gin.Context) {
	req := pb.AdminUsersUnbanPostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.UnbanUser(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleBusinessUnbanPost(g *gin.Context) {
	req := pb.AdminBusinessesUnbanPostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.UnbanBusiness(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleUserDeletePost(g *gin.Context) {
	req := pb.AdminUsersDeletePostRequest{
		Id:      g.Param("id"),
		XUserId: g.GetString("userId"),
	}
	res, err := s.S.DeleteUser(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleCategoriesGet(g *gin.Context) {
	q, err := strconv.Atoi(g.DefaultQuery("query", "0"))
	if err != nil {
		q = 0
	}
	req := pb.CategoriesGetRequest{
		Query:  c.QUERY_CATEGORY_ADMIN(q),
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
		Name:   g.DefaultQuery("name", ""),
	}
	res, err := s.S.ListCategories(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleCategoriesPost(g *gin.Context) {
	req := pb.AdminCategoryPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.PostCategory(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleCategoriesPostDelete(g *gin.Context) {

	req := pb.AdminCategoryPostDeleteRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.DeleteCategory(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleCategoriesPostEdit(g *gin.Context) {

	req := pb.AdminCategoryPostEditRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.EditCategory(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleGroupsGet(g *gin.Context) {
	req := pb.AdminGroupGetRequest{
		Limit:      g.DefaultQuery("limit", "5"),
		Offset:     g.DefaultQuery("offset", "0"),
		CategoryId: g.DefaultQuery("categoryId", ""),
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.GetGroups(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleGroupsPost(g *gin.Context) {
	req := pb.AdminGroupPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.AddGroup(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}
func (s AdminController) HandleGroupsPut(g *gin.Context) {

	req := pb.AdminGroupPutRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.Id = g.Param("id")

	res, err := s.S.EditGroup(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleAdvertiseManagementPost(g *gin.Context) {
	req := pb.AdminAdvertiseManagementPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.CreateAdvertisePackage(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AdminController) HandleAdvertiseManagementGet(g *gin.Context) {
	req := pb.AdminAdvertiseManagementGetRequest{
		Limit:       g.DefaultQuery("limit", "5"),
		Offset:      g.DefaultQuery("offset", "0"),
		ServiceName: g.DefaultQuery("serviceName", ""),
	}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetAdvertisePackages(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleAdvertiseManagementPut(g *gin.Context) {

	req := pb.AdminAdvertiseManagementPutRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.Id = g.Param("id")

	res, err := s.S.EditAdvertisePackage(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s AdminController) HandleAdvertiseManagementDeletePost(g *gin.Context) {
	req := pb.AdminAdvertiseManagementDeletePostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.DeleteAdvertisePackage(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
