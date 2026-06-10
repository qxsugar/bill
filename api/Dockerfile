FROM ccr.ccs.tencentyun.com/wwww/app:alpine-tz-wk

WORKDIR /app

ARG build_message="hello world"
ARG build_version="0.0.1"

RUN echo "$build_message" > /app/build_message
RUN echo "$build_version" > /app/build_version

COPY app /app/
COPY docs /app/docs

EXPOSE 80
CMD ["./app", "api"]
