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