#!/bin/bash

# 清除链码容器 正则：/dev-peer.*.blockchain-real-estate.*/ 匹配上的都会被删除，其中blockchain-real-estate是链码名称，在安装和实例化的时候会指定
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*.blockchain-real-estate.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# 清除不需要的链码镜像 正则：/dev-peer.*.blockchain-real-estate.*/ 匹配上的都会被删除，其中blockchain-real-estate是链码名称，在安装和实例化的时候会指定
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*.blockchain-real-estate.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

echo "区块链 ： 关闭"

echo "开始删除链码生成的docker镜像"
docker-compose down --volumes --remove-orphans

# 调用函数清除链码容器
clearContainers

# 调用函数清除不需要的链码镜像
removeUnwantedImages



echo "一、环境清理"
mkdir -p config
mkdir -p crypto-config
rm -fr config/*
rm -fr crypto-config/*
echo "清理完毕"


echo "二、生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block

echo "三、区块链：启动"
docker-compose up -d        # 按照docker-compose.yaml的配置启动区块链网络并在后台运行
echo "正在等待节点的启动完成，等待10秒"
sleep 10                    # 启动整个区块链网络需要一点时间，所以此处等待10s，让区块链网络完全启动

echo "四、生成通道的TX文件(这个动作会创建一个创世交易，也是该通道的创世交易)"
configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/assetschannel.tx -channelID assetschannel

# 五、在区块链上按照刚刚生成的TX文件去创建通道
# 该操作和上面操作不一样的是，这个操作会写入区块链
echo "五、在区块链上按照刚刚生成的TX文件去创建通道"
docker exec cli peer channel create -o orderer0.blockchainrealestate.com:7050 -c assetschannel -f /etc/hyperledger/config/assetschannel.tx --tls --cafile /etc/hyperledger/orderer/msp/tlscacerts/tlsca.blockchainrealestate.com-cert.pem

echo "六、让节点去加入到通道"
echo "peer0.org1加入到通道"
docker exec cli peer channel join -b assetschannel.block

echo "peer1.org1加入到通道"
docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/assetschannel.block .
docker cp assetschannel.block cli1:/assetschannel.block
docker exec cli1 peer channel join -b /assetschannel.block

# 七、链码安装
# -n 是链码的名字，可以自己随便设置
# -v 就是版本号，就是composer的bna版本
# -p 是目录，目录是基于cli这个docker里面的$GOPATH相对的
# 此处安装的是示例链码，后续课程会自己编写
echo "七、链码安装"
docker exec cli peer chaincode install -n blockchain-real-estate -v 1.0.0 -l golang -p chaincode
docker exec cli1 peer chaincode install -n blockchain-real-estate -v 1.0.0 -l golang -p chaincode

#八、实例化链码
#-n 对应前文安装链码的名字 其实就是composer network start bna名字
#-v 为版本号，相当于composer network start bna名字@版本号
#-C 是通道，在fabric的世界，一个通道就是一条不同的链，composer并没有很多提现这点，composer提现channel也就在于多组织时候的数据隔离和沟通使用
#-c 为传参，传入init参数
echo "八、实例化链码"
docker exec cli peer chaincode instantiate -o orderer0.blockchainrealestate.com:7050 -C assetschannel -n blockchain-real-estate -l golang -v 1.0.0 --tls --cafile /etc/hyperledger/orderer/msp/tlscacerts/tlsca.blockchainrealestate.com-cert.pem -c '{"Args":["init"]}'

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":[""]}'

echo "九、文件权限设置"
sudo chmod -R 777 ../deploy/
