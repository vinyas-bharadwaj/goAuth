# Deployment Guide

## Quick Production Deployment

### Using Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone <your-repo>
   cd GoAuth
   ```

2. **Configure environment:**
   
   **Option A - Using .env file:**
   ```bash
   cp .env.production.example .env
   ```
   
   Edit `.env`:
   ```env
   JWT_SECRET=$(openssl rand -base64 32)
   GOOGLE_CLIENT_ID=your-production-client-id.apps.googleusercontent.com
   ```
   
   **Option B - Direct environment variables in docker-compose.yml:**
   Edit `docker-compose.yml` and set environment variables directly:
   ```yaml
   environment:
     - JWT_SECRET=your-generated-secret
     - GOOGLE_CLIENT_ID=your-client-id
     - MONGODB_URI=mongodb://mongodb:27017
   ```

3. **Start services:**
   ```bash
   docker-compose up -d
   ```
   
   Note: The application will use environment variables from docker-compose.yml if no .env file is present.

4. **Verify:**
   ```bash
   docker-compose logs -f goauth
   ```

### Manual Deployment

1. **Build the binary:**
   ```bash
   make build
   ```

2. **Configure MongoDB:**
   - Ensure MongoDB is running
   - Update `MONGODB_URI` in `.env`

3. **Run the service:**
   ```bash
   ./bin/goauth
   ```

### Production Checklist

- [ ] Set strong `JWT_SECRET` (32+ random characters)
- [ ] Configure MongoDB authentication
- [ ] Enable TLS/SSL for gRPC
- [ ] Set up reverse proxy (nginx/envoy)
- [ ] Configure firewall rules
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure logging aggregation
- [ ] Set up automated backups for MongoDB
- [ ] Implement rate limiting
- [ ] Enable CORS for web clients
- [ ] Set proper `JWT_EXPIRES_IN` (15m recommended)
- [ ] Review and test all Google OAuth redirect URIs

### Monitoring

Monitor these metrics:
- Request rate and latency
- MongoDB connection pool status
- JWT validation failures
- Google OAuth verification failures
- Active user sessions

### Backup Strategy

**MongoDB Backup:**
```bash
# Automated daily backup
mongodump --uri="mongodb://localhost:27017/auth" --out=/backup/$(date +%Y%m%d)

# Restore
mongorestore --uri="mongodb://localhost:27017/auth" /backup/20260102
```

### Scaling

For high-traffic scenarios:
1. Use MongoDB replica sets for redundancy
2. Deploy multiple GoAuth instances behind a load balancer
3. Consider Redis for distributed token blacklisting
4. Implement caching for user lookups

### Security Hardening

1. **Network Security:**
   - Use TLS for all connections
   - Restrict MongoDB access to application only
   - Use VPC/private networks

2. **Application Security:**
   - Regular dependency updates
   - Security scanning (Snyk, Dependabot)
   - Rate limiting per IP/user
   - Input validation and sanitization

3. **Monitoring:**
   - Failed login attempts
   - Unusual traffic patterns
   - Token validation failures
