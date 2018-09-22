FROM golang:1.8.5-jessie

# if left blank app will run with dev settings
# to build production image run:
# $ docker build ./api --build-args app_env=production
ARG app_env
ENV APP_ENV $app_env

# install dep
RUN go get github.com/golang/dep/cmd/dep
RUN go get github.com/lib/pq
RUN go get github.com/dgrijalva/jwt-go
RUN go get golang.org/x/crypto/bcrypt

# it is okay to leave user/GoDoRP as long as you do not want to share code with other libraries
COPY . /go/src/github.com/atmiguel/cerealnotes
WORKDIR /go/src/github.com/atmiguel/cerealnotes

# install packages
# --vendor-only is used to restrict dep from scanning source code
# and finding dependencies
RUN dep ensure --vendor-only


# if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
# if production setting will build binary
# CMD if [ ${APP_ENV} = production ]; \
#	then \
#	api; \
#	else \
#	go get github.com/pilu/fresh && \
#	fresh; \
#	fi

CMD ["go", "run", "main.go"]

EXPOSE 8080