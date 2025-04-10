package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/internal/models"
	"github.com/jimil-28/crowd-monitor/internal/services/firebase"
	"github.com/jimil-28/crowd-monitor/internal/utils"
)

type UserHandler struct {
	firebaseClient *firebase.Client
}

func NewUserHandler(firebaseClient *firebase.Client) *UserHandler {
	return &UserHandler{
		firebaseClient: firebaseClient,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.firebaseClient.GetAllUsers(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.firebaseClient.SaveUser(c, user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User added successfully", user)
}
