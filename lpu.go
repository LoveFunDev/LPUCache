package LPUCache

import (
	"container/list"
	"time"
)

type nodeContent struct {
	key         uint32
	value       uint32
	createdTime int64
}

// LPUCache defines the LPU data structure
type LPUCache struct {
	capacity uint32
	nodeMap  map[uint32]*list.Element // hashMap is used to read the element by Key, the complexity is O(1)
	nodeList *list.List               // doubly linked list is used to store elements and adjust their recently usage
}

const (
	TTL = 15 * time.Minute // TTL is the Time-To-Live
)

// NewLPUCache returns a LPU cache that handles LPU operations
func NewLPUCache(maxCapacity uint32) (*LPUCache, error) {
	cache, err := LoadLPUCache()
	if err != nil {
		newCache := &LPUCache{
			capacity: maxCapacity,
			nodeMap:  make(map[uint32]*list.Element, 0),
			nodeList: list.New(),
		}
		return newCache, nil
	}

	return cache, nil
}

// Get returns the value of the key
func (lpu *LPUCache) Get(key uint32) int {
	if _, exist := lpu.nodeMap[key]; !exist { // Key does not exist
		return -1
	}

	defer SaveLPUCache(lpu)

	nodeItem := lpu.nodeMap[key]
	_, nodeItemValue := getNodeElementKeyValue(nodeItem)

	if isNodeExpired(nodeItem) {
		lpu.removeNodeFromListAndMap(nodeItem)
		return -1
	}

	lpu.nodeList.MoveToFront(nodeItem)
	return int(nodeItemValue)
}

// Put Sets or insert the value
func (lpu *LPUCache) Put(key uint32, value uint32) {
	defer SaveLPUCache(lpu)

	if _, exist := lpu.nodeMap[key]; exist { // Key already exists, just assign the Value and move the element to the front of list
		nodeItem := lpu.nodeMap[key]
		setNodeElementValue(nodeItem, value)
		lpu.nodeList.MoveToFront(nodeItem)
		return
	}

	if lpu.checkIfWillExceedCapacity() { // If inserting node reaches Capacity, evicts the back element(which is the least recently used)
		tailNodeItem := lpu.nodeList.Back()
		lpu.removeNodeFromListAndMap(tailNodeItem)
	}

	newNodeInfo := nodeContent{key: key, value: value, createdTime: time.Now().Unix()}
	nodeItem := lpu.nodeList.PushFront(newNodeInfo)
	lpu.nodeMap[key] = nodeItem
}

// Delete removes the value of the key if the key exists in the LPU cache
func (lpu *LPUCache) Delete(key uint32) int {
	if _, exist := lpu.nodeMap[key]; !exist { // Key does not exist
		return -1
	}

	nodeItem := lpu.nodeMap[key]
	_, nodeItemValue := getNodeElementKeyValue(nodeItem)
	lpu.removeNodeFromListAndMap(nodeItem)
	return int(nodeItemValue)
}

func (lpu *LPUCache) removeNodeFromListAndMap(nodeItem *list.Element) {
	nodeKey, _ := getNodeElementKeyValue(nodeItem)
	lpu.nodeList.Remove(nodeItem)
	delete(lpu.nodeMap, nodeKey)
}

func (lpu *LPUCache) checkIfWillExceedCapacity() bool {
	return uint32(lpu.nodeList.Len())+1 > lpu.capacity
}

func getNodeElementKeyValue(nodeItem *list.Element) (uint32, uint32) {
	content := nodeItem.Value.(nodeContent)
	return content.key, content.value
}

func setNodeElementValue(nodeItem *list.Element, newValue uint32) {
	content := nodeItem.Value.(nodeContent)
	content.value = newValue
}

func isNodeExpired(nodeItem *list.Element) bool {
	content := nodeItem.Value.(nodeContent)
	return time.Duration(time.Now().Unix()-content.createdTime) >= TTL
}
