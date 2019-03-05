# mail-send
本服务用于Open-falcon发送告警邮件

### 运行
1. 配置环境变量: .bashrc
```
export EMAILUSER=mail@example.com    # 发送方邮件用户
export EMAILPASSWORD=123456          # 发送方邮件密码
export EMAILHOST=smtp.exmail.qq.com  # 邮箱服务器地址
export EMAILPORT=25                  # 端口
```

2. 执行环境变量  
```./bashrc```

3. 运行服务  
```./mail-send```

### 告警邮件API：
url: :8011/api/v1/msg/alarm/email/

# 使用

### 测试连通性
```
curl http://127.0.0.1:8011/api/v1/msg/alarm/email/ -d "tos=mail@example.com&subject=testsubject&content=testcontent"
```

### 在Open-Falcon的Alarm组件的配置文件里，配置对应地址即可
```
"api": {
       ...
        "mail": "http://127.0.0.1:8011/api/v1/msg/alarm/email/",
       ...
    },
```
