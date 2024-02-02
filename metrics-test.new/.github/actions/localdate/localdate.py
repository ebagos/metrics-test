import os
from datetime import datetime, timedelta
from dotenv import load_dotenv
from dateutil.tz import gettz

def main():
    load_dotenv()
    tz = os.getenv('TIMEZONE')
    if tz == '':
        tz = 'Asia/Tokyo'
    type = os.getenv('TYPE')
    if type == '':
        type = 'week'
    weekday = os.getenv('WEEKDAY')
    if weekday == '':
        weekday = '0'
    first = None
    last = None
    current = os.getenv('UTC')
    utc = datetime.strptime(current, '%Y-%m-%dT%H:%M:%SZ')
    tz_local = gettz(tz)
    local = utc.astimezone(tz_local)
    # typeがweekの場合、先週の日曜日の00:00:00をfirstに設定し、先週の土曜日の23:59:59をlastに設定する
    if type == 'week':
        first = local - timedelta(7 + local.isoweekday() % 7)
        if weekday != '7':
            first = first + timedelta(days=int(weekday))
        first = datetime(year=first.year, month=first.month, day=first.day, hour=0, minute=0, second=0, tzinfo=tz_local)
        last = first + timedelta(days=6)
        last = datetime(year=last.year, month=last.month, day=last.day, hour=23, minute=59, second=59, tzinfo=tz_local)
    # typeがmonthの場合、先月の1日の00:00:00をfirstに設定し、先月の最終日の23:59:59をlastに設定する
    elif type == 'month':
        # 1ヶ月前の1日を取得
        if local.month == 1:
            first = datetime(year=local.year-1, month=12, day=1, hour=0, minute=0, second=0, tzinfo=tz_local)
        else:
            first = datetime(year=local.year, month=local.month-1, day=1, hour=0, minute=0, second=0, tzinfo=tz_local)
        # localの1ヶ月前の最終日を取得
        last = datetime(year=local.year, month=local.month, day=1, hour=23, minute=59, second=59, tzinfo=tz_local)
        last = last - timedelta(days=1)
        # 上記以外はエラーを出力する
    else:
        print('Error: TYPE is invalid')
        return

    # firstとlastをUTCに変換する
    first = first.astimezone(gettz('utc'))
    last = last.astimezone(gettz('utc'))
    # firstとlastを出力する
    # print('first: ' + first.strftime('%Y-%m-%dT%H:%M:%SZ'))
    # print('last: ' + last.strftime('%Y-%m-%dT%H:%M:%SZ'))
    output_str = '\nfirst=' + first.strftime('%Y-%m-%dT%H:%M:%SZ') + '\nlast=' + last.strftime('%Y-%m-%dT%H:%M:%SZ') + '\n'
    print(output_str)
    output = os.getenv('GITHUB_OUTPUT')
    with open(output, mode='a') as f:
        f.write(output_str)

main()
