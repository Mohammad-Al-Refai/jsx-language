# ht

interpreted html syntax like programming language

## App (entry point)

```html
<App>
...
</App>
```

## Variables

```html
<App>
   </Let id={"name"} value={"Mohammad"}>
</App>
```

```html
<App>
   </Let id={"number"} value={1}>
</App>
```

## If

```html
<App>
   <If condition={1 2 greater}>
    ...
    </If>
</App>
```

## Loop

```html
<App>
   <For var={"i"} from={0} to={10}>
     </Print value={i}>
    </For>
</App>
```

## Function

```html
<App>
   <Function id={"Sum"} args={"x","y"}>
        <Return value={x y +}>
    </Function>
</App>
```

### Call function

```html
<App>
  </Sum x={1} y={2}>
</App>
```

## Array

```jsx
[array] array.length()

# example
<App>
   </Let id={"data"} value={[1,2,"hello",false,423]}>
   </Print value={data array.length()}>
</App>

# output:
7
```

```js
[array] [index] array.at()

# example
<App>
   </Let id={"data"} value={[1,2,"hello",false,423]}>
   </Print value={0 data array.at()}>
</App>

# output:
1
```

```js
[array] [value] array.push()

# example
<App>
   </Let id={"data"} value={[]}>
   </Print value={"hello" data array.push()}>
</App>

# output:
hello
```

```js
[array] array.pop()

# example
<App>
   </Let id={"data"} value={[1,2,3]}>
   </Print value={data array.pop()}>
   </Print value={data}>
</App>

# output:
3
[1 2]
```

### Loop over array

```js

# example
<App>
   </Let id={"data"} value={[1,2,"hello",false,423]}>
   <For var={"i"} from={0} to={data array.length() 1 -}>
      </Print value={i data array.at()}>
   </For>
</App>

# output:
1
2
hello
false
423
```
