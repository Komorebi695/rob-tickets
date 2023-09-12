# 一个简单的抢票系统
## 需求
设计一个算法进行抢票的调度（注意ACID！），并且让每个人都能够在**可接受的时间内**获得自己抢票是否成功的结果。
当然，注意每张票都是独一无二的，都有自己的编码，不能超售，也不能一票多售，更不能在有票的情况下就不售了！
在算法演示时，需要自行设计模拟算法，模拟C/S模型，并且模拟一定人数在一分钟内陆续或同时**访问服务器**，完成抢票操作！

## 数据库设计
> users

编号、用户编号
> tickets

编号、票编号、价格、剩余票数量
> ticket_orders

编号、订单编号、票编号、用户编号、数量、状态
## **版本1.0**
### 逻辑设计
![image.png](https://cdn.nlark.com/yuque/0/2023/png/23153452/1678609065329-1a94db6b-efff-4de1-a6c6-573765f146d3.png#averageHue=%23f8f8f8&clientId=ufece5b38-d46d-4&from=paste&height=322&id=u494855b7&originHeight=322&originWidth=1130&originalType=binary&ratio=1&rotation=0&showTitle=false&size=55567&status=done&style=none&taskId=ubba2e0b4-4900-4f46-a47a-2861a6e5ef8&title=&width=1130)

### 测试数据
| 票数量 | 并发数量 | 耗时 | 失败次数 | 请求失败率 | 备注 |
| --- | --- | --- | --- | --- | --- |
| 10 | 100 | 106.2861ms | 0 | 0 |  |
| 10 | 500 | 369.7883ms | 0 |
|  |
| 10 | 1000 | 766.8147ms | 0 |
|  |
| 10 | 5000 | 3.7108553s | 0 |
|  |
| 10 | 10000 | 7.253048s | 0 |
|  |
| 10 | 15000 | 11.2189516s | 0 |  |  |
| 10 | 20000 | 15.7787627s | 0 |  |  |
| 10 | 25000 | 18.8928289s | 0 |  |  |
| 10 | 30000 |  22.0385819s | 0 |  |  |
| 10 | 35000 | 25.8593406s | 0 |  |  |
| 10 | 40000 | 29.4214239s | 0 |  |  |
| 1000 | 40000 | 32.7520569s | 0 |  |  |
| 5000 | 40000 | 33.8555872s | 0 |  |  |
| 10000 | 40000 | 34.9738608s | 0 |  |  |
| 20000 | 40000 | 38.2229415s | 0 |  |  |
| 30000 | 40000 | 41.2492063s | 0 |  |  |
| 40000 | 40000 | 43.790149s | 0 |  |  |
| 10 | 45000 | 33.2403034s | 108 | 0.240000% | Client.Timeout |
| 10 | 45000 | 33.6592078s | 76 | 0.168888% | Client.Timeout |
| 10 | 45000 | 31.094261s | 0 |  |  |
| 10 | 45000 | 50.9009834s | 804 | 1.786666% | Client.Timeout  && Error |
| 10 | 45000 | 33.9598551s | 166 | 0.368888% | Client.Timeout  && Error |
| 10 | 44000 | 32.9489333s | 171 | 0.388636% | Client.Timeout |
| 10 | 44000 | 34.2827626s | 533 | 1.211364% | Client.Timeout  && Error |
| 10 | 44000 | 31.959539s | 269 | 0.611364% | Client.Timeout |
| 10 | 44000 | 51.3028091s | 141 | 0.141000% | Client.Timeout  && Error |
| 10 | 44000 | 32.2648524s | 273 | 0.320455% | Client.Timeout |
| 10 | 44000 | 35.8412562s | 360 | 0.818182% | Client.Timeout  && Error |
| 10 | 44000 | 31.782626s | 139 | 0.315909% | Client.Timeout  |
| 10 | 44000 | 32.7380922s | 159 | 0.361363% | Client.Timeout  && Error |
| 10 | 44000 | 38.6249777s | 665 | 1.511364% | Client.Timeout  && Error |
| 10 | 44000 | 31.904696s | 111 | 0.252273% | Client.Timeout |
| 10 | 44000 | 中断 |  |  | Client.Timeout && SLOW SQL >= 1s |
| 10 | 44000 | 33.190857s | 46 | 0.104545% | Client.Timeout |
| 10 | 44000 | 32.4256817s | 91 | 0.206818% | Client.Timeout |
| 10 | 44000 | 36.5579036s | 767 | 1.743182% | Client.Timeout &&  SLOW SQL >= 1s && |
| 10 | 44000 | 33.3162491s | 109 | 0.247727% | Client.Timeout |
| 10 | 40000 |  33.8861358s | 0 |  |  |
| 10 | 40000 | 31.4957282s | 0 |  |  |
| 10 | 40000 | 31.0184887s | 0 |  |  |
| 100 | 45000 |  33.8706609s | 99 | 0.220000% | Client.Timeout |
| 100 | 45000 | 33.0276449s | 128 | 0.284444% | Client.Timeout |
| 100 | 45000 | 33.2984907s | 224 | 0.497778% | Client.Timeout && Error |
| 100 | 45000 | 33.5713699s | 95 | 0.211111% | Client.Timeout && Error |
| 100 | 45000 | 33.5108618s | 0 |  |  |
| 100 | 45000 | 33.641454s | 198 | 0.440000% | Client.Timeout |
| 100 | 45000 | 33.1151534s | 263 | 0.584444% | Client.Timeout |
| 100 | 45000 |  30.1010509s | 0 | 关掉了gin的日志 |
| 100 | 45000 | 33.4993167s | 247 | 0.548889% | Client.Timeout && Error |
| 100 | 45000 | 34.3163744s | 389 | 0.864444% | Client.Timeout && Error |
| 100 | 45000 | 33.235657s | 296 | 0657778% | Client.Timeout |
| 100 | 50000 | 40.5437066s | 110 | 0.220000% | Client.Timeout &&  Error  |
| 100 | 50000 | 中断 |  |  |  |
| 100 | 50000 | 45.5639837s | 163 | 0.326000% | Client.Timeout &&  SLOW SQL && |
| 100 | 50000 | 53.1829918s | 492 | 0.984000% | Client.Timeout &&  SLOW SQL && |
| 100 | 50000 | 中断 |  |  | Client.Timeout &&  SLOW SQL && |
| 100 | 50000 | 47.5607266s | 539 | 1.078000% | Client.Timeout &&  SLOW SQL && |
| 100 | 50000 | 53.3928181s | 277 | 0.454000% | Client.Timeout &&  SLOW SQL && |
| 100 | 50000 | 59.5507965s | 10796 | 21.592000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 50000 | 1m15.6294211s | 415 | 0.830000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 50000 | 1m4.7884235s | 234 | 0.468000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 55000 | 1m2.9695642s | 2298 | 4.178182% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 55000 |  56.8557979s |  2718 | 4.941818% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 55000 | 1m0.6392447s | 11490 | 20.890909% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 55000 | 1m0.5575765s | 81 | 0.147273% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 55000 | 55.2305289s | 1166 | 2.120000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 60000 | 58.9001051s | 6003 | 10.005000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 60000 | 1m6.3653592s | 32238 | 53.730000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 100 | 60000 | 59.5837177s | 7890 | 13.150000% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 10000 | 60000 |  56.1324834s | 25043 | 41.738335% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 20000 | 60000 | 1m1.4605706s | 11603 | 19.338333% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 30000 | 60000 | 59.5974752s | 28289 | 47.148335% | Client.Timeout &&  SLOW SQL && 明显停顿 |
| 30000 | 60000 | 59.8599832s | 23144 | 38.573334% | 票有剩余 |


### 出现的问题

1. **Client.Timeout**

context deadline exceeded (Client.Timeout exceeded while awaiting headers)

2. **dial tcp 127.0.0.1:3306: connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.**

3. **SLOW SQL >= 1s**

[1586.294ms] [rows:0] update `tickets` set `ticket_count`=`ticket_count` - 1 where `ticket_id` = 'h1k4J7Dyt0' and ticket_count >= 1;
### 总结
#### 个人思考

1. 当45000并发数时，开始几次也要请求全部成功的情况，会不会跟order表内的数据太多有关？

测试发现：在order表的数据量为1522683时也全部成功了一次。并且在后面的并发量和order表的数据继续增加，并没有出现 随着order表数据增加，而成功次数减少。
结论: 无关。也可能和order表的数据量有关，但是到目前为止，没有出现明显的差异。

2. 被抢的票，数量越多，并发数相同时，请求的时间就相对较久一点吗？

测试发现: 当并发为40000时，随着被抢票数量的增加，总的请求消耗的时间在增加。
结论：是的。

3. 影响最大请求数量的因素可能和服务器的环境配置？

未测试。
猜测：有一定关系。

4. 测试数据是通过自己写的Client进行测试，如果使用压测工具是否更贴近实际？

#### 观察总结

- 并发维持在 40000 左右时，是稳定的，时间大约在30秒左右，没有出现请求超时或则其他错误。同时随着票数的增加，相应的请求消耗的时间也有一定的增加。直观点说，4万人 抢 **一张票**消耗的时间比 4万人 抢 **四万张票** 的时间要少。

- 并发维持在 45000 左右时，开始有请求超时，服务端出现如下错误; 如果把gin框架的日志关掉，可以满足需求。
Errordial tcp 127.0.0.1:3306: connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.

- 并发维持在 50000 左右时，请求超时，服务端出现如上错误和SLOW SQL，有明显停顿。

- 随着并发数量继续增加，出现请求超时，服务端出现错误次数越来越多。
- 并发大于 60000时，当票数量较多时，会出现票有剩余的情况。

## 版本2.0
> 与版本1.0 不同之处：使用Redis数据库+Kafka。

### 逻辑设计
![image.png](https://cdn.nlark.com/yuque/0/2023/png/23153452/1678609083925-a20c6536-7c27-4e0b-b2ce-68ed492f2a47.png#averageHue=%23fafafa&clientId=ufece5b38-d46d-4&from=paste&height=547&id=u2905427b&originHeight=547&originWidth=1146&originalType=binary&ratio=1&rotation=0&showTitle=false&size=64835&status=done&style=none&taskId=ub5f49a00-30b6-4658-8fb0-b05c109e4e8&title=&width=1146)

### 测试数据
| 票数量 | 并发数量 | 耗时 | 失败次数 | 请求失败率 | 备注 |
| --- | --- | --- | --- | --- | --- |
| 10 | 100 | 67.6127ms | 0 |  |  |
| 10 | 500 | 602.3949ms | 0 |  |  |
| 10 | 1000 | 1.5956159s | 0 |  |  |
| 10 | 5000 | 5.8874455s | 0 |  |  |
| 10 | 10000 | 13.2766459s | 0 |  |  |
| 10 | 15000 | 21.1679174s | 0 |  |  |
| 10 | 20000 | 27.5602344s | 0 |  |  |
| 10 | 30000 | 45.4450886s | 0 |  |  |
| 10 | 40000 | 1m0.1203445s | 0 |  |  |
| 10 | 45000 | 53.8913521s | 0 |  |  |
| 10 | 50000 | 59.9181524s | 0 |  |  |
| 100 | 50000 | 1m0.3864073s | 0 |  |  |
| 500 | 50000 | 1m37.5032234s | 0 |  |  |
| 500 | 60000 | 1m21.6188633s | 0 |  |  |
| 500 | 80000 | 1m0.236706s | 0 |  |  |
| 1000 | 100000 | 45.8812738s | 0 |  |  |
| 1000 | 120000 | 1m3.2379635s | 0 |  |  |
| 1400 | 140000 | 1m8.1631063s | 0 |  |  |
| 1600 | 160000 | 57.1305129s | 0 |  |  |
| 1800 | 180000 | 1m2.9356271s | 0 |  |  |
| 1700 | 170000 | 1m23.781692s | 0 |  |  |
| 1750 | 175000 | 1m14.8607791s | 0 |  |  |
| 1790 | 179000 | 56.7539284s | 0 |  |  |
| 1790 | 179000 | 1m0.5236141s | 0 |  |  |
| 1790 | 179000 | 52.3496694s | 0 |  |  |
| 1800 | 180000 | 1m0.1085117s | 0 |  |  |
| 1800 | 180000 | 1m0.1953023s | 0 |  |  |
| 1800 | 180000 | 58.422238s | 0 |  |  |
| 1900 | 190000 | 1m12.4950019s | 101 | 0.053158% | Client.Timeout |
| 2000 | 200000 | 53.4128369s | 507 | 0.253500% | Client.Timeout |
| 1900 | 190000 | 58.3609713s | 179 | 0.094211% |  |
| 1950 | 195000 | 1m20.6553042s | 250 | 0.128205% |  |
| 1900 | 190000 | 1m14.4173575s | 4 | 0.002105% |  |
| 1900 | 190000 | 1m29.0673797s | 149 |  |  |
| 1900 | 190000 | 1m16.5931292s | 0 |  |  |

### 总结
Redis可以承受大量的并发请求，Redis的并发处理能力远大于MySQL，使用Redis作为MySQL的缓存可以很大程度上提高性能。


