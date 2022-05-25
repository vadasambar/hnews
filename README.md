## HNews
`HNews` is a Kubernetes Custom Resource you can use to query HNews articles that you want using filters.

Example:
```yaml
apiVersion: apps.vadasambar.com/v1
kind: HNews
metadata:
  name: hnews-sample
spec:
  filter:
    score: ">300"
    limit: 6
    descendents: ">10"
```
Result:
```yaml
apiVersion: apps.vadasambar.com/v1
kind: HNews
metadata:
  ...
  generation: 3
  name: hnews-sample
  namespace: default
  ...
spec:
  filter:
    descendents: '>10'
    limit: 6
    score: '>300'
    type: story
status:
  link:
  - article_url: https://github.com/dzhang314/YouTubeDrive
    descendents: 260
    hnews_url: https://news.ycombinator.com/item?id=31495049
    score: 607
  - article_url: ""
    descendents: 376
    hnews_url: https://news.ycombinator.com/item?id=31494417
    score: 453
  - article_url: https://e360.yale.edu/digest/bugs-are-evolving-to-eat-plastic-study-finds
    descendents: 221
    hnews_url: https://news.ycombinator.com/item?id=31495836
    score: 321
  - article_url: https://github.com/SymbianSource
    descendents: 149
    hnews_url: https://news.ycombinator.com/item?id=31491744
    score: 379
  - article_url: https://www.catphones.com/en-us/cat-s22-flip/
    descendents: 270
    hnews_url: https://news.ycombinator.com/item?id=31493138
    score: 518
  - article_url: https://gweb-research-imagen.appspot.com/
    descendents: 610
    hnews_url: https://news.ycombinator.com/item?id=31484562
    score: 940
```
### How to use filter?
```yaml
$ kubectl explain hnews.spec.filter
KIND:     HNews
VERSION:  apps.vadasambar.com/v1

RESOURCE: filter <Object>

DESCRIPTION:
     Filter allows you to filter and get the Hacker News articles you want

FIELDS:
   descendents  <string> -required-
     Number of direct (first level) comments in the article. Specify it like: descendents:
     ">=10", descendents: "<10", descendents: "=10", descendents: "!=10"

   limit        <integer> -required-
     Number of Hacker News articles you want.

   score        <string> -required-
     Score of Hacker News articles you are looking for. Specify it like: score:
     ">=10", score: "<10", score: "=10", score: "!=10"

   type <string>
     Type of Hacker News articles you are looking for. Has to be either of:
     job,story,comment,poll,pollopt
```