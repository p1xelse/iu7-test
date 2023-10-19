package delivery

import (
	"net/http"
	"strconv"
	projectUsecase "timetracker/internal/Project/usecase"
	"timetracker/internal/middleware"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Delivery struct {
	ProjectUC projectUsecase.UsecaseI
}

func (del *Delivery) ownerOrAdminValidate(c echo.Context, project *models.Project) error {
	user, ok := c.Get("user").(*models.User)

	if !ok {
		c.Logger().Error("can't get user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	if user.Role == models.Admin.String() || *project.UserID == user.ID {
		return nil
	}

	return models.ErrPermissionDenied
}

// CreateProject godoc
// @Summary      Create project
// @Description  Create project
// @Tags     	 project
// @Accept	 application/json
// @Produce  application/json
// @Param    project body dto.ReqCreateUpdateProject true "project info"
// @Success  200 {object} pkg.Response{body=dto.RespProject} "success update project"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 400 {object} echo.HTTPError "bad req"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /project/create [post]
func (delivery *Delivery) CreateProject(c echo.Context) error {

	var reqProject dto.ReqCreateUpdateProject
	err := c.Bind(&reqProject)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqProject); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	project := reqProject.ToModelProject()
	project.UserID = &userId
	err = delivery.ProjectUC.CreateProject(project)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respProject := dto.GetResponseFromModelProject(project)

	return c.JSON(http.StatusCreated, pkg.Response{Body: *respProject})
}

// GetProject godoc
// @Summary      Show a post
// @Description  Get project by id. Acl: admin, owner
// @Tags     	 project
// @Accept	 application/json
// @Produce  application/json
// @Param id  path int  true  "Project ID"
// @Success  200 {object} pkg.Response{body=dto.RespProject} "success get project"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /project/{id} [get]
func (delivery *Delivery) GetProject(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	project, err := delivery.ProjectUC.GetProject(id)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	err = delivery.ownerOrAdminValidate(c, project)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respProject := dto.GetResponseFromModelProject(project)
	return c.JSON(http.StatusOK, pkg.Response{Body: *respProject})
}

// UpdateProject godoc
// @Summary      Update an project
// @Description  Update an project. Acl: owner
// @Tags     	 project
// @Accept	 application/json
// @Produce  application/json
// @Param    project body dto.ReqCreateUpdateProject true "project info"
// @Success  200 {object} pkg.Response{body=dto.RespProject} "success update project"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /project/edit [post]
func (delivery *Delivery) UpdateProject(c echo.Context) error {

	var reqProject dto.ReqCreateUpdateProject
	err := c.Bind(&reqProject)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqProject); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	project := reqProject.ToModelProject()
	project.UserID = &userId
	err = delivery.ProjectUC.UpdateProject(project)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respProject := dto.GetResponseFromModelProject(project)

	return c.JSON(http.StatusOK, pkg.Response{Body: *respProject})
}

// DeleteProject godoc
// @Summary      Delete an project
// @Description  Delete an project. Acl: owner
// @Tags     	 project
// @Accept	 application/json
// @Param id path int  true  "Project ID"
// @Success  204
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 404 {object} echo.HTTPError "can't find project with such id"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /project/{id} [delete]
func (delivery *Delivery) DeleteProject(c echo.Context) error {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	err = delivery.ProjectUC.DeleteProject(id, userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMyProjects godoc
// @Summary      Get my projects
// @Description  Get my projects.
// @Tags     project
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespProject} "success get projects"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /me/projects [get]
func (delivery *Delivery) GetMyProjects(c echo.Context) error {

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	projects, err := delivery.ProjectUC.GetUserProjects(userId)
	// projects, err := delivery.ProjectUC.GetUserProjectsWithCache(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respProjects := dto.GetResponseFromModelProjects(projects)

	return c.JSON(http.StatusOK, pkg.Response{Body: respProjects})
}

// GetUserProjects godoc
// @Summary      Get user projects
// @Description  Get user projects. Acl: admin, friends
// @Tags     project
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespProject} "success get projects"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /user/{user_id}/projects [get]
func (delivery *Delivery) GetUserProjects(c echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	projects, err := delivery.ProjectUC.GetUserProjects(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respProjects := dto.GetResponseFromModelProjects(projects)

	return c.JSON(http.StatusOK, pkg.Response{Body: respProjects})
}

func handleError(err error) *echo.HTTPError {
	causeErr := errors.Cause(err)
	switch {
	case errors.Is(causeErr, models.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
	case errors.Is(causeErr, models.ErrBadRequest):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	case errors.Is(causeErr, models.ErrPermissionDenied):
		return echo.NewHTTPError(http.StatusForbidden, models.ErrPermissionDenied.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, causeErr.Error())
	}
}

func NewDelivery(e *echo.Echo, pu projectUsecase.UsecaseI, aclM *middleware.AclMiddleware) {
	handler := &Delivery{
		ProjectUC: pu,
	}

	e.POST("/project/create", handler.CreateProject)
	e.POST("/project/edit", handler.UpdateProject)  //acl: owner
	e.GET("/project/:id", handler.GetProject)       //acl: owner, admin
	e.DELETE("/project/:id", handler.DeleteProject) //acl: owner
	e.GET("/me/projects", handler.GetMyProjects)
	e.GET("/user/:user_id/projects", handler.GetUserProjects, aclM.FriendsOrAdminOnly)
}
