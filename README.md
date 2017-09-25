# DazeProxy（服务端）

一个用Golang编写的免费、多功能、高性能代理服务端。

DazeProxy属于Daze代理套件。Daze代理套件包括：

1. [DazeProxy](https://github.com/crabkun/DazeProxy)--Daze代理服务端  
2. [DazeClient](https://github.com/crabkun/DazeClient)--Daze代理客户端  
3. [DazeAdmin](https://github.com/crabkun/DazeAdmin)--DazeProxy的数据库简单管理工具  

## DazeProxy能为你提供什么功能？

- TCP、UDP代理转发（IPv4/IPv6）  
- 多用户  
- 数据传输加密  
- 数据传输伪装  
- 支持外部数据库、用户计时、用户分组
- 模块化（加密和伪装均为模块化，方便第三方开发）

## 对于普通用户

支持TCP、UDP代理并加密传输，同时支持伪装，作用不必多言。  
同时支持不加密不伪装，对于追求低延迟的游戏很有帮助，使自己能搭建游戏加速器成为可能。

## 对于服务器管理员

对于个人服务器管理员，DazeProxy提供了多用户功能，并可以搭配DazeAdmin进行简单的用户管理，让你可以与朋友进行分享。  
对于多服务器管理员，DazeProxy提供了外部数据库和用户计时支持，可以设置用户的过期时间，可以设置服务器所属组，用户可以属于一个/多个组并只能连接所属组的服务器。

## 对于开发者

加密和伪装方式均为模块化设计，并统一和公开了相关接口。第三方如果有更好的想法，可以按照公开的接口进行开发加密方式或者伪装方式。

## 加密和伪装

目前Daze代理套件自带的伪装方式有
- none：不伪装
- http：可伪装成HTTP GET或POST连接  
- tls_handshake：可伪装成TLS1.2连接  

目前Daze代理套件自带的加密方式有
- none：不加密
- keypair-rsa：服务端生成RSA密钥并发送公钥与客户端协商aes密钥，然后进行aes128位cfb模式加密  
- psk-aes-128-cfb：客户端与服务端利用约定好的预共享密钥进行aes128位cfb模式加密  
- psk-aes-256-cfb：客户端与服务端利用约定好的预共享密钥进行aes256位cfb模式加密  
- psk-rc4-md5：客户端与服务端利用约定好的预共享密钥进行rc4加密  

## 哪里下载？
由于某些不可描述的原因，暂停下载一段时间

## 相关教程（持续更新中）
[服务端配置文件详解](https://github.com/crabkun/DazeProxy/wiki/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%A6%E8%A7%A3)   
[快速架设DazeProxy服务器](https://github.com/crabkun/DazeProxy/wiki/%E5%BF%AB%E9%80%9F%E6%9E%B6%E8%AE%BEDazeProxy%E6%9C%8D%E5%8A%A1%E5%99%A8)  
[开启多用户验证并用DazeAdmin管理用户](https://github.com/crabkun/DazeProxy/wiki/%E5%BC%80%E5%90%AF%E5%A4%9A%E7%94%A8%E6%88%B7%E9%AA%8C%E8%AF%81%E5%B9%B6%E7%94%A8DazeAdmin%E7%AE%A1%E7%90%86%E7%94%A8%E6%88%B7)  
[各加密方式的详细解释与区别](https://github.com/crabkun/DazeProxy/wiki/%E5%90%84%E5%8A%A0%E5%AF%86%E6%96%B9%E5%BC%8F%E7%9A%84%E8%AF%A6%E7%BB%86%E8%A7%A3%E9%87%8A%E4%B8%8E%E5%8C%BA%E5%88%AB)  
[各伪装方式的详细解释与区别](https://github.com/crabkun/DazeProxy/wiki/%E5%90%84%E4%BC%AA%E8%A3%85%E6%96%B9%E5%BC%8F%E7%9A%84%E8%AF%A6%E7%BB%86%E8%A7%A3%E9%87%8A%E4%B8%8E%E5%8C%BA%E5%88%AB)    
[连接外部数据库的详细说明 ](https://github.com/crabkun/DazeProxy/wiki/%E8%BF%9E%E6%8E%A5%E5%A4%96%E9%83%A8%E6%95%B0%E6%8D%AE%E5%BA%93%E7%9A%84%E8%AF%A6%E7%BB%86%E8%AF%B4%E6%98%8E)  
[加密与伪装的开发文档 ](https://github.com/crabkun/DazeProxy/wiki/%E5%8A%A0%E5%AF%86%E4%B8%8E%E4%BC%AA%E8%A3%85%E7%9A%84%E5%BC%80%E5%8F%91%E6%96%87%E6%A1%A3)  
[各种常见的问题与答案](https://github.com/crabkun/DazeProxy/wiki/%E5%90%84%E7%A7%8D%E5%B8%B8%E8%A7%81%E7%9A%84%E9%97%AE%E9%A2%98%E4%B8%8E%E7%AD%94%E6%A1%88)

## 感谢（Thanks）
本项目借助了以下开源项目的力量才能完成，非常感谢以下项目以及其作者们！  
- Xorm：[https://github.com/go-xorm/xorm](https://github.com/go-xorm/xorm)  
- Go-MySQL-Driver：[https://github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)  
- go-sqlite3：[https://github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)  

## 开源协议
BSD 3-Clause License

## 声明
本软件仅供技术交流和游戏网络延迟加速，并非侵入或非法控制计算机信息系统的软件，严禁将本软件用于商业及非法用途，如软件使用者不能遵守此规定，请马上停止使用并删除，对于因用户使用本软件而造成任何不良后果，均由用户自行承担，软件作者不负任何责任。您下载或者使用本软件，就代表您已经接受此声明，如产生法律纠纷与本人无关。