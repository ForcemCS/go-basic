package main

import "fmt"

type Product struct {
	Name     string
	Price    float64
	Category string //种类，范畴
	Stock    int    //库存
}

func FilterFunc[T any](items []T, match func(item T) bool) []T {
	var result []T

	for _, item := range items {
		if match(item) {
			result = append(result, item)
		}

	}

	return result

}
func main() {

	products := []Product{
		{"苹果手机", 5999, "电子产品", 10},
		{"鼠标", 99, "电子产品", 0},
		{"水杯", 29, "日用品", 50},
		{"机械键盘", 499, "电子产品", 0},
	}

	//需求 1：找出价格大于 100 的商品
	expensiveProducts := FilterFunc(products, func(p Product) bool {
		return p.Price > 100
	})
	fmt.Println("贵商品:", expensiveProducts)

	outOfStock := FilterFunc(products, func(p Product) bool {
		return p.Stock == 0
	})

	fmt.Println("没库存的:", outOfStock)

}
