package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-lib/db"
)

// ApiDoc 表示API接口文档
type ApiDoc struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        string    `json:"tags"`
	Parameters  string    `json:"parameters"`
	Headers     string    `json:"headers"`
	RequestBody string    `json:"request_body"`
	Responses   string    `json:"responses"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ApiDocListItem 列表项（不包含详细内容）
type ApiDocListItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateApiDocsTable 创建API文档表
func CreateApiDocsTable() error {
	sqlStr := `
	CREATE TABLE IF NOT EXISTS api_docs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		method TEXT NOT NULL DEFAULT 'GET',
		description TEXT,
		category TEXT NOT NULL DEFAULT '默认',
		tags TEXT DEFAULT '[]',
		parameters TEXT DEFAULT '[]',
		headers TEXT DEFAULT '[]',
		request_body TEXT DEFAULT '{}',
		responses TEXT DEFAULT '[]',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_api_docs_category ON api_docs(category);
	CREATE INDEX IF NOT EXISTS idx_api_docs_method ON api_docs(method);
	`
	_, err := db.GetDB().Exec(sqlStr)
	return err
}

// ApiDocCreate 创建API文档
func ApiDocCreate(doc *ApiDoc) error {
	result, err := db.GetDB().Exec(
		"INSERT INTO api_docs (name, path, method, description, category, tags, parameters, headers, request_body, responses) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		doc.Name, doc.Path, doc.Method, doc.Description, doc.Category, doc.Tags, doc.Parameters, doc.Headers, doc.RequestBody, doc.Responses,
	)
	if err != nil {
		return err
	}
	doc.ID, _ = result.LastInsertId()
	return nil
}

// ApiDocGetByID 根据ID获取API文档
func ApiDocGetByID(id int64) (*ApiDoc, error) {
	row := db.GetDB().QueryRow(
		"SELECT id, name, path, method, description, category, tags, parameters, headers, request_body, responses, created_at, updated_at FROM api_docs WHERE id = ?",
		id,
	)
	doc := &ApiDoc{}
	var createdAt, updatedAt sql.NullTime
	err := row.Scan(&doc.ID, &doc.Name, &doc.Path, &doc.Method, &doc.Description, &doc.Category, &doc.Tags, &doc.Parameters, &doc.Headers, &doc.RequestBody, &doc.Responses, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		doc.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		doc.UpdatedAt = updatedAt.Time
	}
	return doc, nil
}

// ApiDocList 获取API文档列表
func ApiDocList(category, method, keyword string) ([]ApiDocListItem, error) {
	query := "SELECT id, name, path, method, description, category, tags, created_at, updated_at FROM api_docs WHERE 1=1"
	args := []interface{}{}

	if category != "" && category != "全部" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if method != "" && method != "全部" {
		query += " AND method = ?"
		args = append(args, method)
	}
	if keyword != "" {
		query += " AND (name LIKE ? OR path LIKE ? OR description LIKE ?)"
		like := "%" + keyword + "%"
		args = append(args, like, like, like)
	}
	query += " ORDER BY category, id"

	rows, err := db.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []ApiDocListItem
	for rows.Next() {
		var doc ApiDocListItem
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&doc.ID, &doc.Name, &doc.Path, &doc.Method, &doc.Description, &doc.Category, &doc.Tags, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		if createdAt.Valid {
			doc.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			doc.UpdatedAt = updatedAt.Time
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// ApiDocUpdate 更新API文档
func ApiDocUpdate(doc *ApiDoc) error {
	_, err := db.GetDB().Exec(
		"UPDATE api_docs SET name = ?, path = ?, method = ?, description = ?, category = ?, tags = ?, parameters = ?, headers = ?, request_body = ?, responses = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		doc.Name, doc.Path, doc.Method, doc.Description, doc.Category, doc.Tags, doc.Parameters, doc.Headers, doc.RequestBody, doc.Responses, doc.ID,
	)
	return err
}

// ApiDocDelete 删除API文档
func ApiDocDelete(id int64) error {
	_, err := db.GetDB().Exec("DELETE FROM api_docs WHERE id = ?", id)
	return err
}

// ApiDocGetCategories 获取所有分类
func ApiDocGetCategories() ([]string, error) {
	rows, err := db.GetDB().Query("SELECT DISTINCT category FROM api_docs ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err == nil && cat != "" {
			categories = append(categories, cat)
		}
	}
	return categories, nil
}

// InitApiDocsFromData 从数据初始化API文档（首次启动时）
func InitApiDocsFromData(docs []map[string]interface{}) error {
	for _, data := range docs {
		tags, _ := json.Marshal(data["tags"])
		parameters, _ := json.Marshal(data["parameters"])
		headers, _ := json.Marshal(data["headers"])
		requestBody, _ := json.Marshal(data["requestBody"])
		responses, _ := json.Marshal(data["responses"])

		doc := &ApiDoc{
			Name:        getString(data, "name"),
			Path:        getString(data, "path"),
			Method:      getString(data, "method"),
			Description: getString(data, "description"),
			Category:    getString(data, "category"),
			Tags:        string(tags),
			Parameters:  string(parameters),
			Headers:     string(headers),
			RequestBody: string(requestBody),
			Responses:   string(responses),
		}

		_, err := db.GetDB().Exec(
			"INSERT OR IGNORE INTO api_docs (name, path, method, description, category, tags, parameters, headers, request_body, responses) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			doc.Name, doc.Path, doc.Method, doc.Description, doc.Category, doc.Tags, doc.Parameters, doc.Headers, doc.RequestBody, doc.Responses,
		)
		if err != nil {
			return fmt.Errorf("初始化文档 %s 失败: %v", doc.Name, err)
		}
	}
	return nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
