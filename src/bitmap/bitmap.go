package bitmap

import "fmt"

type BitMap []byte

const byteSize = 8 //定义的bitmap为byte的数组，byte为8bit

func NewBitMap(n uint) BitMap { // n为数据范围的最大值，例如0~1000，n=1000
	return make([]byte, n/byteSize+1)
}
func (bt BitMap) set(n uint) {
	if n/byteSize > uint(len(bt)) {
		fmt.Println("大小超出bitmap范围")
		return
	}
	byteIndex := n / byteSize   //第x个字节（0,1,2...）
	offsetIndex := n % byteSize //偏移量(0<偏移量<byteSize)
	//bt[byteIndex] = bt[byteIndex] | 1<<offsetIndex //异或1（置位）
	//第x个字节偏移量为offsetIndex的位 置位1
	bt[byteIndex] |= 1 << offsetIndex //异或1（置位）
}
func (bt BitMap) del(n uint) {
	if n/byteSize > uint(len(bt)) {
		fmt.Println("大小超出bitmap范围")
		return
	}
	byteIndex := n / byteSize
	offsetIndex := n % byteSize
	bt[byteIndex] &= 0 << offsetIndex //清零
}
func (bt BitMap) isExist(n uint) bool {
	if n/byteSize > uint(len(bt)) {
		fmt.Println("大小超出bitmap范围")
		return false
	}
	byteIndex := n / byteSize
	offsetIndex := n % byteSize
	//fmt.Println(bt[byteIndex] & (1 << offsetIndex))
	return bt[byteIndex]&(1<<offsetIndex) != 0 //注意：条件是 ！=0，有可能是：16,32等
}

// 实现快速排序
func quickSort(nums []int, start, end int) []int {
	if start < end {
		i, j := start, end
		key := nums[(start+end)/2]
		for i <= j {
			for nums[i] < key {
				i++
			}
			for nums[j] > key {
				j--
			}
			if i <= j {
				nums[i], nums[j] = nums[j], nums[i]
				i++
				j--
			}
		}
		if start < j {
			nums = quickSort(nums, start, j)
		}
		if end > i {
			nums = quickSort(nums, i, end)
		}

	}
	return nums
}

func bucketSort(nums []int, bucketNum int) []int {
	bucket := [][]int{}
	for i := 0; i < bucketNum; i++ {
		tmp := make([]int, 1)
		bucket = append(bucket, tmp)
	}

	//将数据分配到桶中
	for i := 0; i < len(nums); i++ {
		bucket[nums[i]/bucketNum] = append(bucket[nums[i]/bucketNum], nums[i])
	}

	//循环所有的桶进行排序
	index := 0
	for i := 0; i < bucketNum; i++ {
		if len(bucket[i]) > 1 {
			//对每个桶内部进行排序,使用快排
			bucket[i] = quickSort(bucket[i], 0, len(bucket[i])-1)
			for j := 1; j < len(bucket[i]); j++ { //去掉一开始的tmp
				nums[index] = bucket[i][j]
				index++
			}
		}
	}
	return nums
}
