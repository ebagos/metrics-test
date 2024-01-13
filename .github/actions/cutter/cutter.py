import re
import os
from dotenv import load_dotenv
from typing import List

def read_lines_from_file(filename: str) -> List[str]:
    """指定されたファイルから行を読み込み、リストとして返します。"""
    with open(filename, 'r') as file:
        lines = file.readlines()
    return lines

def write_lines_to_file(filename: str, lines: List[str], n: int, m: int) -> None:
    # 指定されたファイルに行を書き込む
    with open(filename, 'w') as file:
        # 書き込まない行を探す        
        empty_line_count = 0
        count = 0
        start_count = 0
        end_cout = 0
        for line in lines:
            if re.match(r'^\s*$', line):
                empty_line_count += 1
                if empty_line_count == n:
                    start_count = count
                if empty_line_count == m:
                    end_cout = count
            count += 1
        
        # 書き込む
        count = 0
        for line in lines:
            if count < start_count or count >= end_cout:    # 1つは空行を残す
                file.write(line)
            count += 1

def main():
    load_dotenv()
    # パラメータ
    n = int(os.getenv('START_BLANK_LINE'))
    m = int(os.getenv('END_BLANK_LINE'))
    input_file = os.getenv('INPUT_FILE')
    output_file = os.getenv('OUTPUT_FILE')

    # ファイル操作
    input_lines = read_lines_from_file(input_file)
    write_lines_to_file(output_file, input_lines, n, m)

if __name__ == '__main__':
    main()
