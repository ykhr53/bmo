FROM alpine
COPY bmo ./bmo
EXPOSE 3000
CMD ["./bmo"]
