# Deployment Guide - LearnHub Backend

## Production Database

**PostgreSQL Database (Render.com)**
- Host: `dpg-d5ca45shg0os73e4ruu0-a.oregon-postgres.render.com`
- Database: `production_db_em4b`
- User: `production_db_em4b_user`
- Connection URL: Check `.env.example`

## Quick Start (Production)

### Option 1: Using Deploy Script

```bash
cd /home/ismail/Documents/Project/Demo_Backend
./deploy_production.sh
```

### Option 2: Manual Deployment

1. **Set Environment Variables**
   ```bash
   export DATABASE_URL="postgresql://production_db_em4b_user:aMDjUPiChzv5mfw2O70NyICuqf0IvqWc@dpg-d5ca45shg0os73e4ruu0-a.oregon-postgres.render.com/production_db_em4b"
   export PORT=8080
   export ENVIRONMENT=production
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Server**
   ```bash
   go run main.go
   ```

   Or build and run:
   ```bash
   go build -o learnhub-server main.go
   ./learnhub-server
   ```

## Local Development

For local development with SQLite:

```bash
# Don't set DATABASE_URL environment variable
# The app will automatically use SQLite

go run main.go
```

## Database Configuration

The app automatically detects the database type:

- **PostgreSQL**: If `DATABASE_URL` environment variable is set
- **SQLite**: If `DATABASE_URL` is not set (default for local dev)

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host/db` |
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment name | `production` or `development` |

## Features

✅ Automatic database migration
✅ Auto-seeding with sample data (if database is empty)
✅ Support for both PostgreSQL and SQLite
✅ CORS enabled for Flutter app
✅ Health check endpoint

## Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "message": "Learning App API is running",
  "environment": "production",
  "database": "postgres"
}
```

### API Base URL
```
http://localhost:8080/api
```

See [API_SUMMARY.md](API_SUMMARY.md) for all endpoints.

## Database Schema

The app will automatically create these tables:
- `users`
- `chapters`
- `videos`
- `quiz_questions`
- `progress`

## Sample Data

On first run, the database will be seeded with:
- 3 Chapters
- 3 Videos
- 15 Quiz Questions (5 per chapter)

## Deployment Platforms

### Render.com
1. Create new Web Service
2. Connect Git repository
3. Set build command: `go build -o server main.go`
4. Set start command: `./server`
5. Add environment variables from `.env.example`

### Railway.app
1. Create new project from GitHub
2. Add PostgreSQL database
3. Set environment variables
4. Deploy automatically

### Heroku
```bash
heroku create learnhub-api
heroku addons:create heroku-postgresql:mini
git push heroku main
```

## Monitoring

Check server status:
```bash
curl http://your-domain.com/health
```

View logs:
```bash
# If running with systemd
journalctl -u learnhub -f

# If running in terminal
# Logs are printed to stdout
```

## Security Notes

⚠️ **Important:**
- Change default CORS settings for production
- Use environment variables for sensitive data
- Never commit `.env` file to git
- Consider adding rate limiting
- Add authentication middleware for sensitive endpoints

## Troubleshooting

### Cannot connect to PostgreSQL
- Check firewall settings
- Verify DATABASE_URL is correct
- Ensure PostgreSQL service is running
- Check if IP is whitelisted

### Database migration failed
- Check database user permissions
- Verify connection string format
- Ensure database exists

### Port already in use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port
export PORT=3000
go run main.go
```

## Production Checklist

- [ ] Set up production PostgreSQL database
- [ ] Configure environment variables
- [ ] Test database connection
- [ ] Run backend server
- [ ] Test API endpoints
- [ ] Configure Flutter app to use production API
- [ ] Set up HTTPS/SSL
- [ ] Configure proper CORS origins
- [ ] Set up monitoring
- [ ] Configure automatic backups

## Support

For issues:
- Check logs for error messages
- Verify environment variables
- Test database connection separately
- Review [README.md](README.md) for full documentation
