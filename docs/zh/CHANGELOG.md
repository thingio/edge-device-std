# CHANGELOG

## 2021.12

1. 确定[系统基本架构](./系统架构/README.md)
2. [设备驱动示例（生成随机数）](https://github.com/thingio/edge-randnum-driver)
3. 新增特性：
    - 支持 command-line flag 指定配置文件路径，如 `xxx -cp etc -cn config` 表示指定配置文件为 `./etc/config.yaml`
    - 定义 MessageBus 接口封装 MQ 操作逻辑，并提供了 MQTTMessageBus 作为默认支持
    - 定义 MetaStore 接口封装元数据操作逻辑，并提供了 FileMetaStore 作为默认支持
    - 区分单向和双向操作，对于双向操作（基于 MessageBus 的 Call 方法实现）而言，使用全局唯一的 UUID - `ReqID` 标识同一组操作
    - 区分属性软读（soft-read）和硬读（hard-read），前者会从设备影子（DeviceTwin）中获取该属性的缓存值，后者会直接从真实设备中读取该属性的值
    - 支持 MQTTS 配置，可参考 [MQTT 使用 TLS 建立安全连接](./系统安全/MQTTS:%20MQTT%20使用%20TLS%20建立安全连接.md)
    - 将设备数据封装为类似于 DeviceData 的结构体，包含属性名、数据类型、数据值、数据采集时间戳等信息

## TODO

1. 支持设备级别的 MQTT(S) QoS 配置
2. 为 `edge-device-manager` & `edge-device-accessor` 提供 HTTP(S) & WS(S) 访问接口及 Client，与业务无关的代码封装在 `edge-device-std`
   中提供基础服务
3. `edge-device-accessor` 确定业务边界、架构设计与代码实现
4. 调整 Protocol、Product 与 Device 元数据结构的定义，如 DeviceStatus 是动态消息，只能由系统更改，不应该放在 Device 这样的元数据中
5. 调研 TCP & MQTT(mosquitto) 如何切分数据包，主要是要确定当传输一个大数据文件（如图片）时，MQTT 如何将其拆分，TCP 如何将其拆分？