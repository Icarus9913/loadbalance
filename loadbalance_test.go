package loadbalance

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type LoadBalance interface {
	Add(...string) error
	Get(string)(string, error)
}

const (
	LbRandom int = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbType int) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LbConsistentHash:
		return NewConsistentHashBalance(10, nil)
	default:
		return &RandomBalance{}
	}
}

func TestLoadBalance(t *testing.T)  {
	rand.Seed(time.Now().Unix())

	var err error
	var result string
	sum := []string{}

	lbType := rand.Intn(4)
	lbType=3
	loadBalance := LoadBalanceFactory(lbType)
	for i:=0;i<5;i++{
		addr := fmt.Sprintf("192.168.%d.%d",rand.Intn(255),rand.Intn(255))
		sum = append(sum,addr)
		if lbType == 2{
			err = loadBalance.Add(addr,strconv.Itoa(rand.Intn(10)))	//权重也可以设置成定值
		}else {
			err = loadBalance.Add(addr)
		}
		if nil!=err{
			panic(err)
		}
	}
	fmt.Println(sum)

	_,ok := loadBalance.(*ConsistentHashBalance);
	for i:=0;i<5;i++{
		if ok{
			result,err = loadBalance.Get(strconv.Itoa(i))
		}else {
			result, err = loadBalance.Get("")
	}

	if nil!=err{
		panic(err)
	}
	fmt.Println(result)
	}
}
