package controllers

import (
	"goblog/database"
	"goblog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateCommentInput struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// CreateComment godoc
// @Summary 创建评论或回复
// @Tags 评论
// @Accept json
// @Produce json
// @Param comment body CreateCommentInput true "评论数据"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /comments [post]
// @Security ApiKeyAuth
func CreateComment(c *gin.Context) {
	var input CreateCommentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效", "detail": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	comment := models.Conment{
		Content:  input.Content,
		UserID:   userID,
		PostID:   input.PostID,
		ParentID: input.ParentID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	if err := database.DB.Preload("User").Preload("Replies.User").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载用户消息失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "评论成功", "comment": comment})
}

// GetCommentsByPostID godoc
// @Summary 获取文章下的评论列表（带子评论）
// @Tags 评论
// @Accept json
// @Produce json
// @Param post_id query int true "文章 ID"
// @Success 200 {object} map[string]interface{}
// @Router /comments [get]
func GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Query("post_id")
	if postIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "必须提供 post_id 参数"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id 无效"})
		return
	}

	var comments []models.Conment

	if err := database.DB.Preload("User").Preload("Replies").Preload("Replies.User").Where("post_id = ? AND parent_id IS NULL", postID).Order("created_at ASC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询评论失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
