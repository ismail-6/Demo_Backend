# Quick Start: Deploy to Google Cloud Run

**5 minutes to deployment** ⚡

## Prerequisites

- Google Cloud account with billing enabled
- gcloud CLI installed (`gcloud --version`)
- PostgreSQL database URL ready

---

## Quick Deploy (3 Steps)

### 1. Install gcloud CLI (if not installed)

```bash
# Linux/macOS
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Verify
gcloud --version
```

### 2. Login and Setup

```bash
# Login
gcloud auth login

# Create/select project
gcloud config set project YOUR_PROJECT_ID
```

### 3. Edit & Deploy

**Edit** `deploy_cloudrun.sh`:
```bash
PROJECT_ID="your-gcp-project-id"              # ← Change this
DATABASE_URL="postgresql://user:pass@host/db" # ← Change this
```

**Run:**
```bash
chmod +x deploy_cloudrun.sh
./deploy_cloudrun.sh
```

**Done!** 🎉 Your API will be live in 3-5 minutes.

---

## What You'll Get

After deployment:
- ✅ Public HTTPS URL: `https://learnhub-backend-xxx.run.app`
- ✅ Auto-scaling (0 to 10 instances)
- ✅ Free SSL certificate
- ✅ Global CDN
- ✅ Auto-healing

---

## Test Your Deployment

```bash
# Save your URL
URL="https://learnhub-backend-xxx.run.app"

# Health check
curl $URL/health

# Login
curl -X POST $URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test"}'

# Get chapters
curl $URL/api/chapters
```

---

## Setup Database (First Time Only)

```bash
# Set your database URL
export DATABASE_URL="your-postgresql-connection-string"

# Run schema setup
psql $DATABASE_URL -f schema.sql
psql $DATABASE_URL -f seed_with_relations.sql
psql $DATABASE_URL -f quiz_answers_schema.sql
```

---

## Cost

**Free Tier:**
- 2 million requests/month FREE
- After that: ~$0.50-$2/month for low traffic

**Total with database:**
- $8-10/month (with small Cloud SQL instance)

---

## Common Issues

**❌ "Project not found"**
```bash
gcloud projects create your-project-id
gcloud config set project your-project-id
```

**❌ "Permission denied"**
```bash
gcloud auth login
```

**❌ "Database connection failed"**
- Check DATABASE_URL is correct
- Ensure database is publicly accessible or use Cloud SQL

---

## Update Deployment

```bash
# Make code changes, then:
./deploy_cloudrun.sh
```

Cloud Run will:
1. Build new image
2. Deploy with zero downtime
3. Keep old version as backup

---

## View Logs

```bash
gcloud run services logs tail learnhub-backend --region us-central1
```

Or visit: https://console.cloud.google.com/run

---

## Files Created

- ✅ `Dockerfile` - Container configuration
- ✅ `.dockerignore` - Build optimization
- ✅ `deploy_cloudrun.sh` - Deployment script
- ✅ `DEPLOY_GCP_CLOUD_RUN.md` - Full documentation
- ✅ `QUICKSTART_CLOUD_RUN.md` - This guide

---

## Next Steps

1. Update `deploy_cloudrun.sh` with your project ID and database URL
2. Run `./deploy_cloudrun.sh`
3. Test your API at the provided URL
4. Setup database tables (first time only)
5. Update frontend to use new API URL

**Need help?** Check `DEPLOY_GCP_CLOUD_RUN.md` for detailed documentation.
