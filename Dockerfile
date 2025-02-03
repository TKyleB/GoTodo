FROM --platform=linux/amd64 debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates
ADD snippetz /usr/bin/snippetz

CMD ["snippetz"]