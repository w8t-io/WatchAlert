const calendarElement = document.getElementById('calendar');
// 获取今天的日期
const today = new Date();
// 用户按钮
const userListContainer = document.getElementById('user-list');
const UpdateUserListContainer = document.getElementById('update-user-list');
const publishBtn = document.getElementById('publish-btn');

// 选择的用户数组
const selectedUsers = [];
let isChecked = ""
let userName = ""
let phone = ""
let userId = ""


// 获取值班数据
async function getDutyData(date) {
    const response = await fetch('/api/v1/dutySystem/select?time='+date);
    return await response.json();
}

// 生成日历
function generateCalendar() {
    // 获取当前月份的第一天
    const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
    let month = firstDay.toLocaleString('default', { month: 'long' });

    // 创建月份标题
    const monthElement = document.createElement('div');
    monthElement.classList.add('month');
    monthElement.textContent = month;
    calendarElement.appendChild(monthElement);

    // 创建星期标题
    const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    const daysOfWeekElement = document.createElement('div');
    daysOfWeekElement.classList.add('days');
    daysOfWeek.forEach(day => {
        const dayElement = document.createElement('div');
        dayElement.classList.add('day');
        dayElement.textContent = day;
        daysOfWeekElement.appendChild(dayElement);
    });
    calendarElement.appendChild(daysOfWeekElement);

    // 创建日期格子
    const daysInMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0).getDate();
    const daysElement = document.createElement('div');
    daysElement.classList.add('days');
    calendarElement.appendChild(daysElement);

    for (let i = 1; i <= daysInMonth; i++) {
        let date = new Date(today.getFullYear(), today.getMonth(), i);
        const dayElement = document.createElement('div');
        dayElement.classList.add('day');
        dayElement.textContent = i;

        // 如果是今天则添加 today 类名
        if (date.toDateString() === today.toDateString()) {
            dayElement.classList.add('today');
        }

        month = (today.getMonth()+1)
        date = today.getFullYear()+"-"+month+"-"+i

        // 获取值班数据并添加到日期格子下方
        getDutyData(date).then(info => {
            const dutyElement = document.createElement('div');
            dutyElement.classList.add('duty');
            dutyElement.textContent = info.data[0].userName;
            dayElement.appendChild(dutyElement);
        });

        // 点击日期格子时显示弹窗
        dayElement.addEventListener('click', () => {
            openModal(date, dayElement.getElementsByClassName('duty')[0].textContent);
        });

        daysElement.appendChild(dayElement);
    }
}

// 生成日历
generateCalendar();

let currentDutyUserName = '';
// 调用后端 API 接口来获取用户数据
fetch('/api/v1/dutySystem/user/select')
    .then(response => response.json())
    .then(data => {
        // 遍历用户数据并展示为列表项
        data.data.forEach(user => {
            const listItem = document.createElement('div');
            listItem.classList.add('user-item');

            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.value = user.ID;

            const username = document.createElement('span');
            username.textContent = user.userName;

            // 使用 data-* 属性存储 phone 和 userId 数据
            listItem.setAttribute('data-phone', user.phone);
            listItem.setAttribute('data-userid', user.userId);

            listItem.appendChild(username);
            const listItemCopy = listItem.cloneNode(true); // 复制listItem

            listItem.appendChild(checkbox);
            listItem.appendChild(username);
            userListContainer.appendChild(listItem);

            const radioBtn = document.createElement('input');
            radioBtn.type = 'radio';
            radioBtn.name = 'user';
            radioBtn.value = user.ID;
            listItemCopy.appendChild(radioBtn)
            UpdateUserListContainer.appendChild(listItemCopy);

            // 监听单选框的点击事件
            radioBtn.addEventListener('click', () => {
                currentDutyUserName = user.userName;
            });
        });
    })
    .catch(error => {
        console.error('无法获取用户数据:', error);
    });

// 监听用户选择
userListContainer.addEventListener('change', (event) => {
    isChecked = event.target.checked;
    userName = event.target.nextElementSibling.textContent;
    phone = event.target.nextElementSibling.parentElement.dataset.phone
    userId = event.target.nextElementSibling.parentElement.dataset.userid

    if (isChecked) {
        // 添加到已选用户列表
        selectedUsers.push({ userid: userId, userName: userName, phone: phone }); // 存储用户id和名称

    } else {
        // 从已选用户列表中移除
        const index = selectedUsers.findIndex(user => user.id === userId);
        if (index > -1) {
            selectedUsers.splice(index, 1);
        }
    }
});

// 监听发布按钮点击事件
publishBtn.addEventListener('click', () => {
    // 发送选中的用户给后端接口
    sendSelectedUsers(selectedUsers);
});

// 发布值班表，发送选中的用户给后端接口
function sendSelectedUsers(selectedUsers) {

    if (selectedUsers.length === 0) {
        alert("请选择要值班的用户!");
        return;
    }

    // 获取 <select> 元素
    const selectElement = document.getElementById('duty-user');
    // 获取当前选中选项的值
    const selectedValue = selectElement.value;

    fetch('/api/v1/dutySystem/create?dutyPeriod='+selectedValue, {
        method: 'POST',
        body: JSON.stringify(selectedUsers),
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log('已发送选中的用户:', data);
        })
        .catch(error => {
            console.error('发送选中的用户时出错:', error);
        });

    location.reload();
}

// 弹窗部分 ----
const modalElement = document.getElementById('modal');
const modalDateElement = document.getElementById('modal-date');
const modalFormElement = document.getElementById('modal-form');

// 显示弹窗
function openModal(date, name) {
    modalDateElement.textContent = date.toLocaleString('default', { month: 'long', day: 'numeric', year: 'numeric' });

    // 根据当前值班用户的名称选中相应的单选框
    const radioBtns = UpdateUserListContainer.querySelectorAll('input[type="radio"]');
    radioBtns.forEach(radioBtn => {
        const listItem = radioBtn.parentNode;
        const username = listItem.querySelector('span').textContent;
        if (username === name) {
            currentDutyUserName = username
            radioBtn.checked = true;
        }
    });

    modalElement.style.display = 'block';
}

// 关闭弹窗
function closeModal() {
    modalElement.style.display = 'none';
}

// 提交表单
function submitForm(event) {
    event.preventDefault();
    const date = new Date(modalDateElement.textContent);
    const name = currentDutyUserName;

    // 构建要发送的数据对象
    const options = { year: 'numeric', month: 'numeric', day: 'numeric' };
    const data = {
        time: date.toLocaleDateString('default', options).replace(/\//g, '-'),
        userName: name
    };

    console.log(data)

    // 发送 POST 请求
    fetch('/api/v1/dutySystem/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(result => {
            console.log('表单提交成功:', result);
            // 在这里可以根据后端返回的结果进行相应的处理
        })
        .catch(error => {
            console.error('表单提交出错:', error);
            // 在这里可以处理请求出错的情况
        });

    // 更新日历上的值班人员名
    const dayElement = event.target.parentNode;
    const dutyElement = dayElement.getElementsByClassName('duty')[0];
    dutyElement.textContent = name;

    closeModal();
}

// 关闭弹窗的按钮事件
modalElement.getElementsByClassName('close')[0].addEventListener('click', closeModal);

// 提交表单的按钮事件
modalFormElement.addEventListener('submit', submitForm);

// 创建规则的弹窗
function openCreateDutyUserForm() {
    var overlay = document.getElementById("create-duty-user-form-overlay");
    overlay.style.display = "block";
}

// 取消创建的规则弹窗
function cancelCreateDutyUserForm() {
    var overlay = document.getElementById("create-duty-user-form-overlay");
    overlay.style.display = "none";
}

function toggleUserIdInput(selectedNotice) {
    var feishuUserIdRow = document.getElementById("feishu-user-id-row");
    var dingdingUserIdRow = document.getElementById("dingding-user-id-row");

    if (selectedNotice === "feishu") {
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

function updateToggleUserIdInput(updateSelectedNotice) {
    var feishuUserIdRow = document.getElementById("update-feishu-user-id-row");
    var dingdingUserIdRow = document.getElementById("update-dingding-user-id-row");

    if (updateSelectedNotice === "feishu") {
        feishuUserIdRow.style.display = "block";
        dingdingUserIdRow.style.display = "none";
    } else if (updateSelectedNotice === "dingding") {
        feishuUserIdRow.style.display = "none";
        dingdingUserIdRow.style.display = "block";
    } else {
        feishuUserIdRow.style.display = "none";
        dingdingUserIdRow.style.display = "none";
    }
}

// 创建值班用户
function submitCreateForm() {
    var username = document.getElementById("username").value;
    var notice = document.querySelector('input[name="notice-type"]:checked').value;
    var feishuUserId = document.getElementById("feiShuUserID").value;
    var phone = document.getElementById("phone").value;
    var email = document.getElementById("email").value;

    var data = {
        userName: username,
        notice: notice,
        feiShuUserID: feishuUserId,
        phone: phone,
        email: email
    };

    console.log(data)

    // 发送数据给后端的API接口
    fetch("/api/v1/dutySystem/user/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(function(response) {
            // 处理响应
            if (response.ok) {
                console.log("数据发送成功");
                // 提取后端返回的userId
                return response.json();
            } else {
                console.error("数据发送失败");
                alert("创建失败");
            }
        })
        .then(function() {
            alert("创建成功");
            cancelCreateDutyUserForm();

            document.getElementById("username").value = "";
            document.querySelector('input[name="notice-type"]:checked').value = "";
            document.getElementById("feiShuUserID").value = "";
            document.getElementById("phone").value = "";
            document.getElementById("email").value = "";
        })
        .catch(function(error) {
            console.error("请求错误:", error);
            alert("创建失败 -> "+ error);
        });
}

// 更新用户弹窗打开时，将选中用户的信息显示在弹窗中
function openUpdateDutyUserForm() {

    if (userId.length === 0){
        alert("请选择要编辑的用户!");
        return;
    }

    fetch('/api/v1/dutySystem/user/getUser?search=' + userId, {
        method: "GET"
    })
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            // 将用户信息显示在弹窗中的相应输入框中
            console.log(data.data[0])
            document.getElementById("update-username").value = data.data[0].userName;
            document.getElementById("update-phone").value = data.data[0].phone;
            document.getElementById("update-email").value = data.data[0].email;
            if (data.data[0].notice === "feishu") {
                document.getElementById("update-feishu").checked = true;
                document.getElementById("update-feishu-user-id-row").style.display = "flex";
                document.getElementById("update-feiShuUserID").value = data.data[0].feiShuUserID;
            }

        })
        .catch(function(error) {
            console.error("请求错误:", error);
            alert("请求失败 -> "+ error);
        });

    // 打开弹窗
    document.getElementById("update-duty-user-form-overlay").style.display = "flex";

}

// 取消更新用户弹窗
function cancelUpdateDutyUserForm() {
    // 清空输入框的值
    document.getElementById("update-username").value = "";
    document.getElementById("update-feiShuUserID").value = "";
    document.getElementById("update-phone").value = "";
    document.getElementById("update-email").value = "";

    // 关闭弹窗
    document.getElementById("update-duty-user-form-overlay").style.display = "none";
}

// 提交更新用户表单
function submitUpdateForm() {

    var username = document.getElementById("update-username").value;
    var notice = document.querySelector('input[name="update-notice-type"]:checked').value;
    var feishuUserId = document.getElementById("update-feiShuUserID").value;
    var phone = document.getElementById("update-phone").value;
    var email = document.getElementById("update-email").value;

    var data = {
        userId: userId,
        userName: username,
        notice: notice,
        feiShuUserID: feishuUserId,
        phone: phone,
        email: email
    };

    console.log(data)

    // 发送数据给后端的API接口
    fetch("/api/v1/dutySystem/user/update", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(function(response) {
            // 处理响应
            if (response.ok) {
                console.log("数据发送成功");
                // 提取后端返回的userId
                return response.json();
            } else {
                console.error("数据发送失败");
                alert("更新失败");
            }
        })
        .then(function() {
            alert("更新成功");
            cancelUpdateDutyUserForm();
        })
        .catch(function(error) {
            console.error("请求错误:", error);
            alert("请求失败 -> "+ error);
        });


    // 关闭弹窗
    cancelUpdateDutyUserForm();
}

// 删除值班用户
function deleteDutyUser(){

    if (userId.length === 0){
        alert("请选择要删除的用户!");
        return;
    }

    fetch('/api/v1/dutySystem/user/delete?userId=' + userId, {
        method: "POST"
    })
        .then(response => {
            if (response.ok) {
                alert("删除成功！");
                location.reload();
            } else {
                alert("删除失败，请重试！");
            }
        })
        .catch(error => {
            console.error("Error:", error);
            alert("删除失败，请重试！");
        });

}
