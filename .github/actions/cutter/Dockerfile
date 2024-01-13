FROM python:3.10-bookworm
COPY requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY cutter.py /cutter.py
ENTRYPOINT ["python", "/cutter.py"]
