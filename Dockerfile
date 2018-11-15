FROM scratch

ADD ./bin/vk-callback-api /usr/bin/vk-callback-api

EXPOSE 9911

CMD ["/usr/bin/vk-callback-api"]
