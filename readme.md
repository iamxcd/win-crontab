# win-crontab
    在win系统下crontab的解决方案

# 使用说明
```
// cron.json
[
  {
    "name": "打开百度",
    "crontab": "* * * * * *",
    "cmd": "start www.baidu.com"
  },
  {
    "name": "执行一个php文件",
    "cron": "* * * * * *",
    "cmd": "cd F:\\www\\test && php ping.php"
  }
]
```