FROM circleci/golang:latest

MAINTAINER Shuhei Kitagawa <shuhei.kitagawa.noreply@gmail.com>

CMD go get github/shuheiktgw/github-label-checker
CMD go get github/shuheiktgw/bump-reviewer