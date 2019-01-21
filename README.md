# burgoking

🍔 **Burger King - Free Burger Code Generator** 🍔

Generate a Burger King's promotion code to get a free burger using Golang.


### Installation

`$ go get github.com/Scotow/burgoking`


### Examples


#### Library

Thelibrary provides two ways to generate codes. A single code can be generating or an auto refreshing pool can be used.


##### Generating a single code

```go
func GenerateCode(meal *Meal) (code string, err error)
```

The `meal` argument is used to fill the first page of the survey.

```go
type Meal struct {
	Restaurant int
	Date       time.Time
}
```

Passing `nil` as an argument generate a random meal using the `func RandomMeal() *Meal` function.

##### Using a pool of codes

A pool of codes holds a fixed amount of codes. When a code of the pool is consumed by the function `GetCode`, a new goroutine is spawned to generate a new code.

An expiration duration can be specified. When a code has stayed in the pool for too long, a timer removes it from the pool and a new generation is launched.

To create a new pool, you may call the following function:

```go
func NewPool(size int, expiration, retry time.Duration) (pool *Pool, err error)
```

Where
* `size` is the total number of code in the pool.
* `expiration` is the duration required for a code to be remove and auto re-generated.
* `retry` is the duration required between two calls of the `GenerateCode` function if the first call failed for any reason.

#### Binaries

The [cmd](https://github.com/Scotow/burgoking/blob/master/cmd) folder contains three examples of program that use the `burgoking` library.

##### Simple command



### Contribution

Feedback are appreciated. Feel free to open an issue or a pull request if needed.

Furthermore, if you went to restaurant which its number isn't in the restaurants [list](https://github.com/Scotow/burgoking/blob/master/meal.go#L9), a merge request to add it is appreciated.


### Disclaimer

*burgoking* provided by *Scotow* is for illustrative purposes only which provides customers with programming information regarding the products. This software is supplied "AS IS" without any warranties and support.

I assumes no responsibility or liability for the use of the software, conveys no license or title under any patent, copyright, or mask work right to the product.

***Enjoy your meal!***
