# SELPG

**Select Page([Selpg](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html))** is a terminal program, which helps print files with given pages. Thr origin author develops it with `C` but for this program, implement it with **`GO`**.

## Usage

```bash
selpg:
  --d string   name of the printer (default "default")
  --e int      End page of file (default -1)
  --f          flag splits page
  --l int      lines in one page (default 72)
  --s int      Start page of file (default -1)
```

## Develop Tutorial

[Implement SELPG with GO](https://xwy27.github.io/Service-Computing/Selpg-WithGO/)

## Test

1. Simple print task in CLI test pass.
1. Without printer, the -d option returns ERROR msg in my test.
