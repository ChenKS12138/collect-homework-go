## testing

![testing](https://github.com/ChenKS12138/collect-homework-go/workflows/testing/badge.svg)

项目的`HTTP测试`

### How To Run

```sh
make test
```

### Detail

- Admin
  - 不允许避免邮箱重复注册
  - 不允许使用错误的账号密码登陆
  - 不允许频繁请求邀请码
  - 不允许使用错误邀请码注册
- Project
  - 超级管理员可以查看所有的项目
  - 仅超级管理员有能力查看所有项目
  - 仅超级管理员有能力编辑所有项目
  - 所有管理员可以删除项目
  - 仅超级管理员可恢复项目
- Storage
  - 上传限制正确的文件后缀名
  - 上传限制正确的文件名
  - 不允许使用不同的 secret 覆盖相同文件名的文件
  - 未正确登陆不可下载
  - 仅超级管理员有能力下载所有人项目
  - 下载得到的文件是一个合法的 zip 文件
  - 正常查看已上传文件数
  - 仅超级管理员有能力获取所有项目的文件列表
