go-rucaptha
-----------

Go-rucaptha is the package to work with API rucaptcha (https://rucaptcha.com)


Installation
------------

```
go get github.com/allush/go-rucaptha
```


Usage
-----

Import package into your project:

```
import "github.com/allush/go-rucaptha"
```

Create new instance of rucaptcha solver:

```
solver := rucaptcha.New("your api key")
```

Call `Solve` method for get answer:

```
answer, _ := solver.Solve("https://raw.githubusercontent.com/allush/go-rucaptha/master/test/captcha.jpg")
```

You can specify additional options for solver by:
```
solver.IsRegsence = true
```
For more, see CaptchaSolver struct definition (https://github.com/allush/go-rucaptha/blob/master/solver.go) and API rucapthca (https://rucaptcha.com/api-rucaptcha)

Example
-----

```
package main

import (
   "github.com/allush/go-rucaptha"
   "fmt"
  )

func main() {
  solver := rucaptcha.New("your api key")
  answer, err := solver.Solve("https://raw.githubusercontent.com/allush/go-rucaptha/master/test/captcha.jpg")
  if err != nil {
    fmt.Println("Error:", err)
  }

  // use the answer
  fmt.Println(*answer)
}
```


Authors
-------

  * [Alexey Lushnikov](https://github.com/allush)
