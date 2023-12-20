// 事件
window.addEventListener('DOMContentLoaded', () => {
    fetchData();
});

// 查询所有对象
function fetchData() {
    fetch('/api/v1/alertNotice/list')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('notice-object-data');
            tableBody.innerHTML = '';

            data.data.forEach(item => {
                const row = document.createElement('tr');

                const checkboxCell = document.createElement('td');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkboxCell.appendChild(checkbox);
                row.appendChild(checkboxCell);

                const uuidCell = document.createElement('td');
                uuidCell.textContent = item.uuid;
                row.appendChild(uuidCell);

                const nameCell = document.createElement('td');
                nameCell.textContent = item.name;
                row.appendChild(nameCell);

                const envCell = document.createElement('td');
                envCell.textContent = item.env;
                row.appendChild(envCell);

                const statusCell = document.createElement('td');
                statusCell.textContent = item.noticeStatus;
                row.appendChild(statusCell);

                const dataSourceCell = document.createElement('td');
                dataSourceCell.textContent = item.dataSource;
                row.appendChild(dataSourceCell);

                const noticeCell = document.createElement('td');
                noticeCell.textContent = item.noticeType;
                row.appendChild(noticeCell);

                const testCell = document.createElement('td');
                const testButton = document.createElement('a');
                testButton.textContent = "Test";
                testButton.classList.add('test-notice-status-link');
                testButton.addEventListener('click', function() {
                    fetch('/api/v1/alertNotice/checkNoticeStatus?uuid='+item.uuid, {
                        method: 'GET',
                    })
                    location.reload()
                });
                testCell.appendChild(testButton);
                row.appendChild(testCell);

                tableBody.appendChild(row);
            });

            // 全选复选框事件监听器
            const selectAllCheckbox = document.getElementById('notice-select-all-checkbox');
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
function openCreateFormNoticeObject() {
    document.getElementById("create-noticeObject-form-overlay").style.display = "block";
}

// 关闭创建对象弹窗
function cancelCreateFormNoticeObject() {
    document.getElementById("create-noticeObject-form-overlay").style.display = "none";
}

// 提交创建对象表单
function submitCreateFormNoticeObject() {

    // 获取输入框的值
    var name = document.getElementById("create-noticeObjectName-input").value;
    var env = document.getElementById("create-noticeObjectEnv-input").value;
    var dataSource = document.getElementById("create-noticeObjectDataSource-input").value;
    var noticeType = document.querySelector('input[name="notice-type"]:checked').value;
    var feishuChatId = document.getElementById("feiShuChatID").value;
    var dutyId = document.getElementById("create-noticeObjectDutyId-input").value;

    // 构建JSON对象
    var data = {
        name: name,
        env: env,
        noticeStatus: "未知",
        dataSource: dataSource,
        noticeType: noticeType,
        feishuChatId: feishuChatId,
        dutyId: dutyId
    };

    // 发送POST请求给后端
    fetch("/api/v1/alertNotice/create", {
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
                document.getElementById("create-noticeObjectName-input").value = "";
                document.getElementById("create-noticeObjectEnv-input").value = "";
                document.getElementById("create-noticeObjectDataSource-input").value = "";
                document.querySelector('input[name="notice-type"]:checked').value = "";
                document.getElementById("feiShuChatID").value = "";
                document.getElementById("create-noticeObjectDutyId-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("创建失败，请重试！");
            }
        })

    // 提交成功后关闭弹窗
    cancelCreateFormRuleGroup();
}

// 打开更新对象弹窗
function openUpdateFormNoticeObject() {

    var updateTableBody = document.getElementById("notice-object-data");
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
    fetch('/api/v1/alertNotice/get?uuid=' + encodeURIComponent(selectedRows[0]))
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('获取数据失败。');
            }
        })
        .then(noticeData => {
            // 将获取到的数据填充到表单中
            document.getElementById("update-noticeObjectName-input").value = noticeData.data.name;
            document.getElementById("update-noticeObjectEnv-input").value = noticeData.data.env;
            document.getElementById("update-noticeObjectDataSource-input").value = noticeData.data.dataSource;
            if (noticeData.data.noticeType === "FeiShu") {
                document.getElementById("update-notice-FeiShu").checked = true;
                document.getElementById("update-chat-id-row").style.display = "flex";
                document.getElementById("update-feiShuChatID").value = noticeData.data.feishuChatId;
            }
            document.getElementById("update-noticeObjectDutyId-input").value = noticeData.data.dutyId;

            // 显示更新弹窗
            var updateFormOverlay = document.getElementById('update-noticeObject-form-overlay');
            updateFormOverlay.style.display = 'block';
        })

}

// 关闭更新对象弹窗
function cancelUpdateFormNoticeObject() {
    document.getElementById("update-noticeObject-form-overlay").style.display = "none";
}

// 提交更新对象表单
function submitUpdateFormNoticeObject() {

    var noticeObjectData = document.getElementById("notice-object-data");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = noticeObjectData.getElementsByTagName("input");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].type === "checkbox" && checkboxes[i].checked) {
            var row = checkboxes[i].parentNode.parentNode;
            selectedRows.push(row.cells[0].nextElementSibling.innerText);
        }
    }

    // 获取输入框的值
    var name = document.getElementById("update-noticeObjectName-input").value;
    var env = document.getElementById("update-noticeObjectEnv-input").value;
    var dataSource = document.getElementById("update-noticeObjectDataSource-input").value;
    var noticeType = document.querySelector('input[name="update-notice-type"]:checked').value;
    var feishuChatId = document.getElementById("update-feiShuChatID").value.trim();
    var dutyId = document.getElementById("update-noticeObjectDutyId-input").value.trim();
    if (dutyId.length === 0){
        dutyId = " "
    }

    // 构建JSON对象
    var data = {
        uuid: selectedRows[0],
        name: name,
        env: env,
        dataSource: dataSource,
        noticeType: noticeType,
        feishuChatId: feishuChatId,
        dutyId: dutyId
    };

    // 发送POST请求给后端
    fetch("/api/v1/alertNotice/update", {
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
                document.getElementById("update-noticeObjectName-input").value = "";
                document.getElementById("update-noticeObjectEnv-input").value = "";
                document.getElementById("update-noticeObjectDataSource-input").value = "";
                document.querySelector('input[name="update-notice-type"]:checked').value = "";
                document.getElementById("update-feiShuChatID").value = "";
                document.getElementById("update-noticeObjectDutyId-input").value = "";

                // 刷新当前页面
                location.reload();

            } else {
                alert("更新失败，请重试！");
            }
        })

    cancelUpdateFormNoticeObject()
}

// 删除对象
function deleteNoticeObject() {

    var noticeObjectData = document.getElementById("notice-object-data");
    var selectedRows = [];

    // 查找被选中的行
    var checkboxes = noticeObjectData.getElementsByTagName("input");
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
            fetch("/api/v1/alertNotice/delete?uuid="+encodeURIComponent(selectedRows[i]), {
                method: "POST"
            })
                .then(response => {
                    if (response.ok) {
                        alert("删除成功！");
                        row.remove(); // 从表格中删除行
                        // 刷新当前页面
                        location.reload();
                    } else {
                        alert("删除失败，请重试！");
                    }
                })
        }
    });

}

// NoticeType 单选项
function toggleFeiShuChatIdInput(selectedNotice) {
    var feishuUserIdRow = document.getElementById("feishu-chat-id-row");
    var dingdingUserIdRow = document.getElementById("dingding-chat-id-row");

    if (selectedNotice === "FeiShu") {
        feishuUserIdRow.style.display = "block";
        dingdingUserIdRow.style.display = "none";
    } else if (selectedNotice === "dingding") {
        feishuUserIdRow.style.display = "none";
        dingdingUserIdRow.style.display = "block";
    } else {
        feishuUserIdRow.style.display = "none";
        dingdingUserIdRow.style.display = "none";
    }
}