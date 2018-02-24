FROM alpine

RUN apk add --no-cache mongodb && \
rm /usr/bin/mongoperf

VOLUME /data/db
EXPOSE 27017

WORKDIR /src
COPY mongorun.sh /src
RUN chmod 755 /src/mongorun.sh
ENTRYPOINT [ "/src/mongorun.sh" ]
CMD [ "mongod", "--bind_ip", "0.0.0.0" ]