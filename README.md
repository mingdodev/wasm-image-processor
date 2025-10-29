# wasm-image-processor

> 경희대학교 마이크로서비스프로그래밍 과목 중 수행하는 프로젝트를 위한 Repository입니다.

<br>

Rust로 작성한 이미지 처리 프로그램을 **WebAssembly(WASM)** 로 컴파일해 웹 브라우저 환경에서 실행해보는 간단한 실습 프로젝트입니다.

- **Rust → WASM** 변환을 통한 웹 실행 테스트  
- 브라우저에서 직접 동작하는 이미지 처리 기능  
- **Docker + Nginx** 로 웹 환경 구성  
- **멀티플랫폼 빌드 지원** (AMD64 / ARM64)

<br>

---

## 프로젝트 1

웹 브라우저에서 실행되는 Rust 이미지 처리 프로그램을 사용해볼 수 있습니다.

```bash
# Docker Hub에서 이미지 다운로드
docker pull mingdodev/wasm-image-processor:1.0

# 컨테이너 실행
docker run -p [호스트포트]:80 mingdodev/wasm-image-processor:1.0
```
- 브라우저에서 `http://localhost:설정한호스트포트` 로 접속

<br>

## 프로젝트 2

API 요청을 통해 이미지 파일을 서버에 저장하고, 저장된 이미지들을 메인 화면에서 확인할 수 있습니다.
서버는 API 핸들러와 정적 리소스 핸들러를 포함합니다.

```bash
# docker-compose.yml이 존재하는 위치에서 실행 (백그라운드 실행)
docker compose up -d
```
- 도커 컴포즈 파일이 있는 위치 기준 `./api/data/images` 경로에 이미지를 저장
- 브라우저에서 `http://localhost:8081` 로 접속