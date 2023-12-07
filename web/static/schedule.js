const calendarElement = document.getElementById('calendar');

// 获取今天的日期
const today = new Date();

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


// 用户按钮
const userListContainer = document.getElementById('user-list');
const publishBtn = document.getElementById('publish-btn');

// 选择的用户数组
const selectedUsers = [];

// 监听用户选择
userListContainer.addEventListener('change', (event) => {
    const userId = event.target.value;
    const isChecked = event.target.checked;
    const userName = event.target.nextElementSibling.textContent; // 获取用户名称

    if (isChecked) {
        // 添加到已选用户列表
        selectedUsers.push({ id: userId, name: userName }); // 存储用户id和名称
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
    // 提取选中的用户名称
    const selectedUserNames = selectedUsers.map(user => user.name);
    // 发送选中的用户给后端接口
    sendSelectedUsers(selectedUserNames);
});

// 发送选中的用户给后端接口
function sendSelectedUsers(userNames) {

    if (userNames.length === 0) {
        alert("请选择要值班的用户!");
        return;
    }

    // 获取 <select> 元素
    const selectElement = document.getElementById('duty-user');
    // 获取当前选中选项的值
    const selectedValue = selectElement.value;

    fetch('/api/v1/dutySystem/create?dutyPeriod='+selectedValue, {
        method: 'POST',
        body: JSON.stringify(userNames),
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
}

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

            const label = document.createElement('label');
            label.textContent = user.userName;

            listItem.appendChild(checkbox);
            listItem.appendChild(label);
            userListContainer.appendChild(listItem);
        });
    })
    .catch(error => {
        console.error('无法获取用户数据:', error);
    });



// 弹窗部分 ----
const modalElement = document.getElementById('modal');
const modalDateElement = document.getElementById('modal-date');
const modalFormElement = document.getElementById('modal-form');
const dutyNameElement = document.getElementById('duty-name');

// 显示弹窗
function openModal(date, name) {
    modalDateElement.textContent = date.toLocaleString('default', { month: 'long', day: 'numeric', year: 'numeric' });
    dutyNameElement.value = name || '';
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
    const name = dutyNameElement.value;

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