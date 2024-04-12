#!/bin/bash

# 设置开始时间
start_time=$(date +%s)

# 循环条件，当前时间减去开始时间小于10秒
while [ $(($(date +%s) - $start_time)) -lt 2 ]; do
    echo "Doing some work..."
    sleep 1 # 假设的任务，这里只是简单地等待1秒
done

# 循环结束后输出错误信息并以非零状态退出
echo "An error occurred." >&2
exit 1