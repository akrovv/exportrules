## Exportes linter

#### Used for verification using private fields/methods

You only need to use the public ** fields** / **methods** in one package.
<table>
<thead><tr><th>Bad Field</th><th>Bad Method</th></tr></thead>
<tbody>
<tr><td>

```go
type User struct {
    Name string
    age int
}

func getAge(user User) int { 
    return user.age
}
```

</td><td>

```go
type User struct {
    Name string
    age int
}

func (u User) setAge(age int) {
    u.age = age
}

func setter[T any](user User, value T) { 
     user.setAge(value)
}
```

</td></tr>
</tbody></table>


### How to start
* git clone __[link]__
* go build -buildmode=plugin cmd/exportes/main.go
* __[!!! you need compiled file from https://github.com/golangci/golangci-lint (make build)]__
* ./golangci-lint -c golangci.yml run __[source code]__