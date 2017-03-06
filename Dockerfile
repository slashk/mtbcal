FROM golang:1.8

RUN go get -u github.com/gobuffalo/buffalo/buffalo
RUN curl https://glide.sh/get | sh

RUN curl -sL https://deb.nodesource.com/setup_6.x -o nodesource_setup.sh && \
    bash nodesource_setup.sh && \
    apt-get install -y nodejs build-essential

ADD . /go/src/github.com/slashk/mtbcal
RUN rm -rf /go/src/github.com/slashk/mtbcal/node_modules

WORKDIR /go/src/github.com/slashk/mtbcal

RUN npm install
RUN buffalo build -o bin/app

EXPOSE 3000:3000

CMD ["./docker-start.sh"]
