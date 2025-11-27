# Production Deployment Guide

This guide covers deploying the CyclingStream platform to production.

## Prerequisites

- Server with Docker and Docker Compose installed
- Domain name with DNS configured
- SSL certificate (Let's Encrypt recommended)
- PostgreSQL database (managed or self-hosted)
- Stripe account with API keys
- CDN account (BunnyCDN or similar) for HLS streaming

## Server Requirements

### Minimum Specifications
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 50GB SSD
- **Network**: 100Mbps

### Recommended Specifications
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 100GB+ SSD
- **Network**: 1Gbps

## Step 1: Server Setup

### 1.1 Initial Server Configuration

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add user to docker group
sudo usermod -aG docker $USER
```

### 1.2 Firewall Configuration

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow RTMP (restrict to trusted IPs in production)
sudo ufw allow from <trusted-ip> to any port 1935

# Enable firewall
sudo ufw enable
```

## Step 2: Database Setup

### Option A: Managed PostgreSQL (Recommended)

Use a managed PostgreSQL service (AWS RDS, DigitalOcean, etc.) and configure connection details.

### Option B: Self-Hosted PostgreSQL

```bash
# Create production docker-compose for database
cat > docker-compose.prod.yml << EOF
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: cyclingstream_postgres_prod
    environment:
      POSTGRES_USER: \${POSTGRES_USER}
      POSTGRES_PASSWORD: \${POSTGRES_PASSWORD}
      POSTGRES_DB: \${POSTGRES_DB}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G

volumes:
  postgres_data:
EOF
```

## Step 3: Environment Configuration

### 3.1 Create Environment Files

Create `.env` files for each component:

**Root `.env`** (for Docker Compose):
```bash
POSTGRES_USER=cyclingstream_prod
POSTGRES_PASSWORD=<strong-random-password>
POSTGRES_DB=cyclingstream_prod
PGADMIN_DEFAULT_EMAIL=admin@yourdomain.com
PGADMIN_DEFAULT_PASSWORD=<strong-random-password>
```

**Backend `.env`**:
```bash
PORT=8080
ENV=production

# Database
DB_HOST=<database-host>
DB_PORT=5432
DB_USER=cyclingstream_prod
DB_PASSWORD=<strong-random-password>
DB_NAME=cyclingstream_prod
DB_SSLMODE=require

# JWT Secret (generate: openssl rand -base64 32)
JWT_SECRET=<strong-random-secret-minimum-32-characters>

# Stripe
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Frontend URL
FRONTEND_URL=https://yourdomain.com
```

**Frontend `.env.local`**:
```bash
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

### 3.2 Generate Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Generate database password
openssl rand -base64 24

# Generate stream keys
openssl rand -hex 32
```

## Step 4: Application Deployment

### 4.1 Clone Repository

```bash
git clone https://github.com/yourusername/cyclingstream.git
cd cyclingstream
```

### 4.2 Build Backend

```bash
cd backend
go build -o bin/cyclingstream-api ./cmd/api/main.go
```

### 4.3 Build Frontend

```bash
cd frontend
npm ci
npm run build
```

### 4.4 Run Database Migrations

```bash
# Install migrate CLI if not already installed
# Download from: https://github.com/golang-migrate/migrate/releases

cd backend
migrate -path migrations -database "postgres://user:pass@host:5432/dbname?sslmode=require" up
```

## Step 5: Reverse Proxy Setup (Nginx)

### 5.1 Install Nginx

```bash
sudo apt install nginx certbot python3-certbot-nginx
```

### 5.2 Configure Nginx

Create `/etc/nginx/sites-available/cyclingstream`:

```nginx
# Backend API
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}

# Frontend
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    root /path/to/cyclingstream/frontend/out;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /_next/static {
        alias /path/to/cyclingstream/frontend/.next/static;
        add_header Cache-Control "public, max-age=31536000, immutable";
    }
}
```

### 5.3 Enable Site and Get SSL

```bash
sudo ln -s /etc/nginx/sites-available/cyclingstream /etc/nginx/sites-enabled/
sudo nginx -t
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com -d api.yourdomain.com
sudo systemctl reload nginx
```

## Step 6: Streaming Server Setup

### 6.1 Owncast Configuration

Update `stream/owncast/config.yaml`:

```yaml
cors:
  allowedOrigins:
    - "https://yourdomain.com"
    - "https://www.yourdomain.com"

server:
  port: 8080
  ip: "0.0.0.0"
```

Update `stream/owncast/docker-compose.yml` to use environment variables for ports.

### 6.2 Start Streaming Server

```bash
cd stream/owncast
docker-compose up -d
```

## Step 7: Process Management

### 7.1 Systemd Service for Backend

Create `/etc/systemd/system/cyclingstream-api.service`:

```ini
[Unit]
Description=CyclingStream API
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/cyclingstream/backend
ExecStart=/path/to/cyclingstream/backend/bin/cyclingstream-api
Restart=always
RestartSec=10
EnvironmentFile=/path/to/cyclingstream/backend/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable cyclingstream-api
sudo systemctl start cyclingstream-api
```

### 7.2 Systemd Service for Frontend

Create `/etc/systemd/system/cyclingstream-frontend.service`:

```ini
[Unit]
Description=CyclingStream Frontend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/cyclingstream/frontend
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=10
Environment=NODE_ENV=production
EnvironmentFile=/path/to/cyclingstream/frontend/.env.local

[Install]
WantedBy=multi-user.target
```

## Step 8: Monitoring and Logging

### 8.1 Log Management

```bash
# View backend logs
sudo journalctl -u cyclingstream-api -f

# View frontend logs
sudo journalctl -u cyclingstream-frontend -f

# View Docker logs
docker-compose logs -f
```

### 8.2 Health Checks

Set up monitoring to check:
- `https://api.yourdomain.com/health` - Should return `{"status":"ok"}`
- Database connectivity
- Disk space
- Memory usage

## Security Checklist

### Pre-Deployment

- [ ] All default passwords changed
- [ ] JWT secret is strong (32+ characters, random)
- [ ] Database uses SSL connections (`DB_SSLMODE=require`)
- [ ] All environment variables set correctly
- [ ] CORS configured for production domains only
- [ ] Firewall rules configured
- [ ] SSL certificates installed and auto-renewal configured
- [ ] Stream keys are secure and unique
- [ ] Stripe webhook secret configured
- [ ] Admin credentials are secure (not hardcoded)

### Post-Deployment

- [ ] HTTPS enforced (HTTP redirects to HTTPS)
- [ ] Security headers configured (HSTS, CSP, etc.)
- [ ] Rate limiting enabled and tested
- [ ] Database backups configured
- [ ] Log rotation configured
- [ ] Monitoring alerts set up
- [ ] Regular security updates scheduled

### Ongoing

- [ ] Regular security updates
- [ ] Monitor logs for suspicious activity
- [ ] Review and rotate secrets quarterly
- [ ] Database backups tested regularly
- [ ] SSL certificate renewal verified

## Troubleshooting

### Backend Won't Start

1. **Check configuration validation**:
   ```bash
   cd backend
   go run cmd/api/main.go
   ```
   Look for validation errors.

2. **Check database connection**:
   ```bash
   psql -h <host> -U <user> -d <database>
   ```

3. **Check logs**:
   ```bash
   sudo journalctl -u cyclingstream-api -n 50
   ```

### Database Connection Issues

1. **Verify credentials** in `.env` file
2. **Check firewall rules** - ensure database port is accessible
3. **Verify SSL mode** - production should use `require` or `verify-full`
4. **Test connection**:
   ```bash
   psql "postgres://user:pass@host:5432/dbname?sslmode=require"
   ```

### Frontend Build Issues

1. **Clear Next.js cache**:
   ```bash
   rm -rf frontend/.next
   npm run build
   ```

2. **Check environment variables**:
   ```bash
   cat frontend/.env.local
   ```

3. **Verify API URL** is correct and accessible

### Streaming Issues

1. **Check Owncast/Nginx logs**:
   ```bash
   docker-compose logs -f
   ```

2. **Verify stream key** is set correctly
3. **Check port conflicts** - ensure streaming server uses different port than backend
4. **Test RTMP connection** from OBS
5. **Verify HLS segments** are being generated

### CORS Errors

1. **Check CORS configuration** in backend routes
2. **Verify frontend URL** matches CORS allowed origins
3. **Check browser console** for specific CORS error
4. **Test with curl**:
   ```bash
   curl -H "Origin: https://yourdomain.com" \
        -H "Access-Control-Request-Method: GET" \
        -X OPTIONS \
        https://api.yourdomain.com/health
   ```

### Payment/Stripe Issues

1. **Verify Stripe keys** are for correct environment (live vs test)
2. **Check webhook endpoint** is accessible from Stripe
3. **Verify webhook secret** matches Stripe dashboard
4. **Check Stripe logs** in dashboard for webhook delivery status

## Backup Strategy

### Database Backups

```bash
# Create backup script
cat > /usr/local/bin/backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/var/backups/cyclingstream"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

pg_dump "postgres://user:pass@host:5432/dbname?sslmode=require" \
  | gzip > "$BACKUP_DIR/backup_$DATE.sql.gz"

# Keep only last 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
EOF

chmod +x /usr/local/bin/backup-db.sh

# Add to crontab (daily at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/backup-db.sh") | crontab -
```

### Application Backups

- Backup environment files (`.env` files)
- Backup configuration files
- Backup SSL certificates
- Store backups in secure, off-site location

## Scaling Considerations

### Horizontal Scaling

- Use load balancer for multiple backend instances
- Use managed database with read replicas
- Use CDN for static assets and HLS streams
- Consider Redis for session management and caching

### Vertical Scaling

- Monitor resource usage
- Upgrade server resources as needed
- Optimize database queries
- Implement caching where appropriate

## Maintenance

### Regular Tasks

- **Weekly**: Review logs, check disk space, verify backups
- **Monthly**: Security updates, dependency updates, performance review
- **Quarterly**: Secret rotation, security audit, capacity planning

### Update Procedure

1. Pull latest code
2. Run tests
3. Run database migrations
4. Build new binaries
5. Deploy with zero-downtime strategy (blue-green or rolling)
6. Verify health checks
7. Monitor for issues

## Support and Resources

- **Documentation**: See `README.md` and `API_DOCUMENTATION.md`
- **Issues**: GitHub Issues
- **Logs**: Systemd journals and Docker logs
- **Monitoring**: Set up external monitoring (UptimeRobot, Pingdom, etc.)

