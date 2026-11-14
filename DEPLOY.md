# Stark Coder DDoS Tool - Public Deployment Package

## Build Instructions (For Developers)

### Prerequisites
- Go 1.21+ installed
- Node.js 18+ installed
- npm or bun

### Step 1: Install Dependencies
```bash
# Frontend dependencies
cd web-client
npm install
cd ..
```

### Step 2: Build Backend (Go)
```bash
go build -o bin/mmb-server.exe ./cmd/server
```

### Step 3: Build Frontend
```bash
cd web-client
npm run build
cd ..
Copy-Item -Path "web-client/dist/public/*" -Destination "bin/web-client" -Recurse -Force
```

### Step 4: Run Server
```bash
powershell -ExecutionPolicy Bypass -File run-server.ps1
```

## Public Deployment

### Option 1: VPS Deployment (Recommended)
1. Deploy to any VPS with public IP
2. Set up systemd service for auto-restart:
   ```ini
   [Unit]
   Description=Stark Coder Network Stresser
   After=network.target
   
   [Service]
   ExecStart=/opt/starkcoder/bin/mmb-server.exe
   WorkingDirectory=/opt/starkcoder
   Restart=always
   Environment="MMB_API_KEY=your_secret_key_here"
   
   [Install]
   WantedBy=multi-user.target
   ```

3. Setup reverse proxy with Nginx/Caddy for HTTPS
4. Open firewall port 3000

### Option 2: Cloudflare Tunnel (Free & Secure)
```bash
cloudflared tunnel --url http://localhost:3000
```
This gives you a free public URL without opening ports.

### Option 3: Render/Railway Hosting
1. Create `Procfile`:
   ```
   web: bin/mmb-server.exe
   ```
2. Push to GitHub
3. Connect to Render/Railway
4. Set environment variables

## Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| `MMB_API_KEY` | Optional | API key for authentication |
| `ALLOW_NO_PROXY` | No | Allow direct connections (true/false) |
| `PORT` | No | Custom port (default: 3000) |

## Security Notes
- Always set `MMB_API_KEY` in production
- Use HTTPS via reverse proxy
- Restrict access by IP if possible
- Never expose on public internet without authentication

## File Structure
```
stark-coder/
├── bin/                    # Compiled binaries
│   ├── mmb-server.exe     # Main server
│   └── web-client/        # Static assets
├── cmd/                    # Source code entry points
├── internal/              # Internal packages
├── pkg/                   # Public packages
├── data/                  # Configuration files
│   ├── proxies.txt       # Proxy list
│   └── uas.txt           # User agents
├── web-client/            # React frontend source
│   ├── src/
│   ├── package.json
│   └── dist/              # Built files
├── go.mod                 # Go module file
├── build.ps1             # Build script
├── run-server.ps1         # Run script
└── Makefile               # Build commands
```
