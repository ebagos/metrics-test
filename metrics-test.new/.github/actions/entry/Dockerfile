FROM python:3.10-bookworm
COPY requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY create.py /create.py
ENTRYPOINT ["python", "/create.py"]
