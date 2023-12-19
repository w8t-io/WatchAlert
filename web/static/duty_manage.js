// 事件
window.addEventListener('DOMContentLoaded', () => {
    fetchData();
});

var dutyManageData = document.getElementById("duty-manage-data");

// 查询所有对象
function fetchData() {
    fetch('/api/v1/dutyManage/list')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('duty-manage-data');
            tableBody.innerHTML = '';

            data.data.forEach(item => {
                const row = document.createElement('tr');

                const checkboxCell = document.createElement('td');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkboxCell.appendChild(checkbox);
                row.appendChild(checkboxCell);

                const idCell = document.createElement('td');
                const idButton = document.createElement('a');
                idButton.textContent = item.id;
                idButton.classList.add('duty-schedule-link');
                idButton.href = '/dutyManage/' + item.id+'/schedule';
                idCell.appendChild(idButton);
                row.appendChild(idCell);

                const nameCell = document.createElement('td');
                nameCell.textContent = item.name;
                row.appendChild(nameCell);

                const descriptionCell = document.createElement('td');
                descriptionCell.textContent = item.description;
                row.appendChild(descriptionCell);

                tableBody.appendChild(row);
            });

            // 全选复选框事件监听器
            const selectAllCheckbox = document.getElementById('dutyManage-select-all-checkbox');
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

// 打开创建对象弹窗
function openCreateFormDutyManage() {
    document.getElementById("create-dutyManage-form-overlay").style.display = "block";
}

// 关闭创建对象弹窗
function cancelCreateFormDutyManage() {
    document.getElementById("create-dutyManage-form-overlay").style.display = "none";
}

// 提交创建对象表单
function submitCreateFormDutyManage() {

    // 获取输入框的值
    var name = document.getElementById("create-dutyManageName-input").value;
    var description = document.getElementById("create-dutyManageDescription-input").value;

    // 构建JSON对象
    var data = {
        name: name,
        description: description,
    };

    // 发送POST请求给后端
    fetch("/api/v1/dutyManage/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.ok) {
                alert("创建成功！");

                // 清空输入框的值
                document.getElementById("create-dutyManageName-input").value = "";
                document.getElementById("create-dutyManageDescription-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("创建失败，请重试！");
            }
        })

    // 提交成功后关闭弹窗
    cancelCreateFormDutyManage();
}

// 打开更新对象弹窗
function openUpdateFormDutyManage() {

    var updateTableBody = document.getElementById("duty-manage-data");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = updateTableBody.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row.cells[0].nextElementSibling.innerText);
        }
    }

    // 确认删除操作
    if (selectedRows.length !== 1) {
        alert("请选择要更新的对象, 并且只能选择 1 个!");
        return;
    }

    // 发送 GET 请求获取规则数据
    fetch('/api/v1/dutyManage/get?id=' + encodeURIComponent(selectedRows[0]))
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('获取数据失败。');
            }
        })
        .then(dutyManageData => {
            console.log(dutyManageData)
            // 将获取到的数据填充到表单中
            document.getElementById("update-dutyManageName-input").value = dutyManageData.data.name;
            document.getElementById("update-dutyManageDescription-input").value = dutyManageData.data.description;

            // 显示更新弹窗
            var updateFormOverlay = document.getElementById('update-dutyManage-form-overlay');
            updateFormOverlay.style.display = 'block';
        })

}

// 关闭更新对象弹窗
function cancelUpdateFormDutyManage() {
    document.getElementById("update-dutyManage-form-overlay").style.display = "none";
}

// 提交更新对象表单
function submitUpdateFormDutyManage() {

    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = dutyManageData.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row.cells[0].nextElementSibling.innerText);
        }
    }

    // 获取输入框的值
    var name = document.getElementById("update-dutyManageName-input").value;
    var description = document.getElementById("update-dutyManageDescription-input").value;

    // 构建JSON对象
    var data = {
        id: selectedRows[0],
        name: name,
        description: description,
    };

    console.log(data)

    // 发送POST请求给后端
    fetch("/api/v1/dutyManage/update", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.ok) {
                alert("更新成功！");

                // 清空输入框的值
                document.getElementById("update-dutyManageName-input").value = "";
                document.getElementById("update-dutyManageDescription-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("更新失败，请重试！");
            }
        })

    cancelUpdateFormDutyManage()
}

// 删除对象
function deleteDutyManage() {

    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = dutyManageData.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row.cells[0].nextElementSibling.innerText);
        }
    }

    console.log(selectedRows)

    // 确认删除操作
    if (selectedRows.length === 0) {
        alert("请选择要删除的记录！");
        return;
    }

    var confirmMessage =
        "您确定要删除选中的 " + selectedRows.length + " 条记录吗？";
    if (!confirm(confirmMessage)) {
        return;
    }

    // 获取选中行的规则名称并发送请求给后端
    selectedRows.forEach(row => {
        for (var i = 0; i < selectedRows.length; i++) {
            // 发送DELETE请求给后端
            fetch("/api/v1/dutyManage/delete?id="+encodeURIComponent(selectedRows[i]), {
                method: "POST"
            })
                .then(response => {
                    if (response.ok) {
                        alert("删除成功！");
                        // 刷新当前页面
                        location.reload();
                    } else {
                        alert("删除失败，请重试！");
                    }
                })
        }
        selectedRows[i]=""
    });

}