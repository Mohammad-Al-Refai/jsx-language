# ht

interpreted html syntax like programming language

## App

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

```html
<App>
   </Let id={"names"} value={[App.name,"A"]}>
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
   <For var={i} from={0} to={10}>
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
