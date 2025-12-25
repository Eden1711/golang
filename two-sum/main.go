package main

import "fmt"

func twoSum(nums []int, target int) []int {
	numToIndexMap := make(map[int]int)

	for i, num := range nums {
		diff := target - num
		idx, found := numToIndexMap[diff]
		if found {

			return []int{i, idx}
		}
		numToIndexMap[num] = i
	}

	return nil
}

func main() {
	result := twoSum([]int{2, 7, 11, 15}, 9)

	fmt.Println(result)

}
