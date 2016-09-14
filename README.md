# Merger

Merge source object into destination. This is especially useful when converting from for example database models to json format.

```
if err := merger.Merge(&dest, src); err != nil {
    panic(err)
}
```

## Todo's

* implement Merger interface, for custom merging strategies
* improve error handling
* add error strategy (continue, fail)

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.

Code released under [the MIT license](LICENSE).
