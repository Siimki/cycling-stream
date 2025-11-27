# Streaming Server Setup

This directory contains configuration files for streaming servers.

## Options

### Owncast (Recommended)

Owncast is a self-hosted live streaming server that's easier to set up and manage.

**Setup:**
1. Copy `stream/owncast/config.yaml` and customize
2. Set `OWNCAST_STREAM_KEY` environment variable
3. Run: `cd stream/owncast && docker-compose up -d`
4. RTMP URL: `rtmp://your-server-ip:1935/live`
5. HLS URL: `http://your-server-ip:8080/hls/stream.m3u8`

### Nginx-RTMP (Alternative)

Nginx with RTMP module provides more control but requires more configuration.

**Setup:**
1. Customize `stream/nginx/nginx.conf`
2. Run: `cd stream/nginx && docker-compose up -d`
3. RTMP URL: `rtmp://your-server-ip:1935/live/your-stream-key`
4. HLS URL: `http://your-server-ip:8080/hls/your-stream-key/index.m3u8`

## OBS Configuration

1. Open OBS Studio
2. Go to Settings → Stream
3. Service: Custom
4. Server: `rtmp://your-server-ip:1935/live` (Owncast) or `rtmp://your-server-ip:1935/live/your-stream-key` (Nginx)
5. Stream Key: Your configured stream key
6. Click OK and Start Streaming

## Testing

1. Start streaming from OBS
2. Check HLS URL in VLC or browser with HLS.js
3. Verify segments are being created
4. Test playback from different networks

## Security Notes

### Stream Key Management

- **Never use default stream keys in production**
- Generate secure random stream keys: `openssl rand -hex 32`
- Store stream keys securely (environment variables, secrets manager)
- Rotate stream keys regularly
- Use different stream keys for different races/events

### Network Security

- Use firewall rules to restrict RTMP access (port 1935) to trusted IPs only
- Restrict HLS access if needed (though it's typically public)
- Consider using VPN or private network for RTMP ingest

### HTTPS Configuration

For production deployments, configure HTTPS for HLS delivery:

**Owncast:**
- Use a reverse proxy (nginx, traefik) with SSL certificates
- Configure Owncast behind the proxy
- Update CORS settings to use HTTPS origins

**Nginx-RTMP:**
- Add SSL configuration to the `http` block
- Use Let's Encrypt or your SSL certificate provider
- Redirect HTTP to HTTPS

Example nginx SSL configuration:
```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location /hls {
        # ... HLS configuration ...
    }
}
```

### Stream Key Authentication

Currently, stream key authentication is handled at the nginx/owncast level. For enhanced security:

1. **Implement backend authentication endpoint** (`/auth/stream`) that:
   - Validates stream keys against database
   - Checks if the stream key is associated with an active race
   - Logs authentication attempts
   - Returns HTTP 200 for valid keys, 403 for invalid

2. **Update nginx configuration** to use the endpoint:
   ```nginx
   on_publish http://your-backend:8080/auth/stream;
   ```

3. **Rate limit** authentication attempts to prevent brute force

### Port Conflicts

- Backend API uses port 8080 by default
- Owncast/Nginx-RTMP HTTP also use port 8080
- **Solution**: Use different ports in production (e.g., 8081 for streaming)
- Update docker-compose.yml with `OWNCAST_HTTP_PORT` or `NGINX_HTTP_PORT` environment variables

