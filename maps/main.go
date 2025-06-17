package main

func main() {
	// This is a placeholder for the main function.
	// You can implement your map-related logic here.
	// For example, you might want to initialize a map,
	// add some key-value pairs, and print them out.

	myMap := make(map[string]int)
	myMap["apple"] = 5
	myMap["banana"] = 3
	myMap["orange"] = 8

	for key, value := range myMap {
		println(key, ":", value)
	}

	println(myMap["abhinav"])
	value, exists := myMap["abhinav"]
	if exists {
		println("abhinav exists with value:", value)
	} else {
		println("abhinav does not exist in the map.")
	}

}
