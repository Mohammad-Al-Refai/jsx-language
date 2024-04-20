# Sogla

## Overview


Have you ever wondered how programming languages actually work under the hood? Do you believe that languages are just magical entities?


In this project, I set out to create a simple stack-oriented programming language and its syntax looks like JSX (because I love react js 😊) and all of that using the Go programming language with 0 dependencies. Inspired by my curiosity and desire to understand the inner workings of languages, this project served as an educational endeavor to demystify the magic behind programming languages, even it just an interpreted language not a compiled language but it worth it.

## Getting Started
Download sogla from the releases then create a file and run:

```bash
./sogla [file path]
```

### How it works?
![image](https://github.com/Mohammad-Al-Refai/ht/assets/55941955/8afca54b-c5d1-4e18-afd0-b16b339c9e4e)

### Comment
start with `#`

### Variables

Math operations: `+` & `-` & `*` & `/` & `%`

Support:

- `string`

- `number`

- `array`

- `boolean`

- object member call

</br>

```jsx
<App>
   <Let id={"name"} value={"Mohammad"}/>
<App/>

<App>
   <Let id={"number"} value={1}/>
<App/>

<App>
   <Let id={"data"} value={[1,"hi",false]}/>
<App/>

<App>
   <Let id={"data"} value={[1,"hi",false]}/>
   <Let id={"dataLength"} value={data array.length()}/>
<App/>
```

## If statement

Logical operators:

- `==`

- `!=`

- `greater` like `>`

- `smaller` like `<`

```jsx
<App>
   <If condition={1 2 greater}>
    ...
    <If/>
<App/>
```

## Loop

```jsx
<App>
   <For var={"i"} from={0} to={10}>
     <Print value={i}/>
    <For/>
<App/>
```

## Function

```jsx
<App>
   <Function id={"Sum"} args={"x","y"}>
        <Return value={x y +}>
    <Function/>
<App/>
```

## Function call

```jsx
<App>
  <Sum x={1} y={2}/>
<App/>
```

## Array

```jsx
// [array] array.length()

<App>
   <Let id={"data"} value={[1,2,"hello",false,423]}/>
   <Print value={data array.length()}/>
<App/>
```

```jsx
// [array] [index] array.at()

<App>
   <Let id={"data"} value={[1,2,"hello",false,423]}/>
   <Print value={0 data array.at()}/>
<App/>
```

```jsx
// [array] [value] array.push()

<App>
   <Let id={"data"} value={[]}>
   <Print value={"hello" data array.push()}/>
<App/>
```

```jsx
// [array] array.pop()

<App>
   <Let id={"data"} value={[1,2,3]}>
   <Print value={data array.pop()}/>
   <Print value={data}/>
<App/>

```

## Loop over array

```jsx
<App>
   <Let id={"data"} value={[1,2,"hello",false,423]}>
   <For var={"i"} from={0} to={data array.length() 1 -}>
      <Print value={i data array.at()}/>
   <For/>
<App/>
```

[More examples](https://github.com/Mohammad-Al-Refai/ht/tree/main/examples)
## Recurses I learnt from

- [Making Programming Language in Python -- Porth](https://www.youtube.com/watch?v=8QP2fDBIxjM&list=PLpM-Dvs8t0VbMZA7wW9aR3EtBqe2kinu4)

- [How To Build A Programming Language From Scratch](https://www.youtube.com/watch?v=8VB5TY1sIRo&list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh&pp=iAQB)
  
- [Stack orinted programming](https://en.wikipedia.org/wiki/Stack-oriented_programming)
