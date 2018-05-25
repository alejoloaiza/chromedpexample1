package main

import (
	"fmt"
	"strings"

	"github.com/schollz/closestmatch"
)

func SplitPrices(input string) (string, string) {

	splitted := strings.Split(input, "$")
	lastpart := splitted[1]
	price := strings.Split(lastpart, " ")
	fmt.Println(price[0], strings.Join(price[1:], " "))
	return price[0], strings.Join(price[1:], " ")
}

// We define the Asset.Type based on similarity of words using schollz/closestmatch
func AssetClassifier(Api string) {
	TypeStuff := []string{"Apartamento", "Casa", "Bodega", "Finca", "Oficina", "Local"}
	BusinessStuff := []string{"Arrendar", "Vender"}
	bagSizes := []int{2, 3, 4, 5}

	cmType := closestmatch.New(TypeStuff, bagSizes)
	cmBusiness := closestmatch.New(BusinessStuff, bagSizes)
	fmt.Println(cmType, cmBusiness)

}
