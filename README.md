# newvelas-adapter

newvelas-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf目录，新建NVLX.ini文件，编辑如下内容：

```ini

# node api url
serverAPI = "http://127.0.0.1:1005"
# fix gas limit
fixGasLimit = ""
# fix gas price
fixGasPrice = ""
# Cache data file directory, default = "", current directory: ./data
dataDir = ""

```

## 资料介绍

### 官网

https://www.velas.com/

### 区块浏览器

https://mainnet-v2.velas.com/
