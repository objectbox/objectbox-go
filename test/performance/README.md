Performance tests
=================

To get good numbers, close all programs before running, build & run outside of IDE:

```bash
cd /tmp
go build /home/ivan/go/src/github.com/objectbox/objectbox-go/test/performance/
./performance
```

You can specify some parameters:
```bash
./performance -h

Usage of ./performance:
  -count int
    	number of objects (default 1000000)
  -db string
    	database directory (default "db")
  -runs int
    	number of times the tests should be executed (default 30)
```
