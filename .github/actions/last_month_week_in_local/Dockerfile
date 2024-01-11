FROM python:3.10-bookworm
COPY requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY last_month_week_in_local.py /last_month_week_in_local.py
ENTRYPOINT ["python", "/last_month_week_in_local.py"]
