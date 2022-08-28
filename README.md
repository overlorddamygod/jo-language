# Jo Programming Language

A toy programming language built with golang.

---
### Usage

```console
user@main:~$ go run main.go example.jo
```
or

```console
user@main:~$ ./jo example.jo
```
---
### Sample Program
```js
name = "Overlord";

if (name == "") {
    name = "World";
}

helloworld = "Hello, " + name;

// Printing Hello, World
print(helloworld);

// Print if number is even or odd from 0 to 10
for ( i = 0; i <= 10; i = i + 1 ) {
    if ( i % 2 == 0 ) {
        print("is even " + i);
    } else {
        print("is odd " + i);
    }
}
```