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
let name = input("Enter your name: ");

let hasName = name != "";

if (!hasName) {
    name = "World";
}

// Printing Hello, World if name is empty
// else Hello, {name} 
print("Hello,", name);


// Print if number is even or odd from 0 to 10

// Helper Functions
fn mod(num, by) {
    return num % by;
}

fn isEven(num) {
    return mod(num, 2) == 0;
}

for ( let i = 0; i <= 10; i = i + 1 ) {
    if ( isEven(i) ) {
        print("is even", i);
    } else {
        print("is odd", i);
    }
}

// Fibonacci
fn fib(num) {
    if (num == 0 || num == 1) {
        return 1;
    }
    return fib(num - 1) + fib(num - 2);
}

for ( let i = 0; i < 10; i = i + 1 ) {
    print("FIB", i, fib(i));
}
```