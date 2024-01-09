FROM python:3.10-bookworm
COPY requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY commits.py /commits.py
ENTRYPOINT ["python", "/commits.py"]
