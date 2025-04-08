# 🔒 Auth Forwarding Image Server

🚀 A lightweight Go server that authenticates requests through a remote server and serves static images from the file system.

![Go Build](https://img.shields.io/badge/built%20with-Go-blue?style=flat-square)  
![Docker Image](https://img.shields.io/badge/dockerized-yes-brightgreen?style=flat-square)

- ✅ simple, secure access to local images based on external authentication 🔐  
- ✅ supports `.png`, `.jpg`, and `.jpeg` formats 🖼️  
- ✅ blazing fast with Go and zero dependencies 🏎️  
- ✅ customizable via environment variables 🛠️  

---

### 🛫 How to use it?

1. Place your image files into a directory mounted into the container (or accessible on host).  
2. Start the app with Docker:

```bash
docker run -p 8080:8080 \
-e AUTH_SERVER_BASE_URL="http://your-auth-server.local/auth" \
-v $(pwd)/images:/app/static \
ghcr.io/tobiasz-gleba/image-hosting-with-forwarded-auth
```

Now you can:
```http
GET http://localhost:8080/cats/image1.jpg
```

And it will:
- Forward request to: `http://your-auth-server.local/auth/cats/image1.jpg`
- If 200 OK → return `cats/image1.png` or `.jpg` or `.jpeg`
- If not → return 401 Unauthorized

---

### 🔧 Available Environment Variables

```env
AUTH_SERVER_BASE_URL=http://localhost:8081/auth
STATIC_DIR=/app/static
```

---

### 💡 Use Cases

- Private photo hosting  
- Auth-gated image previews  
- Secure image CDN  
