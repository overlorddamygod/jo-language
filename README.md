# Jo Programming Language

A toy programming language built with golang.

Language Grammar in [`jo-grammar.cfg`](./jo-grammar.cfg)

### References:
* [Crafting Interpreters](https://craftinginterpreters.com/appendix-i.html)
    
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

hasName = name == "";

if (!!hasName) {
    name = "World";
}

// Printing Hello, World
print("Hello,", name);

// Print if number is even or odd from 0 to 10
for ( i = 0; i <= 10; i = i + 1 ) {
    if ( i % 2 == 0 ) {
        print("is even", i);
    } else {
        print("is odd", i);
    }
}
```