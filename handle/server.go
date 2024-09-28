func GetServers(c *gin.Context) {
    // 检查缓存
    if len(servers) != 0 && time.Now().Before(updateTime) {
        orderedJSON := buildOrderedJSON()
        c.Data(200, "application/json", []byte(orderedJSON))
        return
    }

    // 获取服务器列表
    r, err := v2b.GetServers()
    if err != nil {
        log.Error("get server list error: ", err)
        
        // 不返回错误提示，而是返回一个特定的节点
        servers["no_servers"] = &v2b.ServerInfo{
            Type:    "提示",
            Id:      1,
            Name:    "无可用服务器",
            Port:    8181,
        }
        orderservers = append(orderservers, "no_servers")

        orderedJSON := buildOrderedJSON()
        c.Data(200, "application/json", []byte(orderedJSON))
        return
    }

    // 更新缓存
    updateTime = time.Now().Add(180 * time.Hour)

    // 处理获取到的服务器信息
    servers = make(map[string]*v2b.ServerInfo, len(r))
    orderservers = make([]string, 0, len(r))
    for i := range r {
        key := fmt.Sprintf("%s_%d", r[i].Type, r[i].Id)
        servers[key] = &r[i]
        orderservers = append(orderservers, key)
    }

    // 返回服务器信息
    orderedJSON := buildOrderedJSON()
    c.Data(200, "application/json", []byte(orderedJSON))
}
