# IoTしてメトリクスを汎用的に集められるexporter

## 使い方

ラズパイ向け(arm64)にビルドするには次のようにする。

```console
$ GOARCH=arm64 GOOS=linux go build
```

systemdで管理できるようにするには、次のようにする。

- `temperature_exporter`という名前でユーザー作成
- 適切な場所にビルドしたバイナリを配置。
- unit file作成

具体的には、

```console
$ sudo useradd -r -s /sbin/nologin temperature_exporter
$ sudo cp temperature_exporter /usr/local/bin/temperature_exporter
$ sudo chown temperature_exporter:temperature_exporter /usr/local/bin/temperature_exporter
$ sudo chmod 744 /usr/local/bin/temperature_exporter
$ sudo vim /etc/systemd/system/temperature_exporter.service
```

`temperature_exporter.service`は次のように記述。

```systemd
[Unit]
Description=Temperature Exporter

[Service]
User=temperature_exporter
Group=temperature_exporter
ExecStart=/usr/local/bin/temperature_exporter

[Install]
WantedBy=multi-user.target
```

最後にサービスを読み込み&再起動

```console
$ sudo systemctl daemon-reload
$ sudo systemctl enable temperature_exporter
Created symlink /etc/systemd/system/multi-user.target.wants/temperature_exporter.service → /etc/systemd/system/temperature_exporter.service.
$ sudo systemctl start temperature_exporter
```

動いているか確認する。

```console
$ curl http://localhost:17818/metrics
```

Prometheus側では次のように設定する。

```yaml
scrape_configs:
  - job_name: 'temperature_exporter'
    static_configs:
      - targets: ['localhost:17818']
```