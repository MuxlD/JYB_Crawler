# JYB_Crawler
### Basics
+ basicsData.go  
  > MySql: 配置数据库相关参数  
  > TrainingSchool: 培训机构对象  
  > Type: 培训机构类型  
  > TsUrl: 爬虫中间对象
+ initMysql.go
  > 连接数据库`MysqlInit`，创建数据库表格`CreateTable`及数据库连接接口`GetDB`  
### eduDate
+ allLink.go
  >`getAllEdu`获取机构的连接，通过`TsUrl`并放入通道。
+ allLink_db.go
  >`AllLink`检查数据库中是否储存了`TsUrl`，避免重复爬取。
+ chrome.go
  >爬虫前的准备。
+ chrome_test.go
  >测试`chrome.go`部分函数。
+ dataCollation.go
  >建立爬虫主干
+ findType.go
  >获取所有的机构类型链接，并存入数据库
+ pendLink.go
  >机构信息页面模板`2`爬取程序
+ pendLink_test.go
  >测试可用性
+ sortQuery.go
  >机构信息页面模板`1`爬取程序
+ sortQuery_test.go
  >测试可用性
### elasticsearch
+ init.go
  > 创建`Client`、elasticsearch索引
+ TSbulk_insert.go
  > `BulkInsert`批量插入机构信息函数，`TpBulkInsert`批量插入类型信息。
+ TSmapping.go
  > 索引创建函数，`TsMapping`为培训机构索引信息,`TpMapping`为类型索引信息。
### main
main函数，开启爬虫
