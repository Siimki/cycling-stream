# WebSocket Error Debugging Guide

## Why the Error Object is Empty `{}`

The WebSocket `onerror` event in browsers **does not provide detailed error information**. This is a browser API limitation - the error object is always empty. The actual error details are available in the `onclose` event (close code, reason, etc.).

## Most Realistic Causes (Based on Code Analysis)

### 1. **Stream is Not Live** ⚠️ MOST LIKELY
**Location:** `backend/internal/handlers/chat.go:99-101`

The backend closes the WebSocket connection if the stream status is not "live":
```go
if stream == nil || stream.Status != "live" {
    sendErrorAndClose("Chat is only available for live races")
    return
}
```

**Solution:**
- Ensure the stream status is set to "live" before connecting to chat
- Use the admin API to update stream status: `PUT /admin/races/:id/stream/status` with `{"status": "live"}`
- Check the stream status first: `GET /races/:id/stream/status`

### 2. **Race Doesn't Exist**
**Location:** `backend/internal/handlers/chat.go:86-88`

If the race ID is invalid or doesn't exist, the connection is closed:
```go
if race == nil {
    sendErrorAndClose("Race not found")
    return
}
```

**Solution:**
- Verify the race ID is correct
- Ensure the race exists in the database

### 3. **WebSocket Upgrade Failed**
**Location:** `backend/internal/handlers/chat.go:46-49`

The backend checks if the request is a valid WebSocket upgrade:
```go
if !websocket.IsWebSocketUpgrade(c) {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "WebSocket upgrade required",
    })
}
```

**Possible causes:**
- Incorrect WebSocket URL format
- Missing or incorrect WebSocket headers
- Proxy/firewall blocking WebSocket upgrade

**Solution:**
- Verify the WebSocket URL: `ws://localhost:8080/races/{raceId}/chat/ws`
- Check browser console Network tab for the WebSocket connection
- Ensure no proxy is interfering

### 4. **CORS Issues**
**Location:** `backend/internal/server/routes.go:18-22`

While CORS is configured, WebSocket connections can still have CORS issues:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Content-Type,Authorization,X-CSRF-Token",
}))
```

**Note:** WebSocket connections don't use CORS in the same way as HTTP, but the upgrade request does.

**Solution:**
- Verify the frontend origin is allowed
- Check browser console for CORS errors
- Test with different browsers

### 5. **Backend Server Not Running or Unreachable**
If the backend is down or unreachable, the WebSocket connection will fail immediately.

**Solution:**
- Verify backend is running: `curl http://localhost:8080/health`
- Check if port 8080 is accessible
- Verify network connectivity

### 6. **Authentication Token Issues**
**Location:** `backend/internal/middleware/chat_auth.go`

While authentication is optional for WebSocket connections, invalid tokens might cause issues.

**Solution:**
- Check if token is valid (if provided)
- Try connecting without a token (anonymous)
- Verify token format and expiration

### 7. **Database Connection Issues**
**Location:** `backend/internal/handlers/chat.go:79-83`

If the database query fails, the connection is closed:
```go
race, err := h.raceRepo.GetByID(raceID)
if err != nil {
    logger.WithError(err).Error("Failed to fetch race for chat")
    sendErrorAndClose("Failed to verify race")
    return
}
```

**Solution:**
- Check backend logs for database errors
- Verify database is running and accessible
- Check database connection string

## How to Debug

### 1. Check Browser Console
The improved error handling now logs:
- Connection attempt details
- Close code and reason
- Ready state changes

### 2. Check Backend Logs
The backend logs WebSocket connection attempts:
```bash
# Check backend logs for:
# - "ChatAuthMiddleware: Request received"
# - "Failed to fetch race for chat"
# - "Chat is only available for live races"
```

### 3. Test WebSocket Connection Manually
```javascript
// In browser console:
const ws = new WebSocket('ws://localhost:8080/races/YOUR_RACE_ID/chat/ws');
ws.onopen = () => console.log('Connected');
ws.onerror = (e) => console.log('Error:', e);
ws.onclose = (e) => console.log('Closed:', e.code, e.reason);
ws.onmessage = (e) => console.log('Message:', e.data);
```

### 4. Verify Stream Status
```bash
# Check if stream exists and is live
curl http://localhost:8080/races/{raceId}/stream/status
```

### 5. Check Network Tab
- Open browser DevTools → Network tab
- Filter by "WS" (WebSocket)
- Look for the connection attempt
- Check the status code and response

## Common Close Codes

- **1000**: Normal closure
- **1006**: Abnormal closure (connection lost, no close frame received)
- **1001**: Going away (server shutting down)
- **1002**: Protocol error
- **1003**: Unsupported data type
- **1008**: Policy violation
- **1011**: Internal server error

## Next Steps

1. **Check the browser console** - The improved error handling will show close codes and reasons
2. **Verify stream status** - Ensure the stream is set to "live"
3. **Check backend logs** - Look for error messages from the chat handler
4. **Test with a known live race** - Use a race that definitely has a live stream

## References

- [MDN WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [WebSocket Close Codes](https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent)
- [Fiber WebSocket Documentation](https://docs.gofiber.io/guide/websockets)

