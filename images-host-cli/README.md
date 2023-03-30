
```go
package goPaste_test

import (
	"fmt"
	"goPaste"
	"os"
)

func Example() {
	text, err := goPaste.PasteTXT()
	if err == nil {
		fmt.Println(text)
	} else {
		fmt.Println(err)
	}

	img, imgErr := goPaste.PasteImg(true)
	if imgErr == nil {
		f, err := os.OpenFile("output.png", os.O_WRONLY, 0666)
		if err != nil {
			f, err = os.Create("output.png")
		}
		if err == nil {
			f.Write(img)
			f.Close()
		}
	} else {
		fmt.Println(imgErr)
	}

	img, imgErr = goPaste.PasteImg(false)
	if imgErr == nil {
		f, err := os.OpenFile("output.bmp", os.O_WRONLY, 0666)
		if err != nil {
			f, err = os.Create("output.bmp")
		}
		if err == nil {
			f.Write(img)
			f.Close()
		}
	} else {
		fmt.Println(imgErr)
	}
}
```