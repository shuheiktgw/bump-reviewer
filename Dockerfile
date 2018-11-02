FROM circleci/golang:latest

MAINTAINER Shuhei Kitagawa <shuhei.kitagawa.noreply@gmail.com>

RUN go get github.com/shuheiktgw/github-label-checker
RUN go get github.com/shuheiktgw/bump-reviewer