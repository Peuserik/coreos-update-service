# CoreOs update service



Configure your Coreos hosts

```
GROUP=DevCluster
SERVER=http://xxxxxx:8000/v1/update/
REBOOT_STRATEGY=best-effort
LOCKSMITHD_REBOOT_WINDOW_START=Thu 23:00
LOCKSMITHD_REBOOT_WINDOW_LENGTH=1h30m
```

## Push the version info to the server

Use curl to push the version info to the server

```
curl -v -X PUT http://localhost:8000/version/1122.2.0  --data @version.json
```

1. With a version defined as

```
{
  "VersionId": "1122.2.0",
  "URL": "https://update.release.core-os.net/amd64-usr/1122.2.0/",
  "Name": "update.gz",
  "Hash": "+ZFmPWzv1OdfmKHaGSojbK5Xj3k=",
  "Signature": "cSBzKN0c6vKinrH0SdqUZSHlQtCa90vmeKC7p/xk19M=",
  "Size": 212555113
}


```

2. Push the tracks of the group servers

```
curl -v -X PUT http://localhost:8000/tracks  --data @tracks.json

```
