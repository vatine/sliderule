package main

import (
	//"bytes"
	"fmt"
	"os"
	"sliderule/lib"
)

func main () {
	out, err := os.Create("/Users/ingvar/tmp/slide.svg")
	if err != nil {
		fmt.Println(err)
		return
	}
	scale := sliderule.MakeSlideRule(250.0, 50.0)
	scale.Render(out)
	out.Close()
}
