// 事件
window.addEventListener('DOMContentLoaded', () => {
    fetchData();
});

// 查询规则
function fetchData() {
    var currentURI = window.location.pathname;
    fetch('/api/v1' + currentURI + '/select')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('table-body');
            tableBody.innerHTML = '';

            data.data.ResRule.forEach(item => {

                const row = document.createElement('tr');

                const checkboxCell = document.createElement('td');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkboxCell.appendChild(checkbox);
                row.appendChild(checkboxCell);

                const nameCell = document.createElement('td');
                nameCell.textContent = item.alert;
                row.appendChild(nameCell);

                const exprCell = document.createElement('td');
                exprCell.textContent = item.expr;
                row.appendChild(exprCell);

                const forCell = document.createElement('td');
                forCell.textContent = item.for;
                row.appendChild(forCell);

                const labelsCell = document.createElement('td');
                labelsCell.textContent = JSON.stringify(item.labels);
                row.appendChild(labelsCell);

                const annotationsCell = document.createElement('td');
                annotationsCell.textContent = item.annotations.description;
                row.appendChild(annotationsCell);

                tableBody.appendChild(row);
            });

            const selectAllCheckbox = document.getElementById('select-all-checkbox');
            selectAllCheckbox.addEventListener('change', () => {
                const checkboxes = document.querySelectorAll('tbody input[type="checkbox"]');
                checkboxes.forEach(checkbox => {
                    checkbox.checked = selectAllCheckbox.checked;
                });
            });
        })
        .catch(error => {
            console.log('Error:', error);
        });
}

// 创建规则
function submitCreateForm() {

    // 获取输入框的值
    var name = document.getElementById("name-input").value;
    var expr = document.getElementById("expr-input").value;
    var interval = document.getElementById("for-input").value;
    var label = document.getElementById("label-input").value;
    var description = document.getElementById("annotations-input").value;

    const urlParams = new URLSearchParams(window.location.search);
    const ruleGroupName = urlParams.get('ruleGroupName');

    // 构建JSON对象
    var data = {
        groups: [
            {
                name: ruleGroupName,
                rules: [
                    {
                        alert: name,
                        expr: expr,
                        for: interval,
                        labels: formatTextToJSON(label),
                        annotations: {
                            description: description
                        }
                    }
                ]
            }
        ]
    };

    // 发送POST请求给后端
    var currentURI = window.location.pathname;
    fetch("/api/v1"+currentURI+"/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.ok) {
                alert("规则创建成功！");
                cancelCreateForm(); // 取消弹窗

                // 清空输入框的值
                document.getElementById("name-input").value = "";
                document.getElementById("expr-input").value = "";
                document.getElementById("for-input").value = "";
                document.getElementById("label-input").value = "";
                document.getElementById("annotations-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("规则创建失败，请重试！");
            }
        })
        .catch(error => {
            console.error("Error:", error);
            alert("规则创建失败，请重试！");
        });

    // 提交表单后，隐藏窗口
    cancelCreateForm();
}

// 删除规则
function deleteRow() {
    var tableBody = document.getElementById("table-body");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = tableBody.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row);
        }
    }

    // 确认删除操作
    if (selectedRows.length === 0) {
        alert("请选择要删除的规则！");
        return;
    }

    var confirmMessage =
        "您确定要删除选中的 " + selectedRows.length + " 条规则吗？";
    if (!confirm(confirmMessage)) {
        return;
    }

    // 获取选中行的规则名称并发送请求给后端
    selectedRows.forEach(row => {
        var name = row.cells[1].innerText; // 根据表格结构调整索引

        const urlParams = new URLSearchParams(window.location.search);
        const ruleGroupName = urlParams.get('ruleGroupName');
        // 发送DELETE请求给后端
        fetch("/api/v1/ruleGroup/"+ruleGroupName + "/rule/delete?ruleName="+encodeURIComponent(name), {
            method: "POST"
        })
            .then(response => {
                if (response.ok) {
                    alert("规则删除成功！");
                    row.remove(); // 从表格中删除行
                } else {
                    alert("规则删除失败，请重试！");
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("规则删除失败，请重试！");
            });
    });
}

// 更新规则
function submitUpdateForm() {

    console.log(document.getElementById('update-label-input').value)

    // 构造更新数据对象
    var updateData = {
        groups: [
            {
                name: "rules",
                rules:[
                    {
                        alert: document.getElementById('update-name-input').value,
                        expr: document.getElementById('update-expr-input').value,
                        for: document.getElementById('update-for-input').value,
                        labels: formatTextToJSON(document.getElementById('update-label-input').value),
                        annotations: {
                            description: document.getElementById('update-annotations-input').value
                        }
                    }
                ]
            }
        ]
    };

    console.log(JSON.stringify(updateData))
    const urlParams = new URLSearchParams(window.location.search);
    const ruleGroupName = urlParams.get('ruleGroupName');
    // 发送 POST 请求更新数据
    fetch('/api/v1/ruleGroup/'+ruleGroupName+'/rule/update?ruleName='+encodeURIComponent(document.getElementById('update-name-input').value), {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
    })
        .then(response => {
            if (response.ok) {
                alert("规则更新成功！");
                cancelUpdateForm(); // 取消弹窗
                // 刷新当前页面
                location.reload();
            } else {
                alert("规则更新失败，请重试！");
            }
        });
}

// 创建规则的弹窗
function openCreateForm() {
    var overlay = document.getElementById("create-form-overlay");
    overlay.style.display = "block";
}

// 取消创建的规则弹窗
function cancelCreateForm() {
    var overlay = document.getElementById("create-form-overlay");
    overlay.style.display = "none";
}

// 更新规则的弹窗
function openUpdateForm() {

    var updateTableBody = document.getElementById("table-body");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = updateTableBody.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row);
        }
    }

    // 确认删除操作
    if (selectedRows.length !== 1) {
        alert("请选择要更新的规则, 并且只能选择 1 个!");
        return;
    }

    // 获取规则名称
    var ruleName = row.cells[1].innerText; // 根据表格结构调整索引
    const urlParams = new URLSearchParams(window.location.search);
    const ruleGroupName = urlParams.get('ruleGroupName');
    // 发送 GET 请求获取规则数据
    fetch('/api/v1/ruleGroup/'+ruleGroupName+'/rule/getRule?ruleName=' + encodeURIComponent(ruleName))
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('获取规则数据失败。');
            }
        })
        .then(ruleData => {
            // 将获取到的数据填充到表单中
            console.log(ruleData);
            document.getElementById('update-name-input').value = ruleData.data.alert;
            document.getElementById('update-expr-input').value = ruleData.data.expr;
            document.getElementById('update-for-input').value = ruleData.data.for;
            document.getElementById('update-label-input').value = formatJSONToText(ruleData.data.labels);
            document.getElementById('update-annotations-input').value = ruleData.data.annotations.description;

            // 显示更新弹窗
            var updateFormOverlay = document.getElementById('update-form-overlay');
            updateFormOverlay.style.display = 'block';

        })
        .catch(error => {
            alert(error.message);
        });
}

// 取消更新的规则弹窗
function cancelUpdateForm() {
    // 关闭更新弹窗
    var updateFormOverlay = document.getElementById('update-form-overlay');
    updateFormOverlay.style.display = 'none';
}

// 文本转换 json
// xx=xx -> {"xx":"xx"}
function formatTextToJSON(text) {
    var keyValuePairs = text.split(',');
    var jsonObject = {};

    for (var i = 0; i < keyValuePairs.length; i++) {
        var pair = keyValuePairs[i].trim().split('=');
        var key = pair[0].trim();
        var value = pair[1].trim();

        try {
            // 尝试解析值为 JSON 对象
            value = JSON.parse(value);
        } catch (error) {
            // 解析失败时，将值保留为字符串
        }

        jsonObject[key] = value;
    }

    return jsonObject
}

// json 转换文本
// {"xx":"xx"} -> xx=xx
function formatJSONToText(json) {
    var text = '';

    for (var key in json) {
        if (json.hasOwnProperty(key)) {
            var value = json[key];
            text += key + '=' + value + ',';
        }
    }

    // 移除最后一个逗号
    text = text.slice(0, -1);

    return text;
}