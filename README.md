# wasm-image-processor

> 경희대학교 마이크로서비스프로그래밍 과목 중 수행하는 프로젝트를 위한 Repository입니다.

<br>

Rust로 작성한 이미지 처리 프로그램을 **WebAssembly(WASM)** 로 컴파일해 웹 브라우저 환경에서 실행해보는 간단한 실습 프로젝트입니다.

- **Rust → WASM** 변환을 통한 웹 실행 테스트  
- 브라우저에서 직접 동작하는 이미지 처리 기능  
- **Docker + Nginx** 기반 웹 서비스 구성  
- **Docker Compose / Kubernetes(minikube)** 를 이용한 배포 실습
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

<br>

## 프로젝트 3

Docker Compose 기반의 프로젝트 2를 **Kubernetes 환경**으로 확장한 버전입니다.
**Deployment, Service(NodePort), PersistentVolumeClaim** 등을 활용해 **Pod를 관리하고 스케일링 가능한 구조**를 구성했습니다.

웹 클라이언트는 이미지 업로드 요청을 API 서버로 전송하고, API 서버는 업로드된 이미지를 PVC에 저장합니다.

Web Pod는 웹 클라이언트, API Pod는 서버 로직을 담당하며, Kubernetes Service를 통해 서로 통신합니다.

<br>

### 배포 방법

```bash
# Kubernetes 리소스 생성
kubectl apply -f k8s.yml
```

<br>

### 웹 클라이언트 접속

Web Pod는 NodePort로 외부에 노출되므로, 아래 명령어를 통해 로컬에서 접속 가능한 URL을 확인합니다.

```bash
minikube service web --url
```
- 출력된 URL로 브라우저에서 접속하면 웹 클라이언트가 실행됩니다.  
- 예시: `http://127.0.0.1:59866`

<br>

### 이미지 저장 (Persistent Volume)

API Pod는 이미지 파일을 <strong>PersistentVolumeClaim(PVC)</strong>에 저장합니다.  
PVC는 파드의 생명주기와 독립적이므로, 파드가 재시작되더라도 저장된 이미지는 유지됩니다.

<br>

### API 요청 처리

웹 클라이언트는 `localhost:8080`으로 API 요청을 보내도록 구현되어 있기 때문에,
다음 명령어로 호스트의 8080 포트를 API Service와 연결해 API 요청이 가능하도록 설정합니다:

```bash
kubectl port-forward svc/api 8080:8080
```
이를 통해 브라우저에서 업로드된 이미지를 API Pod로 전달하고, Kubernetes 클러스터 내부에서 저장 및 조회할 수 있습니다.

<br>

### 로그 확인

두 개 이상의 API 파드를 실행한 경우, 다음 명령어로 파드별 처리 로그를 확인할 수 있습니다:

```bash
kubectl logs -f -l app=api --prefix
```