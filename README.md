# wasm-image-processor

> 경희대학교 마이크로서비스프로그래밍 과목 중 수행하는 프로젝트를 위한 Repository입니다.

Rust로 작성한 이미지 처리 프로그램을 **WebAssembly(WASM)** 로 컴파일해 웹 브라우저 환경에서 실행해보는 간단한 실습 프로젝트입니다.

<br>

- **Rust → WASM** 변환을 통한 웹 실행 테스트  
- 브라우저에서 직접 동작하는 이미지 처리 기능  
- **Docker + Nginx** 로 웹 환경 구성  
- **멀티플랫폼 빌드 지원** (AMD64 / ARM64)

<br>

---

## 실행 방법

```bash
# Docker Hub에서 이미지 다운로드
docker pull mingdodev/wasm-image-processor:1.0

# 컨테이너 실행
docker run -p [호스트포트]:80 mingdodev/wasm-image-processor:1.0
```
브라우저에서 `http://localhost:설정한호스트포트` 로 접속