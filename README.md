# PathPick

PathPick은 파일 시스템을 탐색하고 선택한 파일들의 내용을 클립보드로 복사하는 CLI 도구입니다.
(만든 이유는 claude 토큰 사용량을 줄이기 위해, 대화가 길어졌을 때 끊고 새 대화에서 시작하기 위함.)

## 설치

```bash
# 빌드 및 설치
git clone https://github.com/cjnghn/pathpick
cd pathpick
go build -o pathpick ./cmd/pathpick
sudo cp pathpick /usr/local/bin/
sudo chmod +x /usr/local/bin/pathpick
```

## 사용법

```bash
# 현재 디렉토리에서 시작
pathpick

# 특정 디렉토리에서 시작
pathpick -path ~/projects

# 특정 패턴의 파일만 표시
pathpick -pattern "*.go"
```

## 키 조작법

- `↑`/`↓`: 항목 이동
- `←`/`→`: 상위/하위 디렉토리 이동
- `Space`: 항목 선택/해제
- `y`: 선택한 파일 내용을 클립보드에 복사
- `q` or `ESC`: 종료

## 특징

- tree 스타일의 파일 시스템 탐색
- 디렉토리 선택 시 하위 항목 모두 선택
- 선택한 파일들의 내용을 클립보드에 복사
- 파일 패턴 필터링 지원
