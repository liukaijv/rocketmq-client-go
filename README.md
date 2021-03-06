## 遇到的问题

使用protobuf的proto.Marshal后的数据传输被截断

```
# api.proto

message RegisterServerRequest{
    required int32      id      =1;
    required string     type    =2;
    required string     key     =3;
    required string     name    =4;
    required int32      port    =5;
    required string     ip      =6;
    optional int32      onLine  =7;
}

# produce.go

# 生产者
registerRequest := api.RegisterServerRequest{
    Id:   proto.Int(1),
    Type: proto.String("edge"),
    Key:  proto.String("123456"),
    Name: proto.String("server1"),
    Port: proto.Int(0),
    Ip:   proto.String("127.0.0.1"),
}
data, err := proto.Marshal(&registerRequest)
if err != nil {
    fmt.Println("proto.Marshal error:", err)
}

msg := &rocketmq.Message{Topic: topic, Body: string(data)}

result, err := producer.SendMessageSync(msg)
if err != nil {
    fmt.Println("Error:", err)
}

# 发送的数据

00000000  08 01 12 04 65 64 67 65  1a 06 31 32 33 34 35 36  |....edge..123456|
00000010  22 07 73 65 72 76 65 72  31 28 00 32 09 31 32 37  |".server1(.2.127|
00000020  2e 30 2e 30 2e 31                                 |.0.0.1|

# 实际传输的数据为

00000000  08 01 12 04 65 64 67 65  1a 06 31 32 33 34 35 36  |....edge..123456|
00000010  22 07 73 65 72 76 65 72  31 28 00                 |".server1(.     |

# 到c++层时，二进制数据遇到0x00被截断了

```

## 解决

* 现在把：Message.Body修改为[]bye类型，生产时往Message.property写一个key为_bodyLen标记发送的[]byte长度，消费时取出指定长度的[]byte； 如果使用跨语言，生产者property里每次需要写入数据和度
* 也可以不修改，需要proto.Marshal后再base64一下，但这样多了一次base64的decode和encode


## RocketMQ Client Go
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![TravisCI](https://travis-ci.org/apache/rocketmq-client-go.svg)](https://travis-ci.org/apache/rocketmq-client-go)

( The alpha version of Apache RocketMQ Go in Pure Go has been released, Welcome have a try! [native version](https://github.com/apache/rocketmq-client-go/tree/native) )

* The client is using cgo to call [rocketmq-client-cpp](https://github.com/apache/rocketmq-client-cpp), which has been proven robust and widely adopted within Alibaba Group by many business units for more than three years.

----------
## Features
At present, this SDK supports
* sending message in synchronous mode
* sending message in orderly mode
* sending message in oneway mode
* consuming message using push model
* consuming message using pull model

----------
## How to use
* Step-by-step instruction are provided in [RocketMQ Go Client Introduction](./doc/Introduction.md)
* Consult [RocketMQ Quick Start](https://rocketmq.apache.org/docs/quick-start/) to setup rocketmq broker and nameserver.

----------
## Apache RocketMQ Community
* [RocketMQ Community Projects](https://github.com/apache/rocketmq-externals)

----------
## Contact us
* Mailing Lists: <https://rocketmq.apache.org/about/contact/>
* Home: <https://rocketmq.apache.org>
* Docs: <https://rocketmq.apache.org/docs/quick-start/>
* Issues: <https://github.com/apache/rocketmq-client-go/issues>
* Ask: <https://stackoverflow.com/questions/tagged/rocketmq>
* Slack: <https://rocketmq-community.slack.com/>
 
---------- 
## How to Contribute
  Contributions are warmly welcome! Be it trivial cleanup, major new feature or other suggestion. Read this [how to contribute](http://rocketmq.apache.org/docs/how-to-contribute/) guide for more details. 
   
   
----------
## License
  [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html) Copyright (C) Apache Software Foundation
