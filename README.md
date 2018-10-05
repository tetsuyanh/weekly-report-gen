# weekly-report-gen

generate your weekly report which integrated web services.

## web services

- github.com
  - issues
  - Pull requests
- asana
  - tasks
- esa

## grouping

group items by specifying multiple keywords.

the head of keyword become group label.

for exmaple, you specify [渋谷,Shibuya,shibuya]

representation like this
```
- 渋谷
  - /shibuya-ward/issue
  - /東京/渋谷/parks
  - /Shibuya-story/2018.10.05_dialy
```

## format

now, markdown only

exsample
```
## this week
- GroupA
  - [groupA/note](https://hoge.com)

```
