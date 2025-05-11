package controllers

import (
	"fmt"
	"goblog/database"
	"goblog/models"
	"goblog/pkg/cache"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePostInput struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	IsDraft     bool     `json:"is_draft"`
	IsTop       bool     `json:"is_top"`
	IsRecommend bool     `json:"is_recommend"`
	Tags        []string `json:"tags"`
}

type UpdatePostInput struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	IsDraft     *bool   `json:"is_draft"`
	IsTop       *bool   `json:"is_top"`
	IsRecommend *bool   `json:"is_recommend"`
}

// CreatePost godoc
// @Summary 创建文章
// @Description 登录用户可发布新文章，并附带标签
// @Tags 文章
// @Accept json
// @Produce json
// @Param post body CreatePostInput true "文章数据"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /posts [post]
// @Security ApiKeyAuth
func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	var tags []*models.Tag

	for _, tagName := range input.Tags {
		var tag models.Tag
		if err := database.DB.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err == nil {
			tags = append(tags, &tag)
		}
	}

	post := models.Post{
		Title:       input.Title,
		Content:     input.Content,
		UserID:      userID,
		IsDraft:     input.IsDraft,
		IsTop:       input.IsTop,
		IsRecommend: input.IsRecommend,
		Tags:        tags,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "文章创建成功", "post": post})

	for i := 1; i <= 5; i++ {
		cacheKey := fmt.Sprintf("posts:page:%d:limit:10", i)
		cache.Rdb.Del(cache.Ctx, cacheKey)
	}
}

// GetPosts godoc
// @Summary 获取文章列表
// @Description 支持分页获取文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /posts [get]
func GetPosts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	cacheKey := fmt.Sprintf("posts:page:%s:limit:%s", page, limit)

	var posts []models.Post

	if hit, err := cache.GetJSON(cacheKey, &posts); hit && err == nil {
		c.JSON(http.StatusOK, gin.H{"from": "cache", "posts": posts})
		return
	}

	offset, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(limit)
	offset = (offset - 1) * pageSize
	// page, _ := strconv.Atoi(pageStr)
	// limit, _ := strconv.Atoi(limitStr)
	// offset := (page - 1) * limit

	if err := database.DB.Preload("User").Preload("Tags").Order("created_at desc").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		fmt.Println("查询出错:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
		return
	}
	cache.SetJSON(cacheKey, posts, 30*time.Second)
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetPostByID godoc
// @Summary 获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [get]
func GetPostByID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}

	var post models.Post

	if err := database.DB.Preload("User").Preload("Tags").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// UpdatePost godoc
// @Summary 修改文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Param post body UpdatePostInput true "更新数据"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /posts/{id} [put]
// @Security ApiKeyAuth
func UpdataPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}

	var post models.Post

	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改该文章"})
		return
	}

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体无效"})
		return
	}

	updatdData := map[string]interface{}{}
	if input.Title != nil {
		updatdData["title"] = *input.Title
	}
	if input.Content != nil {
		updatdData["content"] = *input.Content
	}
	if input.IsDraft != nil {
		updatdData["is_draft"] = *input.IsDraft
	}
	if input.IsTop != nil {
		updatdData["is_top"] = *input.IsTop
	}
	if input.IsRecommend != nil {
		updatdData["is_recommend"] = *input.IsRecommend
	}

	if err := database.DB.Model(&post).Updates(updatdData).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功", "post": post})
}

// DeletePost godoc
// @Summary 删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /posts/{id} [delete]
// @Security ApiKeyAuth
func DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除该文章"})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
