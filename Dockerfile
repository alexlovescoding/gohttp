FROM scratch

COPY public-html /public-html
COPY bin/main /

VOLUME /public-html

EXPOSE 80
EXPOSE 443

CMD ["/main"]
