# ipapk
ipa or apk parser written in golang, aims to extract app information

## INSTALL
	$ go get github.com/phinexdaz/ipapk
  
## USE
```go
package main

import (
	"fmt"
	"github.com/phinexdaz/ipapk"
)

func main() {
	apk, _ := ipapk.NewAppParser("test.apk")
	fmt.Println(apk)
}
```
