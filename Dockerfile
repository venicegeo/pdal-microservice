FROM cflinusfs2-with-pdal2

ENTRYPOINT ["/pzsvc-pdal"]

COPY pzsvc-pdal /pzsvc-pdal
RUN chmod a+x /pzsvc-pdal
