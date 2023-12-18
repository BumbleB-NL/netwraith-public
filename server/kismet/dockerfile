ARG ARCH=
FROM ${ARCH}debian:bullseye
ENV DEBIAN_FRONTEND=noninteractive

RUN apt update -y
RUN apt install ca-certificates wget gpg -y
RUN wget -O - https://www.kismetwireless.net/repos/kismet-release.gpg.key --quiet | gpg --dearmor | tee /usr/share/keyrings/kismet-archive-keyring.gpg
RUN echo 'deb [signed-by=/usr/share/keyrings/kismet-archive-keyring.gpg] https://www.kismetwireless.net/repos/apt/release/bullseye bullseye main' | tee /etc/apt/sources.list.d/kismet.list
RUN apt update -y
RUN apt install kismet -y
RUN useradd -r netwraith
RUN usermod -aG kismet netwraith
RUN mkdir /home/netwraith
RUN chown netwraith:netwraith /home/netwraith

USER netwraith
RUN mkdir -p /home/netwraith/.kismet/logs
RUN mkdir /home/netwraith/.kismet/plugins

EXPOSE 2501

ENTRYPOINT ["kismet"]
