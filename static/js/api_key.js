// API Key管理功能
document.addEventListener('DOMContentLoaded', function() {
    // 加载API Key列表
    loadApiKeys();

    // 表单提交事件
    document.getElementById('createApiKeyForm').addEventListener('submit', function(e) {
        e.preventDefault();
        createApiKey();
    });
});

// 加载API Key列表
function loadApiKeys() {
    fetch('/auth/api_key')
        .then(response => response.json())
        .then(data => {
            renderApiKeys(data.api_keys);
        })
        .catch(error => {
            console.error('加载API Key失败:', error);
            document.getElementById('apiKeysContainer').innerHTML = '<div class="alert alert-danger">加载API Key失败，请刷新页面重试</div>';
        });
}

// 渲染API Key列表
function renderApiKeys(apiKeys) {
    const container = document.getElementById('apiKeysContainer');
    
    if (apiKeys.length === 0) {
        container.innerHTML = '<div class="alert alert-info">暂无API Key，请创建新的API Key</div>';
        return;
    }

    let html = '';
    apiKeys.forEach(key => {
        const usagePercentage = key.is_permanent ? 0 : Math.min(100, (key.current_usage / key.max_usage) * 100);
        const badgeClass = key.is_permanent ? 'badge-permanent' : 'badge-limited';
        const badgeText = key.is_permanent ? '永久有效' : `限制使用 ${key.max_usage} 次`;
        
        html += `
            <div class="api-key-item">
                <div class="row">
                    <div class="col-md-8">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5>${key.name}</h5>
                            <span class="badge ${badgeClass}">${badgeText}</span>
                        </div>
                        <div class="mb-2">
                            <strong>API Key:</strong>
                            <div class="api-key-value mt-1">${key.key}</div>
                        </div>
                        <div class="mb-1">
                            <div class="d-flex justify-content-between">
                                <span>使用情况:</span>
                                <span>${key.current_usage}/${key.max_usage}</span>
                            </div>
                            <div class="usage-bar">
                                <div class="usage-fill" style="width: ${usagePercentage}%"></div>
                            </div>
                        </div>
                        <div class="text-muted" style="font-size: 0.85rem;">
                            创建时间: ${new Date(key.created_at).toLocaleString()}
                        </div>
                    </div>
                    <div class="col-md-4 d-flex align-items-center justify-content-end">
                        <button class="btn btn-danger btn-sm" onclick="deleteApiKey(${key.id}, '${key.name}')">
                            <i class="fa fa-trash mr-1"></i>
                            删除
                        </button>
                    </div>
                </div>
            </div>
        `;
    });
    
    container.innerHTML = html;
}

// 创建API Key
function createApiKey() {
    const form = document.getElementById('createApiKeyForm');
    const formData = new FormData(form);
    const data = {
        name: formData.get('name'),
        max_usage: parseInt(formData.get('max_usage')),
        is_permanent: formData.get('is_permanent') === 'on'
    };

    fetch('/auth/api_key', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(result => {
        if (result.api_key) {
            // 关闭模态框
            const modal = bootstrap.Modal.getInstance(document.getElementById('createApiKeyModal'));
            modal.hide();
            
            // 显示成功消息
            showSuccessMessage('API Key创建成功');
            
            // 重置表单
            form.reset();
            
            // 重新加载API Key列表
            loadApiKeys();
        } else {
            alert('创建失败: ' + (result.error || '未知错误'));
        }
    })
    .catch(error => {
        console.error('创建API Key失败:', error);
        alert('创建API Key失败，请重试');
    });
}

// 删除API Key
function deleteApiKey(id, name) {
    if (confirm(`确定要删除API Key "${name}"吗？此操作不可恢复。`)) {
        fetch(`/auth/api_key/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(result => {
            if (result.message) {
                // 显示成功消息
                showSuccessMessage('API Key删除成功');
                
                // 重新加载API Key列表
                loadApiKeys();
            } else {
                alert('删除失败: ' + (result.error || '未知错误'));
            }
        })
        .catch(error => {
            console.error('删除API Key失败:', error);
            alert('删除API Key失败，请重试');
        });
    }
}

// 显示成功消息
function showSuccessMessage(message) {
    const modal = new bootstrap.Modal(document.getElementById('successModal'));
    document.getElementById('successModalMessage').textContent = message;
    modal.show();
}