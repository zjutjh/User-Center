FROM scratch
MAINTAINER  cbluebird
WORKDIR /app
COPY user_center /app
EXPOSE 6060
CMD ["./user_center"]