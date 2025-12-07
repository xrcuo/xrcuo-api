package common

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatsHandler 处理统计信息展示页面
func StatsHandler(c *gin.Context) {
	// 获取统计信息
	stats := GlobalStats.GetStats()

	// 定义HTML模板
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API使用统计</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .stats-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        .stat-card {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
        }
        .stat-number {
            font-size: 2.5rem;
            font-weight: bold;
            color: #4CAF50;
        }
        .stat-label {
            font-size: 1rem;
            color: #666;
            margin-top: 5px;
        }
        .detail-section {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }
        .detail-section h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 1.5rem;
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .method-bar {
            display: inline-block;
            height: 20px;
            background-color: #4CAF50;
            border-radius: 10px;
            margin-right: 10px;
        }
        .method-label {
            display: inline-block;
            min-width: 60px;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <h1>API使用统计</h1>
    
    <div class="stats-container">
        <div class="stat-card">
            <div class="stat-number">{{.TotalCalls}}</div>
            <div class="stat-label">总调用次数</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">{{.DailyCalls}}</div>
            <div class="stat-label">今日调用次数</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">{{len .MethodCalls}}</div>
            <div class="stat-label">HTTP方法类型</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">{{len .PathCalls}}</div>
            <div class="stat-label">API路径数量</div>
        </div>
    </div>

    <div class="detail-section">
        <h2>按HTTP方法统计</h2>
        {{range $method, $count := .MethodCalls}}
            <div style="margin-bottom: 10px;">
                <span class="method-label">{{$method}}</span>
                <div class="method-bar" style="width: {{percentage $.TotalCalls $count}}%"></div>
                <span>{{$count}}次</span>
            </div>
        {{end}}
    </div>

    <div class="detail-section">
        <h2>按API路径统计</h2>
        <table>
            <thead>
                <tr>
                    <th>API路径</th>
                    <th>调用次数</th>
                </tr>
            </thead>
            <tbody>
                {{range $path, $count := .PathCalls}}
                    <tr>
                        <td>{{$path}}</td>
                        <td>{{$count}}</td>
                    </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <div class="detail-section">
        <h2>最近调用记录</h2>
        <table>
            <thead>
                <tr>
                    <th>时间</th>
                    <th>路径</th>
                    <th>方法</th>
                    <th>IP</th>
                    <th>状态码</th>
                </tr>
            </thead>
            <tbody>
                {{range .LastCallDetails}}
                    <tr>
                        <td>{{.Timestamp.Format "2006-01-02 15:04:05"}}</td>
                        <td>{{.Path}}</td>
                        <td>{{.Method}}</td>
                        <td>{{.IP}}</td>
                        <td>{{.StatusCode}}</td>
                    </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <script>
        // 定期刷新页面
        setInterval(function() {
            window.location.reload();
        }, 5000); // 每5秒刷新一次
    </script>
</body>
</html>
`

	// 创建模板并添加自定义函数
	t, err := template.New("stats").Funcs(template.FuncMap{
		"percentage": func(total, count int64) int {
			if total == 0 {
				return 0
			}
			return int((float64(count) / float64(total)) * 100)
		},
	}).Parse(tmpl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模板解析错误"})
		return
	}

	// 渲染模板
	if err := t.Execute(c.Writer, stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模板渲染错误"})
	}
}
