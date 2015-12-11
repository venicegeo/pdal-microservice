FROM venicegeo/minimal-pdal
EXPOSE 8080
WORKDIR /app
# copy binary into image
COPY pdal /app/
ENTRYPOINT ["./pdal"]
