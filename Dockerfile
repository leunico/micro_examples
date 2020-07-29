FROM alpine
ADD myauth-api /myauth-api
ENTRYPOINT [ "/myauth-api" ]
