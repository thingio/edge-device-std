# MQTTS：MQTT 使用 TLS 建立安全连接

本文以 Docker 镜像的形式部署 Mosquitto 为例进行演示。

测试环境：

- ~~ARCH                   : x86_64~~
- ~~OS                     : Ubuntu 20.04.3 LTS~~
- Docker                 : Docker version 20.10.11, build dea9396
- Docker Compose         : docker-compose version 1.25.5, build 8a1c60f6
- Mosquitto Docker Image : eclipse-mosquitto:2.0-openssl

## TLS 证书生成

```shell

mkdir ca && cd ca
# 生成 CA 私钥
# -des3 表示给生成的 RSA 私钥加密，需要设置密码，在使用该私钥时需要输入密码
openssl genrsa -des3 -out ca.key 2048
# 生成 CA 证书
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt

mkdir broker && cd broker
# 生成 broker 端私钥
# 该命令将产生一个不加密的 RSA 私钥，其中参数 2048 表示私钥的长度
# 如果需要为产生的 RSA 私钥加密，需加上 -des3，对私钥文件加密之后，后续使用该密钥的时候都要求输入密码
openssl genrsa -out broker.key 2048
# 生成 broker 端请求文件
# 该命令使用上一步产生的私钥生成一个签发证书所需要的请求文件 broker.csr，使用该文件向 CA 发送请求才会得到 CA 签发的证书
openssl req -new -out broker.csr -key broker.key
# 生成 broker 端证书
# 该命令将使用 broker 私钥、CA 证书、CA 私钥向 CA 请求生成一个证书文件 broker.crt
openssl x509 -req -days 3650 -in broker.csr -CA ../ca.crt -CAkey ../ca.key -CAcreateserial -out broker.crt

cd .. && mkdir client && cd client
# 生成 client 端私钥
openssl genrsa -out client.key 2048
# 生成 client 端请求文件
openssl req -new -out client.csr -key client.key
# 生成 client 端证书
openssl x509 -req -days 3650 -in client.csr -CA ../ca.crt -CAkey ../ca.key -CAcreateserial -out client.crt

cd ..
# 验证证书
openssl x509 -noout -text -in ca.crt
openssl x509 -noout -text -in ./broker/broker.crt
openssl x509 -noout -text -in ./client/client.crt
```

`注：生成 CA 证书时指定的 Common Name 必须和 Broker & client 不一样。此外，当使用双向认证时，Broker & client 的 Common Name 需要统一为 broker 所在服务器域名，这里使用 mosquitto。`

## Mosquitto Broker 开安全

### 创建用户

```shell
# 使用 mosquitto_passwd 命令创建用户，需要在 mosquitto.conf 指定 password_file
mosquitto_passwd -c ./passwd_file admin
```

### Mosquitto 配置

> mosquitto.conf

```text
listener 8883 # TLS 端口

cafile /mosquitto/config/ca/ca.crt
certfile /mosquitto/config/ca/broker.crt
keyfile /mosquitto/config/ca/broker.key
password_file /mosquitto/config/passwd_file
require_certificate true
allow_anonymous false
```

### Mosquitto 启动脚本

> 目录结构

```shell
.
├── conf
│   ├── ca
│   │   ├── broker
│   │   │   ├── broker.crt
│   │   │   ├── broker.key
│   │   │   └── ca.crt
│   ├── mosquitto.conf
│   ├── mosquitto.conf.template
│   └── passwd_file
├── data
│   └── mosquitto.db
├── docker-compose.yml
└── log
    └── mosquitto.log

6 directories, 21 files
```

> docker-compose.yml

```yaml
version: "2"

services:
  mosquitto-openssl:
    image: eclipse-mosquitto:2.0-openssl
    container_name: mosquitto-openssl
    volumes:
      - ./conf/mosquitto.conf:/mosquitto/config/mosquitto.conf
      - ./conf/passwd_file:/mosquitto/config/passwd_file
      - ./conf/ca/broker/:/mosquitto/config/ca/
      - ./data/:/mosquitto/data/
      - ./log/:/mosquitto/log/
    ports:
      - 8883:8883
    privileged: true
    restart: on-failure
```



### 测试

#### 单向认证（实际使用时单向认证即可）

> 客户端不对服务器的证书链和 Host 进行校验。

```shell
mosquitto_sub -h 172.16.251.163 -p 8883 -u admin -P 123456 --tls-version tlsv1.3 --debug --cafile ca.crt --cert client.crt --key client.key -t /greet --insecure

mosquitto_pub -h 172.16.251.163 -p 8883 -u admin -P 123456 --tls-version tlsv1.3 --debug --cafile ca.crt --cert client.crt --key client.key -t /greet -m "hello world" --insecure
```

#### 双向认证

> 客户端须对服务器的证书链和 Host 进行校验，客户端需要设置 DNS：`127.16.251.163 mosquitto`。

```shell
# 设置 DNS
# sudo vim /etc/hosts
# 172.16.251.163 mosquitto
mosquitto_sub -h mosquitto -p 8883 -u admin -P 123456 --tls-version tlsv1.3 --debug --cafile ca.crt --cert client.crt --key client.key -t /greet

mosquitto_pub -h mosquitto -p 8883 -u admin -P 123456 --tls-version tlsv1.3 --debug --cafile ca.crt --cert client.crt --key client.key -t /greet -m "hello world"
```

## 参考

1. [MQTTS: How to use MQTT with TLS?](https://openest.io/en/2020/01/03/mqtts-how-to-use-mqtt-with-tls/)