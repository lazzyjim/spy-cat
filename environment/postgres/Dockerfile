FROM postgres:11

ARG UID

ADD schema.sql /
ADD initialize_db.sh /
RUN chmod +x /initialize_db.sh

# make user id the same inside container as on hostmachine
RUN usermod -u "$UID" postgres
