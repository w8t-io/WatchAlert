## DockerCompose

目录结构
```yaml
[root@master01 w8t]# tree
.
├── config
│   └── config.yaml
└── docker-compose.yaml

1 directory, 2 files
```
配置文件

- [config.yaml](../../config/config.yaml)

Docker-Compose 
> 注意：w8t-web 的 command。
>
> REACT_APP_BACKEND_PORT=9001 yarn start
>
> 参数解析：
>
> - REACT_APP_BACKEND_PORT：有特殊需要需要修改后端端口，需要在这里指定后端端口。

- [docker-compose.yaml](docker-compose.yaml)


启动项目
```shell
# docker-compose -f docker-compose.yaml up -d
# docker-compose -f docker-compose.yaml ps
   Name                  Command               State                 Ports              
----------------------------------------------------------------------------------------
w8t-mysql     docker-entrypoint.sh mysqld      Up      0.0.0.0:3306->3306/tcp, 33060/tcp
w8t-redis     docker-entrypoint.sh redis ...   Up      6379/tcp                         
w8t-service   /app/watchAlert                  Up      0.0.0.0:9002->9001/tcp           
w8t-web       /bin/sh -c REACT_APP_BACKE ...   Up      0.0.0.0:80->3000/tcp      
```

初始化数据
- [README](../sql/README.md)