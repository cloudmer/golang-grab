# golang-grab
golang 彩票分析

重庆,天津,新疆 时时彩,台湾五分彩 一些算法  

采用golang 多并发 来分析不同彩种开奖数据

php 只用来抓取各彩种开奖数据

golang 分析 开奖数据

安装步骤:

 * 安装golang 环境 
 * cd $GOPATH/src
 * git clone https://github.com/aiyeyun/golang-grab.git "xmn"
 * go get github.com/go-gomail/gomail
 * go get github.com/go-ini/ini
 * go get github.com/go-sql-driver/mysql
 * go install
 * cd $GOPATH/bin
 * ./xmn 运行