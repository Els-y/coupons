# 优惠券分发系统

## 启动

```bash
docker-compose up -d
```

默认开放 20080 端口

## 代码流程

![code flow](assets/code_flow.png)

## TODO

- [x] ~~数据表设计问题有问题，无法区分用户的优惠券是从哪个商家获取的，无法判定是否有重复领取~~ 优惠券名称唯一
- [x] 优惠券查询、用户查询等添加缓存
- [x] 使用消息队列异步减库存
- [ ] 添加 LRU 缓存，未命中再去请求 Redis
- [ ] 优惠券 (username, couponName) 索引
- [ ] 数据库调优
- [ ] Redis 单点故障
- [ ] 性能测试
  - [ ] 后端单实例和多实例差异
  - [ ] 优惠券表有无索引差异
