# Jo Programming Language

An interpreted dynamic typed toy programming language built with golang. [Sample program here.](#sample-program)
More examples in [examples folder](./examples).

It is a mixture of many programming languages ( C, C++, Javascript, Python, etc). Syntax wise and Performance wise.

Language Grammar in [`jo-grammar.cfg`](./jo-grammar.cfg)

### References:
* [Crafting Interpreters](https://craftinginterpreters.com/appendix-i.html)
    
---
### Usage
Requires Golang

**Without Building:**
```console
user@main:~$ go run cmd/interpreter/main.go example.jo
```

**By building**
```console
// Build for mac
user@main:~$ GOOS=darwin GOARCH=arm64 go build -o jo cmd/interpreter/main.go

// Build for windows
user@main:~$ GOOS=windows GOARCH=amd64 go build -o jo cmd/interpreter/main.go

// Run
user@main:~$ ./jo example.jo
```
---
### Build
Builds for all OS and architecture as in [Makefile](./Makefile).
Built executables found at [./bin folder](./bin)
```console
user@main:~$ make
```
---
# JoLang Docs
* [Data Types](#data-types)
    * [Boolean](#data-types)
    * [Int](#data-types)
    * [Float](#data-types)
    * [String](#data-types)
    * [Struct](#data-types)
* [Operators](#operators)
    * [Arithmetic Operator](#arithmetic-operator)
    * [Assignment Operator](#assignment-operator)
    * [Relational Operator](#relational-operator)
    * [Logical Operator](#logical-operator)
* [Statements](#statements)
    * [Declaration Statements](#declaration-statements)
        * [Variable Declaration](#variable-declaration)
        * [Function Declaration](#function-declaration)
        * [Struct Declaration](#struct-declaration)
    * [Conditional Statements](#conditional-statements)
        * [If-Else Statement](#if-else-statement)
    * [Looping Statements](#looping-statements)
        * [For Loop Statement](#for-loop-statement)
        * [While Loop Statement](#while-loop-statement)
* [Built-In Functions](#built-in-functions)
    * [Input](#input)
    * [Print](#print)
---
# Data Types
```js
let a = true; // boolean
let a = 2; // int
let b = 2.2; // float
let name = "Jo" // string

struct Person {
    fn walk() {
        print("walking");
    }
}

let p = Person(); // struct instance

// Coming Soon: array and map [maybe lol]
```

---

# Operators
### Arithmetic Operator
| Operator |             Meaning of Operator            | Same as |
|:--------:|:------------------------------------------:|:-------:|
| +        | addition or unary plus                     | a = b   |
| -        | subtraction or unary minus                 | a = a+b |
| *        | multiplication                             | a = a-b |
| /        | division                                   | a = a*b |
| %        | remainder after division (modulo division) | a = a/b |
### Assignment Operator
| Operator | Example |
|:--------:|:-------:|
| =        | a = 2+3;|
| +=        | a += 6; |
| -=        | a -= 9; |
| *=        | a *= 4; |
| /=        | a /= 2; |
| %=        | a %= 0; |
| &&=        | a &&= false;|
| \|\|=        | a \|\|= true;|
| \|=        | a \|= 3;|
| &=        | a &= 6;|

### Relational Operator
| Operator |    Meaning of Operator   |          Example         |
|:--------:|:------------------------:|:------------------------:|
| ==       | Equal to                 | 5 == 3 is evaluated to 0 |
| >        | Greater than             | 5 > 3 is evaluated to 1  |
| <        | Less than                | 5 < 3 is evaluated to 0  |
| !=       | Not equal to             | 5 != 3 is evaluated to 1 |
| >=       | Greater than or equal to | 5 >= 3 is evaluated to 1 |
| <=       | Less than or equal to    | 5 <= 3 is evaluated to 0 |
### Logical Operator
| Operator |                       Meaning                       |                                Example                               |
|:--------:|:---------------------------------------------------:|:--------------------------------------------------------------------:|
| &&       | Logical AND. True only if all operands are true     | If c = 5 and d = 2 then, expression ((c==5) && (d>5)) equals to 0.   |
| \|\|     | Logical OR. True only if either one operand is true | If c = 5 and d = 2 then, expression ((c==5) \|\| (d>5)) equals to 1. |
| !        | Logical NOT. True only if the operand is 0          | If c = 5 then, expression !(c==5) equals to 0.                       |
---
# Statements
## Declaration Statements

### Variable Declaration

```js
let variable_name = 2 + 3;
```

### Function Declaration
```js
fn add(a, b) {
    return a + b;
}
```

### Struct Declaration
* Struct uses `self` keyword inside the struct methods/functions to assign attributes to the struct. [ ***like `this` keyword*** ] 
```c

struct Person {
    // Doesnt have constructor for now in the language. So we create a function to initialize and return the same struct object
    fn init(name, address, age) {
        self.name = name;
        self.address = address;
        self.age = age;
        return self;
    }
    fn printinfo() {
        print("Name", self.getName());
        print("Address", self.getAddress());
        print("Age", self.getAge());
    }
    fn walk(steps, limit) {
        for (let i = 0; i < steps; i = i + 1) {
            if (i >= limit) {
                break;
            }
            print("Person Walking", i + 1, "step");
        }
    }
    fn getName() {
        return self.name;
    }
    fn getAddress() {
        return self.address;
    }
    fn getAge() {
        return self.age;
    }
}

// Making // Making instance of the struct and calling a fucntion as constructor
let p1 = Person().init("John", "USA", 20);
// Calling methods of the Person struct instance

p1.printinfo();
p1.walk(5, 10);
```
---
## Conditional Statements
### If-Else Statement
```js
let a = 3;
if ( a == 0 ) {
    print(a, "is zero");
} elif ( a % 2 == 0 ) {
    print(a, "is even");
} else {
    print(a, "is odd");
}
```
---
## Looping Statements
### For Loop Statement
```js
for (let i = 0; i < 10; i += 1) {
    if (i == 5 ) {
        break; // continue
    }
    print(i);
}
```
### While Loop Statement
```js
let i = 0;
while (i < 10) {
    print(i);
    i += 1;
}
```
---
## Built-In Functions
### Input
Get input from the console.
Accepts a single argument.
```
let name = input("Enter your name:");
```

### Print
Print to the console.
Accepts any number of arguments.
```
let name = input("Enter your name:");
print("Your name is", name);
```
---
## Sample Program
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

for ( let i = 0; i <= 10; i += 1 ) {
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

for ( let i = 0; i < 10; i += 1 ) {
    print("FIB", i, fib(i));
}
```