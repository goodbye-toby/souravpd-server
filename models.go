package main
import(
	"fmt"
)
type Blog struct{
	Title string `json:"title"`
	Description string `json:"description"`
	Html string `json:"html"`
	CreatedAt string `json:"date"`
}

func main(){
	fmt.Println(Blog{"k","a","b","h"})
}
