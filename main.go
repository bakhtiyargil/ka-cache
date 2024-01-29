package main

import "fmt"

func main() {

	var cac = NewCache(5)
	cac.Set("a", "neymar")
	cac.Set("b", "jr")
	cac.Set("b", "valdes")
	cac.Set("b", "iniesta")
	cac.Set("c", "nunez")
	for key, value := range cac.cacheMap {
		fmt.Printf("%s: %s\n", key, value.value)
	}

	fmt.Println("head: " + cac.head.value)
	cac.Get("a")
	fmt.Println("head: " + cac.head.value)

	fmt.Println("Ka-cache!")

}
