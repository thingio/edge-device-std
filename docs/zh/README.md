# EDS

**EDS**（Edge Device Std）是为 IOT 场景设计并开发的设备接入层 SDK，提供快速将多种协议的设备接入平台的能力。

## 特性

1. 轻量：只使用 MQTT 作为设备驱动服务与设备管理服务之间数据交换的中间件，无需引入多余组件；
2. 通用：
    1. 将 MQTT 封装为 MessageBus，向上层模块提供支持；
    2. 基于 MessageBus 定义和实现了通用的**元数据操作**与**物模型操作**规范。

## 术语

### 物模型

[物模型](https://blog.csdn.net/zjccoder/article/details/107050046)，作为一类设备的抽象，描述了设备的：

- 属性（property）：用于描述设备状态，支持读取和写入；
- 方法（method）：设备可被外部调用的能力或方法，可设置输入和输出参数，相比与属性，方法可通过一条指令实现更复杂的业务逻辑；
- 事件（event）：用于描述设备主动上报的事件，可包含多个输出参数，参数必须是某个“属性”。

产品（Product）：即物模型。

设备（Device）：与真实设备一一对应，必须属于并且仅属于某一个产品。

### 元数据

元数据包括协议（Protocol）、产品（Product）及设备（Device）：

- 协议对应了一个特定协议的描述信息；
- 产品对应了一个特定产品的描述信息；
- 设备对应了一个特定设备的描述信息。

### 设备驱动服务

负责一种特定设备协议的接入实现，主要功能包括：

- 协议注册：将当前设备协议注册到设备管理服务；
- 设备初始化：从设备管理服务获取产品 & 设备元数据，加载设备驱动；
- 元数据监听：监听设备管理服务产品 & 设备元数据的变更，加载/卸载/重加载驱动；
- 设备属性读写：从真实设备或驱动缓存中读取产品定义的属性数据；
- 设备方法调用：调用设备支持的方法；
- 设备事件监听：将设备主动推送的数据转发到 MessageBus 中；
- 设备状态检查：检查真实设备的健康状态，并周期性推送到 MessageBus 中；
- 设备驱动服务状态检查：检查设备驱动服务的健康状态，并周期性推送到 MessageBus 中。

### 设备管理服务

负责管理接入的设备驱动服务，主要功能包括：

- 协议管理：接收设备服务的注册请求；
- 产品管理：基于特定协议定义产品 & 产品 CRUD；
- 设备管理：基于特定产品定义设备 & 设备 CRUD。

## Topic 约定

因为 EDS 基于 MQTT 实现数据交换，所以我们基于 MQTT 的 Topic & Payload 概念定义了我们自己的数据格式及通信规范。

### 物模型操作

对于物模型来说，Topic 的一般格式为 `DATA/${Version}/${OptMode}/${ProtocolID}/${ProductID}/${DeviceID}/${FuncID}/${OptType}[/${ReqID}]`
：

- `Version`：DATA Topic 的版本（与 META Topic 的版本相互独立），在 `edge-device-std/version` 中定义；
- `OptMode`，操作模式，可选 `DOWN | UP | UP-ERR`：
    - `DOWN`：对应设备管理服务（北向）向设备驱动服务（南向）发起的操作；
    - `UP` 对应设备驱动服务（南向）向设备管理服务（北向）反馈的正常数据；
    - `UP-ERR` 对应设备驱动服务（南向）向设备管理服务（北向）反馈的错误信息；
- `ProtocolID`：协议元数据的 UUID，协议（设备驱动服务）唯一；
- `ProductID`：产品元数据的 UUID，产品唯一；
- `DeviceID`：设备元数据的 UUID，设备唯一；
- `FuncID`：当前操作的 ID，对应 `property | method | event` 的 ID；
- `OptType`，操作类型：
    - 双向操作：
        - 对于 `property`，可选 `READ | HARD-READ | WRITE`，分别对应于设备属性的软读、硬读和写入操作；
        - 对于 `method`，可选 `CALL`，对应设备方法的调用操作；
    - 单向操作：
        - 对于 `property`，可选 `PROP`，对应设备属性的上报操作；
        - 对于 `event`，可选 `EVENT`，对应设备事件的上报操作（与 `PROP` 的区别在于，`PROP` 由设备驱动主动向真实设备获取，而 `EVENT` 由真实设备主动向设备驱动推送）；
        - 特别地，对于设备状态检测功能而言，定义了特定的操作类型 `STATUS`，用于设备状态的上报操作；
- `ReqID`：对于双向操作而言，如果不能绑定属于该操作的请求与响应，则当多个操作并发执行时，会导致请求与响应的混乱， 因此需要一个唯一的 UUID 来表示标识同一组操作。

#### 示例

1. 从设备读取属性（软读）：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/    +/+/+/                            READ/+
PUB: DATA/${VERSION}/  DOWN/${prot_id}/  ${prod_id}/${dev_id}/${prop_id}/  READ/${req_id}

# 设备驱动服务（南向）
SUB: DATA/${VERSION}/  DOWN/${prot_id}/  +/+/+/                            READ/+
PUB: DATA/${VERSION}/  UP/${prot_id}/    ${prod_id}/${dev_id}/${prop_id}/  READ/${req_id}  --payload ${props}
```

2. 从设备读取属性（硬读）：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/    +/+/+/                            HARD-READ/+
PUB: DATA/${VERSION}/  DOWN/${prot_id}/  ${prod_id}/${dev_id}/${prop_id}/  HARD-READ/${req_id}

# 设备驱动服务（南向）
SUB: DATA/${VERSION}/  DOWN/${prot_id}/  +/+/+/                            HARD-READ/+
PUB: DATA/${VERSION}/  UP/${prot_id}/    ${prod_id}/${dev_id}/${prop_id}/  HARD-READ/${req_id}  --payload ${props}
```

3. 往设备写入属性：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/    +/+/+/                            WRITE/+
PUB: DATA/${VERSION}/  DOWN/${prot_id}/  ${prod_id}/${dev_id}/${prop_id}/  WRITE/${req_id}  --payload ${props}

# 设备驱动服务（南向）
SUB: DATA/${VERSION}/  DOWN/${prot_id}/  +/+/+/                            WRITE/+
PUB: DATA/${VERSION}/  UP/${prot_id}/    ${prod_id}/${dev_id}/${prop_id}/  WRITE/${req_id}  --payload ${result}
```

4. 调用设备方法：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/    +/+/+/                              CALL/+
PUB: DATA/${VERSION}/  DOWN/${prot_id}/  ${prod_id}/${dev_id}/${method_id}/  CALL/${req_id}  --payload ${ins}

# 设备驱动服务（南向）
SUB: DATA/${VERSION}/  DOWN/${prot_id}/  +/+/+/                              CALL/+
PUB: DATA/${VERSION}/  UP/${prot_id}/    ${prod_id}/${dev_id}/${method_id}/  CALL/${req_id}  --payload ${outs}
```

5. 设备属性推送：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/   +/+/+/                            PROP/+

# 设备驱动服务（南向）
PUB: DATA/${VERSION}/  UP/${prot_id}/   ${prod_id}/${dev_id}/${prop_id}/  PROP/*  --payload ${props}
```

6. 设备事件推送：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/   +/+/+/                             EVENT/+

# 设备驱动服务（南向）
PUB: DATA/${VERSION}/  UP/${prot_id}/   ${prod_id}/${dev_id}/${event_id}/  EVENT  --payload ${props}
```

7. 设备状态推送：

```text
# 设备管理服务（北向）
SUB: DATA/${VERSION}/  UP/${prot_id}/   +/+/+/                            STATUS/+

# 设备驱动服务（南向）
PUB: DATA/${VERSION}/  UP/${prot_id}/   ${prod_id}/${dev_id}/-/           STATUS/*  --payload ${driver_status}
```

### 元数据操作

对于元数据操作来说，Topic 的一般格式为 `META/${Version}/${OptMode}/${ProtocolID}/${OptType}[/${ID}]`：

- `Version`：META Topic 的版本（与 DATA Topic 的版本相互独立），在 `edge-device-std/version` 中定义；
- `OptMode`，操作模式，可选 `DOWN | UP`：
    - `DOWN`：对应设备管理服务（北向）向设备驱动服务（南向）发起的操作；
    - `UP` 对应设备驱动服务（南向）向设备管理服务（北向）反馈的正常数据；
- `OptType`，操作类型，可选 `INIT | PRODUCT | DEVICE`：
    - 单向操作：
        - `PRODUCT`，对应于产品的增/删/改操作，由设备驱动服务根据缓存判断具体操作类型：
            - 首先判断 Payload 是否为空，如果为空则表示产品删除；
            - 否则判断设备驱动服务缓存是否存在 `${prod_id}`，如果不存在则表示产品创建；
            - 否则表示产品更新；
        - `DEVICE`，对应于设备的增/删/改操作，由设备驱动服务根据缓存判断具体操作类型：
            - 首先判断 Payload 是否为空，如果为空则表示设备删除；
            - 否则判断设备驱动服务缓存是否存在 `${prod_id}`，如果不存在则表示设备创建；
            - 否则表示设备更新；
        - `STATUS`，对应于设备驱动服务的状态上报操作，需要包含状态码、协议元数据、上次 `APPEND` 操作的时间戳（初始为 0）；
        - `APPEND`，对应于设备管理服务向设备驱动服务发起的产品 & 设备元数据追加操作，设备管理服务会根据 `STATUS` 中的时间戳向设备驱动服务增量发送更新的产品 & 设备；
- `ID`：
    - 对于单向操作而言，用于指定产品/设备的增/删/改的操作对象。

#### 示例

1. 产品元数据更新：

```text
# 设备管理服务（北向）
PUB: META/${VERSION}/  DOWN/${prot_id}/  PRODUCT/${prod_id}  --payload ${product}

# 设备驱动服务（南向）
SUB: META/${VERSION}/  DOWN/${prot_id}/  PRODUCT/+
```

2. 设备元数据更新：

```text
# 设备管理服务（北向）
PUB: META/${VERSION}/  DOWN/${prot_id}/  DEVICE/${dev_id}  --payload ${device}

# 设备驱动服务（南向）
SUB: META/${VERSION}/  DOWN/${prot_id}/  DEVICE/+
```

3. 设备驱动服务状态推送：

```text
# 设备管理服务（北向）
SUB: META/${VERSION}/  UP/${prot_id}/  STATUS/+

# 设备驱动服务（南向）
PUB: META/${VERSION}/  UP/${prot_id}/  STATUS/*   --payload ${driver_status}
```

4. 驱动服务元数据初始化：

```text
# 设备管理服务（北向）
PUB: META/${VERSION}/  DOWN/${prot_id}/  INIT/*  --payload ${products & devices}

# 设备驱动服务（南向）
SUB: META/${VERSION}/  DOWN/${prot_id}/  INIT/+
```