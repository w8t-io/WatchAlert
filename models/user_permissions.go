package models

type UserPermissions struct {
	Key string `json:"key"`
	API string `json:"api"`
}

/*
[
    {
        "name":"UserRegister",
        "api":"/api/v1/auth/register"
    },
    {
        "name":"UserUpdate",
        "api":"/api/v1/auth/updateUser"
    },
    {
        "name":"UserList",
        "api":"/api/v1/auth/listUser"
    },
    {
        "name":"UserDelete",
        "api":"/api/v1/auth/deleteUser"
    },
    {
        "name":"UserChangePass",
        "api":"/api/v1/auth/changePass"
    }
]
*/
