// 事件
window.addEventListener('DOMContentLoaded', () => {
    fetchData();
});

// 查询规则组
function fetchData() {
    fetch('/api/v1/ruleGroup/select')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('rule-group-data');
            tableBody.innerHTML = '';

            data.data.forEach(item => {
                const row = document.createElement('tr');

                const checkboxCell = document.createElement('td');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkboxCell.appendChild(checkbox);
                row.appendChild(checkboxCell);

                const nameCell = document.createElement('td');
                const nameButton = document.createElement('a');
                nameButton.textContent = item.Name;
                nameButton.classList.add('rule-group-link');
                nameButton.href = '/ruleGroup/' + item.Name+'/rule?ruleGroupName=' + encodeURIComponent(item.Name);
                nameCell.appendChild(nameButton);
                row.appendChild(nameCell);

                const numCell = document.createElement('td');
                numCell.textContent = item.RuleNumber;
                row.appendChild(numCell);

                const createTimeCell = document.createElement('td');
                createTimeCell.textContent = item.Description;
                row.appendChild(createTimeCell);

                const descriptionCell = document.createElement('td');
                descriptionCell.textContent = item.CreatedAt;
                row.appendChild(descriptionCell);

                tableBody.appendChild(row);
            });

            // 全选复选框事件监听器
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

// 打开创建规则组弹窗
function openCreateFormRuleGroup() {
    document.getElementById("create-ruleGroup-form-overlay").style.display = "block";
}

// 关闭创建规则组弹窗
function cancelCreateFormRuleGroup() {
    document.getElementById("create-ruleGroup-form-overlay").style.display = "none";
}

// 提交创建规则组表单
function submitCreateFormRuleGroup() {

    // 获取输入框的值
    var name = document.getElementById("create-ruleGroupName-input").value;
    var description = document.getElementById("create-ruleGroupDescription-input").value;

    // 构建JSON对象
    var data = {
        name: name,
        description: description
    };

    // 发送POST请求给后端
    fetch("/api/v1/ruleGroup/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.ok) {
                alert("规则组创建成功！");

                // 清空输入框的值
                document.getElementById("create-ruleGroupName-input").value = "";
                document.getElementById("create-ruleGroupDescription-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("规则组创建失败，请重试！");
            }
        })

    // 提交成功后关闭弹窗
    cancelCreateFormRuleGroup();
}

// 打开更新规则组弹窗
function openUpdateFormRuleGroup() {

    var updateTableBody = document.getElementById("rule-group-data");
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

    // 发送 GET 请求获取规则数据
    fetch('/api/v1/ruleGroup/getRuleGroup?ruleGroupName=' + encodeURIComponent(ruleName))
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
            document.getElementById('update-ruleGroupName-input').value = ruleData.data[0].Name;
            document.getElementById('update-ruleGroupDescription-input').value = ruleData.data[0].Description;

            // 显示更新弹窗
            var updateFormOverlay = document.getElementById('update-form-overlay');
            updateFormOverlay.style.display = 'block';

        })
        .catch(error => {
            alert(error.message);
        });


}

// 关闭更新规则组弹窗
function cancelUpdateFormRuleGroup() {
    document.getElementById("update-form-overlay").style.display = "none";
}

// 提交更新规则组表单
function submitUpdateFormRuleGroup() {

    // 构造更新数据对象
    var updateData = {
        name: document.getElementById('update-ruleGroupName-input').value,
        description: document.getElementById('update-ruleGroupDescription-input').value
    };

    console.log(JSON.stringify(updateData))

    // 发送 POST 请求更新数据
    fetch('/api/v1/ruleGroup/update?ruleGroupName='+encodeURIComponent(document.getElementById('update-ruleGroupName-input').value), {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
    })
        .then(response => {
            if (response.ok) {
                alert("规则更新成功！");
                cancelUpdateFormRuleGroup();
                // 刷新当前页面
                location.reload();
            } else {
                alert("规则更新失败，请重试！");
            }
        });
}

// 删除规则组
function deleteRowRuleGroup() {

    var ruleGroupData = document.getElementById("rule-group-data");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = ruleGroupData.getElementsByTagName("input");
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

        // 发送DELETE请求给后端
        fetch("/api/v1/ruleGroup/delete?ruleGroupName="+encodeURIComponent(name), {
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