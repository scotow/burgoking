# burgoking

⚠️ As of September 2022, Burger King migrated their survey system to a third party company, resulting to changes that require an almost complete rebuild of this project. Unfortunately, I cannot find the time nor the will to do it. So, I replaced the generation algorithm with something local and static. In almost all cases, this generator will work (I never encountered a BK employee that checked the numbers). ⚠️

🍔 **Burger King - Free Burger Code Generator** 🍔

Generate a Burger King's promotion code to get a free burger using Golang.


### Installation

`$ go get github.com/scotow/burgoking`


### Examples


#### Library

The library provides two ways to generate codes. A single code can be generating or an auto refreshing pool can be used.


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
* `size` is the total number of codes in the pool.
* `expiration` is the duration required for a code to be remove and auto re-generated.
* `retry` is the duration required between two calls of the `GenerateCode` function if the first call failed for any reason. This prevents the pool from spaming the website.

#### Binaries

The [cmd](https://github.com/Scotow/burgoking/blob/master/cmd) folder contains three examples of program that use the `burgoking` library.

##### Simple code

The [code](https://github.com/Scotow/burgoking/blob/master/cmd/code) command is a simple wrapper of the library. It generates a code and prints it to `stdout`.

Here are the available options:

```
Usage of code:
  -n int
    	number of code to generate (default 1)
  -p	
  	generate each code on a different goroutine (default false)
```

##### Simple API

The [api](https://github.com/Scotow/burgoking/blob/master/cmd/api) command starts a simple web server that returns a new promotion code to each request.

##### (Old) Web friendly server

The [web](https://github.com/Scotow/burgoking/blob/master/cmd/web_original) command is a ready-for-demo binary that serves a simple, yet beautiful UI generating promotion codes using a pool.

The program may setup a second private pool that requires a `Authorization` HTTP header.

Here are the available options:

```
Usage of web:
  -D string
    	private code expiration (default "24h0m0s")
  -N int
    	private code pool size (default 1)
  -R string
    	private code regeneration interval (default "30s")
  -c string
    	contact address on error
  -d string
    	public code expiration (default "24h0m0s")
  -k string
    	authorization token for private and direct code (disable if empty)
  -n int
    	public code pool size (default 3)
  -p int
    	listening port (default 8080)
  -r string
    	public code regeneration interval (default "30s")
```

### Contribution

Feedbacks are appreciated. Feel free to open an issue or a pull request if needed.

Furthermore, if you went to a restaurant which its number isn't in the restaurants [list](https://github.com/Scotow/burgoking/blob/master/meal.go#L9), a merge request to add it is appreciated.


### Disclaimer

*burgoking* provided by *Scotow* is for illustrative purposes only which provides customers with programming information regarding the products. This software is supplied "AS IS" without any warranties and support.

I assumes no responsibility or liability for the use of the software, conveys no license or title under any patent, copyright, or mask work right to the product.

### JavaScript version

A deprecated and no more maintained version of this project made with JavaScript (NodeJS) and CasperJS is availiable [here](https://github.com/Scotow/bk-generator).

***Enjoy your meal!***
