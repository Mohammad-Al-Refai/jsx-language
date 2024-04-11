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
   </Let id={name} value={"Mohammad"}x>
</App>
```

```html
<App>
   </Let id={number} value={1}>
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
   <If condition={1>2}>
    ...
    </If>
</App>
```

## Loop

```html
<App>
   <For var={i} from={0} to={10}>
     <Print value={For.i}/>
    </For>
</App>
```

## Function

```html
<App>
   <Function id={"Sum"} params={x,y}>
        <Return value={Sum.x+Sum.y}>
    </Function>
</App>
```

### Call function

```html
<App>
   </Print value={<App.Sum x={10} y={20}>}>
</App>
```
