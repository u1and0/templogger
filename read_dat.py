#!/usr/bin/env python3
from io import StringIO
import pandas as pd
import os

def read_dat(exe, *filenames, **kwargs):
    """templogger用datファイルの読み込み
    JSON中間ファイルを生成せずに
    DATファイルから直接pandasオブジェクトを返します
    """
    files = ' '.join(filenames)
    jsonbin = get_ipython().getoutput(
        f'{exe} --format json {files}')
        # !./templogger --format json {files} 実行するのと同じ
    jsonstr = StringIO(jsonbin[0])
    df = pd.read_json(jsonstr, **kwargs)
    df.Time = pd.to_datetime(df.Time)
    return df
