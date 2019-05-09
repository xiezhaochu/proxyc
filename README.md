# Proxyc
Proxyc是一种终端命令行代理工具，在需要使用代理的命令前，加上proxyc即可。

如使用 proxyc go get XXX，可以get到一些之前因为某些原因无法get到的库。

可以在用户主目录下配置名为.proxyc的json文件，格式如下，来设置代理信息。

```json
{
    "https_proxy":"127.0.0.1:1111",
    "http_proxy":"127.0.0.1:1111"
}
```



如果没有相关文件，则默认代理地址都是127.0.0.1:1080 。

目前测试在Windows和Linux环境下都可以正常使用。