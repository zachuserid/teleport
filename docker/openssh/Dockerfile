FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
    rsyslog \
    vim \
    iproute2 \
    net-tools \
    openssh-server
RUN mkdir /var/run/sshd
RUN echo 'root:screencast' | chpasswd

# sshd_config
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/LogLevel INFO/LogLevel DEBUG3/' /etc/ssh/sshd_config
RUN echo "TrustedUserCAKeys /etc/ssh/tca.pub" >> /etc/ssh/sshd_config
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGlMvOoWd4pRYPoE+G/jgTkgJUeMBkfQeYopzj8+UjdO71ueKYqyqUfoUc4iquQpSJsZig8/Vfwqv9h0+mOnKvCYfNs+RnkLdk3m8unn/ybS/z4g6Y8JvwTSzQV9vGHxk/TbceGqhU+81Dd7vn5mlsp6NFW/CH0oNUWoQOdpF38DOEirvbgDG6zc6maocLU/qaX/bPfl/c0xRZYA3r2QF8/HhbUvHFWfvnT4Il/mGxup3W6OKWTAncSq+mh8aLD4b/H/3vcMk7HwknX5A3GOF6rLprT+3s7DL/MVe2sThqk0raBvCxIAkQetUf7mt/FQ0PX1jLoEUQy/MBwvEYvHar" >> /etc/ssh/tca.pub

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D", "-E", "/ssh.log"]
