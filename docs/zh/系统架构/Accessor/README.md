# Accessor

1. Accessor 功能：
    1. 设备数据持久化与查询；
        1. 流程：
            1. 设备元数据增加 `recording` 字段，表示是否将采集的设备进行落库；
            2. accessor 启动后，连接到 manager，先从 manager 获取全量设备元数据，再监听设备元数据的变更（使用 edge-device-std 定义的 operations ）；
            3. 如果设备的 `recording` 字段为 true，则在 accessor 中启动一个 Recorder 监听 event 以及 props 对应的主题，将采集到的数据落库，否则跳过或者关闭
               Recorder；
        2. 是否需要支持多数据源，如 InfluxDB | TDEngine 等？
        3. 是否可使用 ORM 框架减少开发量，如 gorm | sqlx | ent 等？
        4. 是否集成数据可视化工具，如 Grafana | Superset | Metabase 等？
    2. 是否需要支持数据推送：
        1. 前端通过 WebSocket 获取设备推送的数据？
        2. 后端直接通过 MessageBus 订阅设备对应的主题？
2. 对于 [当前的架构设计](https://app.diagrams.net/#G1E_OcUtDI-vPk-1XZFqoLdzKV46zRjKUY) 来说，Manager 负责操作元数据，Accessor 负责操作设备数据。是否可以让
   Manager 集成 Accessor，即 Manager 同时负责操作元数据和设备数据？
3. `props *` 的设计跟 `event` 是否有重叠？或者说是否可以用 `event` 表示 `props *`？