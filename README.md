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
```
$ kubectl get hnews
NAME           TYPE    SCORE   LIMIT   DESCENDENTS   LASTSYNCEDAT
hnews-sample   story   >300    6       >10           2022-05-26T04:28:09Z
```
```yaml
apiVersion: apps.vadasambar.com/v1
kind: HNews
metadata:
  ...
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
  lastSyncedAt: "2022-05-26T04:29:03Z"
  link:
  - article_url: https://www.ftc.gov/business-guidance/blog/2022/05/twitter-pay-150-million-penalty-allegedly-breaking-its-privacy-promises-again
    descendents: 247
    hnews_url: https://news.ycombinator.com/item?id=31510865
    score: 904
  - article_url: ""
    descendents: 1640
    hnews_url: https://news.ycombinator.com/item?id=31503201
    score: 742
  - article_url: ""
    descendents: 144
    hnews_url: https://news.ycombinator.com/item?id=31508009
    score: 456
  - article_url: https://uxdesign.cc/the-forgotten-benefits-of-low-tech-user-interfaces-57fdbb6ac83
    descendents: 385
    hnews_url: https://news.ycombinator.com/item?id=31502193
    score: 365
  - article_url: https://github.com/SymbianSource
    descendents: 186
    hnews_url: https://news.ycombinator.com/item?id=31491744
    score: 428
  - article_url: https://www.theglobeandmail.com/canada/article-the-great-junk-transfer-inheritance-decluttering-canada/
    descendents: 419
    hnews_url: https://news.ycombinator.com/item?id=31499766
    score: 443
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

# To run it locally
1. Install the CRDs first:
```
make install
```
2. And then run the controller
```
make run
```
That's all. Make sure your kube context is pointing to the correct cluster.
3. To uninstall, just do
```
make uninstall
# Ctrl+c on make run to stop the controller
```
4. To run tests,
```
make test
```

# To run it in your cluster (untested)
1. Build the image
```
make docker-build IMG=hnews-controller:latest
```
2. Push the image
```
make docker-push IMG=hnews-controller:latest # IMG should be same as the one used in `make docker-build`
```
3. Deploy the manifests
```
make deploy IMG=hnews-controller:latest # IMG should be the same as one used in `make docker-build` amd `make docker-push`
```