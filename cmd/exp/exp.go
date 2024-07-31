package main

import "fmt"

func main() {
	arr := []int{4, 7, 2, 9, 7, 4, 9}
	
	// Count occurrences of each number
	countMap := make(map[int]int)
	for _, num := range arr {
		countMap[num]++
	}

	fmt.Println(countMap)
	
	// Create a new slice with numbers that don't repeat odd times
	result := []int{}
	for _, num := range arr {
		if countMap[num]%2 == 0 {
			result = append(result, num)
		}
	}

	fmt.Println(result)
}


