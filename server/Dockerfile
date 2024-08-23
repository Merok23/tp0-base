FROM python:3.9.7-slim
COPY ./server /server
RUN python -m unittest server/tests/test_common.py
ENTRYPOINT ["/bin/sh"]
