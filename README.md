<p align="center">
    <img src="./README.png"><br>
    Exported glucose measurements from your freestyle libre in your monitoring setup
<p>

## Notes

> [!CAUTION]
> Im pretty sure this is not official, so i wont take any warrenty of anything. This is experimental and just a fun project. Use at your own risk.

## Usage

### Prepare

Even if you run the Libre app on your phone, you also have to install the [LibreLinkUp](https://www.librelinkup.com/) app. This is the app for family or friends to receive glucose values. This is how this software gets its data from. It fetches it from that api.

1. So install the app and create an account
2. Link it to the Libre app

### Running the exporter

The easiest way to install and running the exporter is through docker. Please check the [release](https://github.com/xsteadfastx/glucose_exporter/releases) page for the latest release.

#### Configuration

| Environment variable | Description                                                                                  |
| -------------------- | -------------------------------------------------------------------------------------------- |
| `EMAIL`              | Email login data for LibreLinkUp                                                             |
| `PASSWORD`           | LibreLinkUp password. Consider using `PASSWORD_FILE`                                         |
| `PASSWORD_FILE`      | File with the account password in it. Nice to use in combination with docker compose secrets |
| `CACHE_DIR`          | Where to store cache data. Defaults to `/var/cache/glucose_exporter`                         |
| `DEBUG`              | Enabling debug logging                                                                       |

## Exported metrics

| Metrics              | Description                              |
| -------------------- | ---------------------------------------- |
| `value_in_mg_per_dl` | The glucose level itself                 |
| `trend_arrow`        | A integer representing the glucose trend |

### Trend interpretation

| Trend | Meaning |
| ----- | ------- |
| 1     | ⬇️      |
| 2     | ↘️      |
| 3     | ➡️      |
| 4     | ↗️      |
| 5     | ⬆️      |

## Thanks

- https://github.com/FokkeZB/libreview-unofficial
- https://gist.github.com/khskekec/6c13ba01b10d3018d816706a32ae8ab2
- https://libreview-unofficial.stoplight.io
