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

CREATE TABLE IF NOT EXISTS `goods`(
    `goodsid`   BIGINT(20)  default 0       PRIMARY KEY COMMENT '编号',
    `goodsname`       VARCHAR(128) NOT NULL COMMENT '名称',
    `type`      int default 0               COMMENT '商品类型',
    `source`    int default 0               COMMENT '来源',
    `url`     VARCHAR(64) NOT NULL          COMMENT '链接',
    `imgurl`    VARCHAR(64) NOT NULL        COMMENT '图片链接',
    `brand`     int default 0               COMMENT '品牌',
    `status`    char  NOT NULL default '0'  COMMENT '状态',
    `date`      int default 0               COMMENT '日期',
    `time`      int default 0               COMMENT '时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `goods_describ`(
    `goodsid`   int  default 0              PRIMARY KEY COMMENT '编号',
    `describe` blob                        COMMENT '描述'
   )ENGINE=InnoDB DEFAULT CHARSET=utf8;