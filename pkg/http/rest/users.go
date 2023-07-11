package rest

import (
	"errors"
	"food_delivery_api/pkg/model"
	"food_delivery_api/pkg/service"
	"food_delivery_api/pkg/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddUser(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body model.User
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		res, err := s.AddUser(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func AddUsers(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		res, err := s.AddUsers(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func GetUsers(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res []model.User
		var ttl int64
		var err error

		qp := model.QueryPagination{}
		err = c.ShouldBindQuery(&qp)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		if res, ttl, err = s.GetUsers(qp); err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		paginate(c, res, ttl)
	}
}

func GetUser(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res model.User

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: errors.New("id must be uint").Error()})
			return
		}
		res.ID = uint(id)

		res, err = s.GetUser(res)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func GetLoggedUser(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res model.User

		token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", -1)

		id, err := util.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}
		res.ID = uint(id)

		res, err = s.GetUser(res)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func EditUser(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body model.User
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: errors.New("id must be uint").Error()})
			return
		}
		body.ID = uint(id)

		res, err := s.EditUser(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func RemoveUser(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body model.User

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: errors.New("id must be uint").Error()})
			return
		}
		body.ID = uint(id)

		res, err := s.RemoveUser(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}
