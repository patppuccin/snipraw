FROM scratch

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/snipraw /snipraw

EXPOSE 8245

ENTRYPOINT ["/snipraw"]
CMD ["--host", "0.0.0.0", "--dir", "/snippets"]