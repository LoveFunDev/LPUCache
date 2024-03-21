package LPUCache

import (
	"container/list"
	"errors"
	"fmt"
)

// LPUConfigNode stores node content
type LPUConfigNode struct {
	Key         uint32 `json:"Key"`
	Value       uint32 `json:"Value"`
	CreatedTime int64  `json:"CreatedTime"`
}

// LPUConfigCache stores LPU cache
type LPUConfigCache struct {
	Capacity       uint32           `json:"Capacity"`
	Nodes          []*LPUConfigNode `json:"Nodes"`
	MappedNodeKeys []uint32         `json:"MappedNodeKeys"`
}

const (
	ConfigFile = "/etc/lpuCache_config.json"
)

// LoadLPUCache loads LPUCache from configuration file
func LoadLPUCache() (*LPUCache, error) {
	configCache, err := loadLPUConfigCache()
	if err != nil {
		return nil, err
	}
	return convertConfigCache(configCache)
}

// SaveLPUCache stores LPU Cache
func SaveLPUCache(lpu *LPUCache) {
	configCache := &LPUConfigCache{
		Capacity: lpu.capacity,
	}

	configCache.Nodes = []*LPUConfigNode{}
	for e := lpu.nodeList.Front(); e != nil; e.Next() {
		nodeItem := e.Value.(nodeContent)
		configNode := &LPUConfigNode{Key: nodeItem.key, Value: nodeItem.value, CreatedTime: nodeItem.createdTime}
		configCache.Nodes = append(configCache.Nodes, configNode)
	}

	configCache.MappedNodeKeys = []uint32{}
	for key := range lpu.nodeMap {
		configCache.MappedNodeKeys = append(configCache.MappedNodeKeys, key)
	}

	if err := Store(ConfigFile, configCache); err != nil {
		fmt.Println(err)
	}
}

func loadLPUConfigCache() (*LPUConfigCache, error) {
	retrievedData := &LPUConfigCache{}

	err := Fetch(ConfigFile, retrievedData)
	if err != nil {
		return nil, err
	}

	return retrievedData, nil
}

func convertConfigCache(configCache *LPUConfigCache) (*LPUCache, error) {
	cache := &LPUCache{
		capacity: configCache.Capacity,
	}

	cache.nodeList = list.New()
	for _, node := range configCache.Nodes {
		newNodeInfo := nodeContent{key: node.Key, value: node.Value, createdTime: node.CreatedTime}
		_ = cache.nodeList.PushBack(newNodeInfo)
	}

	cache.nodeMap = make(map[uint32]*list.Element, 0)
	for _, key := range configCache.MappedNodeKeys {
		node, err := findNodeByKey(key, cache.nodeList)
		if err != nil {
			return nil, err
		}
		cache.nodeMap[key] = node
	}
	return cache, nil
}

func findNodeByKey(key uint32, l *list.List) (*list.Element, error) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(nodeContent).key == key {
			return e, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find Key: %d", key))
}
