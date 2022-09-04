# Jo Programming Language

An interpreted toy programming language built with golang. [Sample program here.](#sample-program)

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
* [Declarations](#declarations)
    * [Variable Declaration](#variable-declaration)
    * [Function Declaration](#function-declaration)
    * [Struct Declaration](#struct-declaration)
* [Statements](#statements)
    * [If-Else Statement](#if-else-statement)
    * [For Loop Statement](#for-loop-statement)
* [Built-In Functions](#built-in-functions)
    * [Input](#input)
    * [Print](#print)
---
## Data Types
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

## Operators
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
| Operator |                       Meaning                       |                                Example                               |   |   |
|:--------:|:---------------------------------------------------:|:--------------------------------------------------------------------:|---|---|
| &&       | Logical AND. True only if all operands are true     | If c = 5 and d = 2 then, expression ((c==5) && (d>5)) equals to 0.   |   |   |
| \|\|     | Logical OR. True only if either one operand is true | If c = 5 and d = 2 then, expression ((c==5) \|\| (d>5)) equals to 1. |   |   |
| !        | Logical NOT. True only if the operand is 0          | If c = 5 then, expression !(c==5) equals to 0.                       |   |   |
---
## Declarations

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
```c
struct Person {
    fn walk(steps, limit) {
        for (let i = 0; i < steps; i = i + 1) {
            if (i >= limit) {
                break;
            }
            print("Person Walking", i + 1, "step");
        }
    }
    fn talk() {
        print("Person", "Talking");
    }
}

// Making instance of the struct
let p1 = Person();

// Calling methods of the Person struct instance
p1.walk(5, 10);
```
---
## Statements
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

### For Loop Statement
```js
for (let i = 0; i < 10; i = i + 1) {
    if (i == 5 ) {
        break; // continue
    }
    print(i);
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