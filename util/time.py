import datetime


def current_datetime():
    return datetime.datetime.now(tz=datetime.timezone.utc)


def parse_datetime(dt_str: str) -> datetime.datetime:
    return datetime.datetime.fromisoformat(dt_str).astimezone(datetime.timezone.utc)
