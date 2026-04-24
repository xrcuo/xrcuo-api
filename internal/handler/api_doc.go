package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/common"
)

// ApiDocCreateRequest 创建API文档请求
type ApiDocCreateRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Path        string                   `json:"path" binding:"required"`
	Method      string                   `json:"method" binding:"required"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Tags        []string                 `json:"tags"`
	Parameters  []map[string]interface{} `json:"parameters"`
	Headers     []map[string]interface{} `json:"headers"`
	RequestBody map[string]interface{}   `json:"requestBody"`
	Responses   []map[string]interface{} `json:"responses"`
}

// ApiDocUpdateRequest 更新API文档请求
type ApiDocUpdateRequest struct {
	ID          int64                    `json:"id" binding:"required"`
	Name        string                   `json:"name" binding:"required"`
	Path        string                   `json:"path" binding:"required"`
	Method      string                   `json:"method" binding:"required"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Tags        []string                 `json:"tags"`
	Parameters  []map[string]interface{} `json:"parameters"`
	Headers     []map[string]interface{} `json:"headers"`
	RequestBody map[string]interface{}   `json:"requestBody"`
	Responses   []map[string]interface{} `json:"responses"`
}

// ApiDocList 获取API文档列表
func ApiDocList(c *gin.Context) {
	category := c.Query("category")
	method := c.Query("method")
	keyword := c.Query("keyword")

	docs, err := models.ApiDocList(category, method, keyword)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取列表失败")
		return
	}

	common.SuccessResponse(c, docs, "获取成功")
}

// ApiDocGet 获取单个API文档
func ApiDocGet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "无效的ID")
		return
	}

	doc, err := models.ApiDocGetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, common.CodeNotFound, "文档不存在")
		return
	}

	common.SuccessResponse(c, doc, "获取成功")
}

// ApiDocCreate 创建API文档
func ApiDocCreate(c *gin.Context) {
	var req ApiDocCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	tags, _ := json.Marshal(req.Tags)
	parameters, _ := json.Marshal(req.Parameters)
	headers, _ := json.Marshal(req.Headers)
	requestBody, _ := json.Marshal(req.RequestBody)
	responses, _ := json.Marshal(req.Responses)

	doc := &models.ApiDoc{
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
		Category:    req.Category,
		Tags:        string(tags),
		Parameters:  string(parameters),
		Headers:     string(headers),
		RequestBody: string(requestBody),
		Responses:   string(responses),
	}

	if doc.Category == "" {
		doc.Category = "默认"
	}

	if err := models.ApiDocCreate(doc); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "创建失败: "+err.Error())
		return
	}

	common.SuccessResponse(c, doc, "创建成功")
}

// ApiDocUpdate 更新API文档
func ApiDocUpdate(c *gin.Context) {
	var req ApiDocUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	tags, _ := json.Marshal(req.Tags)
	parameters, _ := json.Marshal(req.Parameters)
	headers, _ := json.Marshal(req.Headers)
	requestBody, _ := json.Marshal(req.RequestBody)
	responses, _ := json.Marshal(req.Responses)

	doc := &models.ApiDoc{
		ID:          req.ID,
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
		Category:    req.Category,
		Tags:        string(tags),
		Parameters:  string(parameters),
		Headers:     string(headers),
		RequestBody: string(requestBody),
		Responses:   string(responses),
	}

	if doc.Category == "" {
		doc.Category = "默认"
	}

	if err := models.ApiDocUpdate(doc); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "更新失败: "+err.Error())
		return
	}

	common.SuccessResponse(c, doc, "更新成功")
}

// ApiDocDelete 删除API文档
func ApiDocDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "无效的ID")
		return
	}

	if err := models.ApiDocDelete(id); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "删除失败")
		return
	}

	common.SuccessResponse(c, nil, "删除成功")
}

// ApiDocCategories 获取所有分类
func ApiDocCategories(c *gin.Context) {
	categories, err := models.ApiDocGetCategories()
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取分类失败")
		return
	}

	common.SuccessResponse(c, categories, "获取成功")
}
