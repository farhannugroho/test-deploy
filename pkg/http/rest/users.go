package rest

import (
	"errors"
	"food_delivery_api/pkg/model"
	"food_delivery_api/pkg/svc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addUser(s svc.Service) gin.HandlerFunc {
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

func getUsers(s svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res []model.User
		var err error

		if res, err = s.GetUsers(); err != nil {
			c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{Success: true, Data: res})
	}
}

func getUser(s svc.Service) gin.HandlerFunc {
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

func editUser(s svc.Service) gin.HandlerFunc {
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

func removeUser(s svc.Service) gin.HandlerFunc {
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
