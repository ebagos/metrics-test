FROM python:3.10-bookworm
COPY requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY localdate.py /localdate.py
ENTRYPOINT ["python", "/localdate.py"]
