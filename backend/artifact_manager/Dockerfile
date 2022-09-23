FROM registry.access.redhat.com/ubi8/python-38
WORKDIR /opt/app-root/src
COPY backend/artifact_manager/requirements.txt /opt/app-root/src
RUN pip install --upgrade pip && pip install -Ur requirements.txt
CMD ["/bin/bash", "-e"]
