package LPUCache

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_LPUCache_Basic(t *testing.T) {
	we := NewGomegaWithT(t)

	cache := NewLPUCache(2)

	cache.Put(1, 1)
	we.Expect(cache.nodeList.Len()).To(Equal(1))

	cache.Put(2, 2)
	we.Expect(cache.nodeList.Len()).To(Equal(2))

	ret := cache.Get(1) // returns 1
	we.Expect(ret).To(Equal(1))

	cache.Put(3, 3) // evicts key 2
	we.Expect(cache.nodeList.Len()).To(Equal(2))
	_, exist := cache.nodeMap[2]
	we.Expect(exist).To(BeFalse())

	ret = cache.Get(2) // return -1
	we.Expect(ret).To(Equal(-1))

	cache.Put(4, 4) // evicts key 1
	we.Expect(cache.nodeList.Len()).To(Equal(2))
	_, exist = cache.nodeMap[1]
	we.Expect(exist).To(BeFalse())

	ret = cache.Get(1) // return -1
	we.Expect(ret).To(Equal(-1))

	ret = cache.Get(3) // returns 3
	we.Expect(ret).To(Equal(3))

	ret = cache.Get(4) // returns 4
	we.Expect(ret).To(Equal(4))

	ret = cache.Delete(3) // returns 3
	we.Expect(cache.nodeList.Len()).To(Equal(1))
	_, exist = cache.nodeMap[3]
	we.Expect(exist).To(BeFalse())
	we.Expect(ret).To(Equal(3))

	ret = cache.Get(3) // returns -1
	we.Expect(ret).To(Equal(-1))
}
