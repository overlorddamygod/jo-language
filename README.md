# Jo Programming Language

An interpreted dynamic typed toy programming language built with golang. [Sample program here.](#sample-program)
More examples in [examples folder](./examples).

It is a mixture of many programming languages ( C, C++, Javascript, Python, etc). Syntax wise and Performance wise.

Language Grammar in [`jo-grammar.cfg`](./jo-grammar.cfg)

```
print("Hello World");
```
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
    * [Null](#null)
    * [Boolean](#boolean)
    * [Int](#int--float)
    * [Float](#int--float)
    * [String](#string)
    * [Array](#array)
    * [Struct](#struct)
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
        * [Switch Statement](#switch-statement)
    * [Looping Statements](#looping-statements)
        * [For Loop Statement](#for-loop-statement)
        * [While Loop Statement](#while-loop-statement)
    * [Import Export Statements](#import-export-statements)
* [Built-In Libraries](#built-in-functions)
    * [Input](#input)
    * [Print](#print)
    <details>
    <summary>Math Struct</summary>

    * [Math Struct](#math-struct)
        * [Attributes](#math-struct)
            * [pi](#math-struct)
            * [e](#math-struct)
        * [Methods](#math-struct)
            * [random](#math-struct)
            * [pow](#math-struct)
            * [exp](#math-struct)
            * [log](#math-struct)
            * [log2](#math-struct)
            * [log10](#math-struct)
            * [sqrt](#math-struct)
            * [abs](#math-struct)
            * [sin](#math-struct)
            * [cos](#math-struct)
            * [tan](#math-struct)
            * [round](#math-struct)
            * [ceil](#math-struct)
            * [floor](#math-struct)
    </details>
    
---
# Data Types
* All types have the following methods
    * len - length of the data ***works only for string and array***
    * type - returns the type of the data

### Null
Methods:
* type - returns the data type
* getInt - returns the integer value
* getFloat - returns the float value
* getString - returns the string value of the number

Example:
```js
let a = null;

print(a.type());
if (a) {
    print("Is not null");
}
```

### Boolean
Methods:
* type - returns the data type

Example:
```js
let a = true;

print(a.type());
if (a) {
    print("Is true");
}
```

### Int & Float

Methods:
* type - returns the data type
* getInt - returns the integer value
* getFloat - returns the float value
* getString - returns the string value of the number

Example:
```js
let numInt = 2;
let numFloat = 2.3;

print(numInt.type());
print(numFloat.type());

print(numInt.getFloat());
print(numFloat.getInt());

print(numInt.getString());
```
### String
Methods:
* type - returns the data type
* len - returns the string length
* get - get string from an index
* slice - get part of the string
* lower - convert string to lowercase letters.
* upper - convert string to uppercase letters.
* split - splits a string into an array of substrings.
* getInt - returns the integer value if parsable
* getFloat - returns the float value if parsable
* getString - returns the string value of the number
* replace - replace a string with new string

Example:
```js
let num = "29.2";

print(num);
print(num.len());
print(num.type());
print(num.getInt());
print(num.getFloat());

print("abc".upper());
print("ABCDEF".lower());

print("String Sliced:", num.slice(0, 3));
print("String Get from index:", num.get(2));
print("String Split:", num.split("."));

// replace method takes 2 or 3 args
print("hello world hello".replace("world", "earth")); // takes 2 args if we want to replace only 1 occurance
print("hello world hello".replace("hello", "earth", -1)); // takes 3 args if we want to replace `n` occurance. if n < 0 all the matched string are replaced
```
### Array
Methods:
* type - returns the data type
* len - returns the array length
* get - get data from an index
* slice -  get part of the array
* push -  push to the end of array
* pop -  pop from the end of array
* set -  set value to an index
* join -  join the contents of the array with a join string
* contains -  return bool if the array contains the data

Example:
```js
let arr = [1,2,3,4,5];
print("Initial Array:", arr);

print("Array Length:", arr.len());
print("Array Sliced:", arr.slice(1,5));
print("Array Get from index:", arr.get(2));

print("rray Push:", arr.push(6));
print("Array Pop:", arr.pop());

print("Array Set:", arr.set(2, 69));

print("Array Join:", arr.join(" "));
print("Array Contains:", arr.contains(2));

print("Final Array:", arr);
```

### Struct
Methods:
* type - returns the struct type

Example:
```js
struct Person {
    fn init(name) {
        print("constructor called");
        self.name = name;
    }
    fn walk() {
        print("walking");
    }
}

let p = Person("Nick"); // struct instance
print(p.type());
p.walk();
```
---

# Operators
### Arithmetic Operator
| Operator |             Meaning of Operator            |
|:--------:|:------------------------------------------:|
| +        | addition or unary plus                     |
| -        | subtraction or unary minus                 |
| *        | multiplication                             |
| /        | division                                   |
| %        | remainder after division (modulo division) |
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
| Operator |    Meaning of Operator   |
|:--------:|:------------------------:|
| ==       | Equal to                 |
| >        | Greater than             |
| <        | Less than                |
| !=       | Not equal to             |
| >=       | Greater than or equal to |
| <=       | Less than or equal to    |
### Logical Operator
| Operator |                       Meaning                       |
|:--------:|:---------------------------------------------------:|
| &&       | Logical AND. True only if all operands are true     |
| \|\|     | Logical OR. True only if either one operand is true |
| !        | Logical NOT. True only if the operand is 0          |
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
* Constructor method name should be `init`
```c
struct Person {
    // Constructor
    fn init(name, address, age) {
        self.name = name;
        self.address = address;
        self.age = age;
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

// Making instance of the struct
let p1 = Person("John", "USA", 20);

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
### Switch Statement
```js
let a = 3;
switch (a) {
    case 1:
        print("1");
        break;
    case 2, 3:
        print("2 or 3");
    default:
        print("NOT 1 or 2 or 3");
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
## Import Export Statements
```js
// car.js
struct Car {
    fn init(name) {
        self.name = name;
    }
    fn drive() {
        print("Driving a car named " + self.name);
    }
}
export Car;
```
```js
// main.js
import "car" as Car
// import "car"

let car1 = Car("lambo");
car1.drive();
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

### Math Struct
Has attributes and methods to perform mathematical tasks on numbers
```js
// Attributes
print(math.pi); // Value of Pi
print(math.e); // Value of Euler's number

// Methods
print(math.random()); // random float between 0 and 1
print(math.pow(2,2)); // x to the power y
print(math.exp(2)); // value of e ^ x

print(math.log(2)); // natural log
print(math.log2(2)); // log base 2
print(math.log10(2)); // log base 10

print(math.sqrt(4)); // square root
print(math.abs(-69)); // absolute (positive)

print(math.sin(6)); // sine value
print(math.cos(29)); // cosine value
print(math.tan(30)); // tan value

print(math.round(7.89)); // rounds to nearest integer
print(math.ceil(7.89)); // rounds up to nearest integer
print(math.floor(7.89)); // rounds down to nearest integer
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