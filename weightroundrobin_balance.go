package loadbalance

import (
	"errors"
	"strconv"
)

/*
	加权轮询负载
	1. currentWeight = currentWeight + effectiveWeight
	2. 选中最大的currentWeight节点为选中节点
	3. currentWeight = currentWeight - totalWeight
	totalWeight = sum(effectiveWeight)
*/
type WeightNode struct {
	addr            string
	Weight          int // 初始化时对节点约定的权重
	currentWeight   int // 节点临时权重，每轮都会变化
	effectiveWeight int // 有效权重，默认与weight相同
}

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
}

func (w *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("params len require 2")
	}
	parseInt, err := strconv.ParseInt(params[1], 10, 64)
	if nil != err {
		return err
	}
	node := &WeightNode{
		addr:   params[0],
		Weight: int(parseInt),
	}
	node.effectiveWeight = node.Weight
	w.rss = append(w.rss, node)
	return nil
}

func (w *WeightRoundRobinBalance) Next() string {
	var best *WeightNode
	total := 0
	for i := 0; i < len(w.rss); i++ {
		weight := w.rss[i]
		// 1. 计算有效权重
		total += weight.effectiveWeight
		// 2. 修改当前节点临时权重
		weight.currentWeight += weight.effectiveWeight
		// 3. 有效权重默认与权重相同，通讯异常时-1，通讯成功+1，直到恢复到weight大小
		if weight.effectiveWeight < weight.Weight {
			weight.effectiveWeight++
		}
		// 4. 选中最大临时权重节点
		if nil == best || weight.currentWeight > best.currentWeight {
			best = weight
		}
	}
	if nil == best {
		return ""
	}
	// 5. 变更临时权重为 临时权重-有效权重之和
	best.currentWeight -= total
	return best.addr
}

func (w *WeightRoundRobinBalance) Get(string) (string, error) {
	addr := w.Next()
	if addr == "" {
		return "", errors.New("No data")
	}
	return addr, nil
}
