package controllers

import (
	"goblog/database"
	"goblog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeInput struct {
	TargetID   uint   `json:"target_id" binding:"required"`
	TargetType string `json:"target_type" binding:"required,oneof=post comment"`
}

// Like godoc
// @Summary 点赞文章或评论
// @Tags 点赞
// @Accept json
// @Produce json
// @Param like body LikeInput true "点赞对象"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /likes [post]
// @Security ApiKeyAuth
func Like(c *gin.Context) {
	var input LikeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	var like models.Like
	result := database.DB.Where("user_id = ? AND target_id ? AND target_type = ?", userID, input.TargetID, input.TargetType).First(&like)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "你已经点过赞了"})
		return
	}

	newLike := models.Like{
		UserID:     userID,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
	}

	if err := database.DB.Create(&newLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "点赞成功"})
}

// GetLikeCount godoc
// @Summary 获取点赞总数
// @Tags 点赞
// @Accept json
// @Produce json
// @Param target_id query int true "目标 ID"
// @Param target_type query string true "目标类型（post/comment）"
// @Success 200 {object} map[string]int
// @Router /likes/count [get]
func GetLikeCount(c *gin.Context) {
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	if targetID == "" || targetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 target_id 或 target_type"})
		return
	}

	var count int64
	if err := database.DB.Model(&models.Like{}).Where("target_id = ? AND target_type = ?", targetID, targetType).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"like_count": count})
}

// CheckIfLiked godoc
// @Summary 检查用户是否已点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param target_id query int true "目标 ID"
// @Param target_type query string true "目标类型（post/comment）"
// @Success 200 {object} map[string]bool
// @Router /likes/check [get]
// @Security ApiKeyAuth
func CheckIfLiked(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	if targetID == "" || targetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}

	var like models.Like
	result := database.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&like)
	c.JSON(http.StatusOK, gin.H{"liked": result.RowsAffected > 0})
}
