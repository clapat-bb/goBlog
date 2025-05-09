package controllers

import (
	"goblog/database"
	"goblog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPostsByTag godoc
// @Summary 获取指定标签下的所有文章
// @Tags 标签
// @Accept json
// @Produce json
// @Param name path string true "标签名"
// @Success 200 {object} map[string]interface{}
// @Router /tags/{name}/posts [get]
func GetPostByTag(c *gin.Context) {
	tagName := c.Param("name")
	var tag models.Tag

	if err := database.DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Preload("Tags").Order("created_at desc")
	}).First(&tag, "name = ?", tagName).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":   tag.Name,
		"posts": tag.Posts,
		"count": len(tag.Posts),
	})
}
