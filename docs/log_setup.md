# 로그 설정 가이드

## 개요

이 프로젝트는 [logrus](https://github.com/sirupsen/logrus)를 로거로, [lumberjack](https://github.com/natefinch/lumberjack)을 파일 로테이션 라이브러리로 사용합니다.

로그 설정은 `configs/prod.yaml`의 `log` 섹션에서 관리합니다.

---

## 설정 옵션

### 기본 옵션

| 키       | 타입   | 기본값   | 설명                                                          |
| -------- | ------ | -------- | ------------------------------------------------------------- |
| `output` | string | `stdout` | 로그 출력 대상. `stdout`, `file`, `both` 중 하나              |
| `level`  | string | `debug`  | 로그 레벨. `debug`, `info`, `warn`, `error`, `silent`         |

### 파일 출력 옵션 (`output: "file"` 또는 `output: "both"` 일 때 유효)

| 키              | 타입   | 예시       | 설명                                                        |
| --------------- | ------ | ---------- | ----------------------------------------------------------- |
| `savePath`      | string | `logs`     | 로그 파일 디렉토리. 상대경로 시 `$APP_HOME` 기준            |
| `fileName`      | string | `app.log`  | 로그 파일 이름. 백업 시 `app-{timestamp}.log` 형식으로 저장 |
| `sizePerFileMb` | int    | `100`      | 파일 최대 크기 (MB). 초과 시 로테이션                        |
| `maxOfDay`      | int    | `10`       | 보관할 백업 파일 수                                          |
| `maxAge`        | int    | `7`        | 로그 파일 보관 기간 (일)                                     |
| `compress`      | bool   | `false`    | 오래된 백업 파일 gzip 압축 여부                             |

---

## 설정 예시

### stdout 출력

```yaml
log:
  output: "stdout"
  level: "info"
```

### 파일 출력

```yaml
log:
  output: "file"
  level: "info"
  savePath: "logs"        # $APP_HOME/logs/ 디렉토리에 저장
  fileName: "app.log"     # 실제 파일: $APP_HOME/logs/app.log
  sizePerFileMb: 100      # 100MB 초과 시 로테이션
  maxOfDay: 10            # 최대 10개 백업 파일 보관
  maxAge: 7               # 7일 이상 된 파일 삭제
  compress: true          # 백업 파일 gzip 압축
```

### stdout + 파일 동시 출력

```yaml
log:
  output: "both"
  level: "debug"
  savePath: "logs"
  fileName: "app.log"
  sizePerFileMb: 100
  maxOfDay: 10
  maxAge: 7
  compress: false
```

---

## 로그 레벨

| 레벨     | 설명                       |
| -------- | -------------------------- |
| `debug`  | 디버그 이상 모든 로그 출력 |
| `info`   | 정보 이상 로그 출력        |
| `warn`   | 경고 이상 로그 출력        |
| `error`  | 에러 이상 로그 출력        |
| `silent` | 모든 로그 억제             |

`LOG_LEVEL` 환경 변수로 설정 파일의 레벨을 덮어쓸 수 있습니다.

---

## 파일 로테이션 동작

lumberjack이 다음 조건에 따라 자동으로 로그 파일을 로테이션합니다:

- **크기 기반**: 파일 크기가 `sizePerFileMb`를 초과하면 현재 파일을 타임스탬프가 포함된 이름으로 백업하고 새 파일을 생성합니다.
- **백업 수 제한**: 백업 파일이 `maxOfDay` 개를 초과하면 오래된 파일부터 삭제합니다.
- **기간 제한**: `maxAge` 일이 지난 백업 파일을 삭제합니다.
- **압축**: `compress: true` 설정 시 백업 파일을 `.gz`로 압축합니다.

백업 파일 이름 형식: `app-2006-01-02T15-04-05.000.log`
