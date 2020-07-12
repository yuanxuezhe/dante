use dante;

drop table userinfo;

CREATE TABLE IF NOT EXISTS `userinfo`(
`userid`    int NOT NULL,
`username`  VARCHAR(32) NOT NULL,
`passwd`    VARCHAR(16) NOT NULL,
`sex`       char  NOT NULL,
`phone`     bigint default 0,
`email`     VARCHAR(100),
`status`    char  NOT NULL default '0' ,
`registerdate` int default 0,
PRIMARY KEY ( `userid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
insert into `userinfo` (`userid`,`username`,`passwd`,`sex`,`phone`,`email`,`status`,`registerdate`)
values (1, 'woshinibaba', '1', 0,18664324256,'446968454@qq.com',0,'20200712');

drop table goods;
CREATE TABLE IF NOT EXISTS `goods`(
    `goodsid`   BIGINT(20)  default 0       PRIMARY KEY COMMENT '编号',
    `goodsname`       VARCHAR(128) NOT NULL COMMENT '名称',
    `type`      int default 0               COMMENT '商品类型',
    `source`    int default 0               COMMENT '来源',
    `url`     VARCHAR(128) NOT NULL         COMMENT '链接',
    `imgurl`    VARCHAR(128) NOT NULL       COMMENT '图片链接',
    `brand`     int default 0               COMMENT '品牌',
    `status`    int default 0               COMMENT '状态',
    `date`      int default 0               COMMENT '日期',
    `time`      int default 0               COMMENT '时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (1, 'goods01', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/911a7002e11608fb581fffcde05d5257.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (2, 'goods02', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/2_1536049430.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (3, 'goods03', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/911a7002e11608fb581fffcde05d5257.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (4, 'goods04', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/2_1536049430.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (5, 'goods05', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/911a7002e11608fb581fffcde05d5257.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (6, 'goods06', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/2_1536049430.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (7, 'goods07', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/911a7002e11608fb581fffcde05d5257.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (8, 'goods08', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/2_1536049430.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (9, 'goods09', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/2_1536049430.jpg');
insert into `goods` (`goodsid`,`goodsname`,`type`,`source`,`url`,`imgurl`) values (10, 'goods10', 1, 0,'https://item.jd.com/100012545852.html','https://boweisou.oss-cn-shenzhen.aliyuncs.com/yy/images/911a7002e11608fb581fffcde05d5257.jpg');

drop table goods_describ;
CREATE TABLE IF NOT EXISTS `goods_describ`(
    `goodsid`   int  default 0              PRIMARY KEY COMMENT '编号',
    `describe` blob                        COMMENT '描述'
   )ENGINE=InnoDB DEFAULT CHARSET=utf8;