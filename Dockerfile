FROM debian:stable-slim
COPY out /bin/out
ENV PORT 8000
CMD ["/bin/out"]