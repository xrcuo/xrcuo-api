// 全局变量保存图表实例
let pathStatsChart = null;

// 计算百分比函数
function calculatePercentage(part, total) {
    if (total === 0) return 0;
    return Math.round((part / total) * 100);
}

// 更新统计卡片数据
function updateStatCards(stats) {
    // 更新总调用次数
    document.querySelector('.stat-card:nth-child(1) .stat-number').textContent = stats.total_calls;
    
    // 更新今日调用次数
    document.querySelector('.stat-card:nth-child(2) .stat-number').textContent = stats.daily_calls;
    
    // 更新HTTP方法类型数量
    document.querySelector('.stat-card:nth-child(3) .stat-number').textContent = Object.keys(stats.method_calls).length;
    
    // 更新API路径数量
    document.querySelector('.stat-card:nth-child(4) .stat-number').textContent = Object.keys(stats.path_calls).length;
}

// 更新HTTP方法统计
function updateMethodStats(stats) {
    const container = document.querySelector('.detail-section:nth-child(4)');
    const methodContainer = container.querySelector('.method-bar-container') ? container : container;
    
    // 移除旧的方法统计
    methodContainer.querySelectorAll('.method-bar-container').forEach(el => el.remove());
    
    // 添加新的方法统计
    Object.entries(stats.method_calls).forEach(([method, count]) => {
        const percentage = calculatePercentage(count, stats.total_calls);
        const barContainer = document.createElement('div');
        barContainer.className = 'method-bar-container';
        
        barContainer.innerHTML = `
            <div class="method-info">
                <span class="method-name">${method}</span>
                <span class="method-count">${count}</span>
            </div>
            <div class="method-bar">
                <div class="method-bar-fill" style="width: ${percentage}%"></div>
            </div>
        `;
        
        methodContainer.appendChild(barContainer);
    });
}

// 更新API路径统计表格
function updatePathStatsTable(stats) {
    const tbody = document.querySelector('#path-stats-table tbody');
    
    // 移除旧的表格行
    tbody.innerHTML = '';
    
    // 添加新的表格行
    Object.entries(stats.path_calls).forEach(([path, count]) => {
        const percentage = calculatePercentage(count, stats.total_calls);
        const row = document.createElement('tr');
        
        row.innerHTML = `
            <td>${path}</td>
            <td class="count-column">${count}</td>
            <td>${percentage}%</td>
            <td>
                <div class="progress path-progress">
                    <div class="progress-bar" role="progressbar" style="width: ${percentage}%;" aria-valuenow="${percentage}" aria-valuemin="0" aria-valuemax="100"></div>
                </div>
            </td>
        `;
        
        tbody.appendChild(row);
    });
    
    // 重新应用排序（如果有）
    const sortHeaders = document.querySelectorAll('.sortable');
    sortHeaders.forEach(header => {
        if (header.classList.contains('asc')) {
            header.click();
            header.click();
        }
    });
}

// 更新API路径统计图表
function updatePathStatsChart(stats) {
    const ctx = document.getElementById('path-stats-chart');
    if (!ctx) return;
    
    // 准备数据
    let pathStats = Object.entries(stats.path_calls).map(([path, count]) => ({
        path: path,
        count: count
    }));
    
    // 按调用次数排序
    pathStats.sort((a, b) => b.count - a.count);
    
    // 最多显示10个路径，其他合并为"其他"
    let topPaths = pathStats.slice(0, 10);
    if (pathStats.length > 10) {
        const othersCount = pathStats.slice(10).reduce((sum, item) => sum + item.count, 0);
        topPaths.push({ path: "其他", count: othersCount });
    }
    
    const labels = topPaths.map(item => item.path);
    const data = topPaths.map(item => item.count);
    
    // 生成渐变色
    const gradient = ctx.getContext('2d').createLinearGradient(0, 0, 0, 300);
    gradient.addColorStop(0, 'rgba(102, 126, 234, 0.8)');
    gradient.addColorStop(1, 'rgba(102, 126, 234, 0.2)');
    
    // 如果图表已存在，更新数据
    if (pathStatsChart) {
        pathStatsChart.data.labels = labels;
        pathStatsChart.data.datasets[0].data = data;
        pathStatsChart.update();
    } else {
        // 创建新图表
        pathStatsChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: '调用次数',
                    data: data,
                    backgroundColor: gradient,
                    borderColor: 'rgba(102, 126, 234, 1)',
                    borderWidth: 2,
                    borderRadius: 5,
                    barPercentage: 0.6,
                    categoryPercentage: 0.6
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        padding: 12,
                        titleFont: {
                            size: 14
                        },
                        bodyFont: {
                            size: 13
                        },
                        cornerRadius: 8
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        grid: {
                            color: 'rgba(0, 0, 0, 0.05)'
                        },
                        ticks: {
                            callback: function(value) {
                                return value.toLocaleString();
                            }
                        }
                    },
                    x: {
                        grid: {
                            display: false
                        },
                        ticks: {
                            maxRotation: 45,
                            minRotation: 45,
                            callback: function(value) {
                                const label = this.getLabelForValue(value);
                                return label.length > 15 ? label.substring(0, 15) + '...' : label;
                            }
                        }
                    }
                },
                animation: {
                    duration: 1000,
                    easing: 'easeOutQuart'
                }
            }
        });
    }
}

// 更新最近调用记录
function updateLastCallDetails(stats) {
    const tbody = document.querySelector('.detail-section:last-child tbody');
    
    // 移除旧的表格行
    tbody.innerHTML = '';
    
    // 添加新的表格行
    stats.last_call_details.forEach(detail => {
        const date = new Date(detail.timestamp);
        const formattedTime = date.toLocaleString('zh-CN');
        let statusClass = 'status-500';
        
        if (detail.status_code === 200) {
            statusClass = 'status-200';
        } else if (detail.status_code === 400) {
            statusClass = 'status-400';
        } else if (detail.status_code === 404) {
            statusClass = 'status-404';
        }
        
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${formattedTime}</td>
            <td>${detail.path}</td>
            <td><span class="badge bg-primary">${detail.method}</span></td>
            <td>${detail.ip}</td>
            <td><span class="status-badge ${statusClass}">${detail.status_code}</span></td>
        `;
        
        tbody.appendChild(row);
    });
}

// 刷新统计数据
async function refreshStats() {
    try {
        const response = await fetch('/api/stats');
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        const stats = await response.json();
        
        // 更新各个部分
        updateStatCards(stats);
        updateMethodStats(stats);
        updatePathStatsTable(stats);
        updatePathStatsChart(stats);
        updateLastCallDetails(stats);
        
        console.log('Stats refreshed successfully');
    } catch (error) {
        console.error('Failed to refresh stats:', error);
    }
}

// 定期刷新统计数据
setInterval(refreshStats, 5000); // 每5秒刷新一次

// 表格排序功能
document.addEventListener('DOMContentLoaded', function() {
    // 获取排序表头
    const sortableHeaders = document.querySelectorAll('.sortable');
    
    sortableHeaders.forEach(header => {
        header.addEventListener('click', function() {
            const table = this.closest('table');
            const tbody = table.querySelector('tbody');
            const rows = Array.from(tbody.querySelectorAll('tr'));
            const sortType = this.getAttribute('data-sort');
            const isAscending = this.classList.contains('asc');
            
            // 移除所有排序图标
            document.querySelectorAll('.sortable i').forEach(icon => {
                icon.className = 'fa fa-sort';
            });
            
            // 设置当前排序图标
            const icon = this.querySelector('i');
            if (isAscending) {
                icon.className = 'fa fa-sort-desc';
            } else {
                icon.className = 'fa fa-sort-asc';
            }
            
            // 排序行
            rows.sort((a, b) => {
                let aValue, bValue;
                
                if (sortType === 'count') {
                    aValue = parseInt(a.querySelector('.count-column').textContent);
                    bValue = parseInt(b.querySelector('.count-column').textContent);
                }
                
                if (isAscending) {
                    return bValue - aValue;
                } else {
                    return aValue - bValue;
                }
            });
            
            // 更新表格
            rows.forEach(row => tbody.appendChild(row));
            
            // 切换排序方向
            this.classList.toggle('asc');
        });
    });
    
    // 初始化图表
    fetch('/api/stats')
        .then(response => response.json())
        .then(stats => {
            updatePathStatsChart(stats);
        })
        .catch(error => {
            console.error('Failed to initialize chart:', error);
        });
});