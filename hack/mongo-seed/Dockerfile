FROM mongo:5.0

ARG DATABASE_URI
ARG USERNAME
ARG PASSWORD

ENV DATABASE_URI ${DATABASE_URI}
ENV USERNAME ${USERNAME}
ENV PASSWORD ${PASSWORD}

COPY ./testdata/*.json /testdata
COPY seed-database.sh /
RUN chmod +x /seed-database.sh
RUN /seed-database.sh ${USERNAME} ${PASSWORD} ${DATABASE_URI}
