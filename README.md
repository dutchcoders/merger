# Merger

Simple merge source object into destination.

```
if err := merger.Merge(&dest, src); err != nil {
    panic(err)
}
```

## Todo's

* implement Merger interface, for custom merging strategies
* improve error handling
* add error strategy (continue, fail)
