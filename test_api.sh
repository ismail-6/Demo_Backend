#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing Learning App API"
echo "========================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/health" | jq .
echo -e "\n"

# Test 2: Login
echo "2. Testing Login..."
curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test_user_123"}' | jq .
echo -e "\n"

# Test 3: Get All Chapters
echo "3. Getting All Chapters..."
curl -s "$BASE_URL/api/chapters" | jq .
echo -e "\n"

# Test 4: Get Chapter 1 Video
echo "4. Getting Chapter 1 Video..."
curl -s "$BASE_URL/api/chapters/1/video" | jq .
echo -e "\n"

# Test 5: Get Chapter 1 Quiz
echo "5. Getting Chapter 1 Quiz..."
curl -s "$BASE_URL/api/chapters/1/quiz" | jq .
echo -e "\n"

# Test 6: Save Video Progress
echo "6. Saving Video Progress (Chapter 1, 33 seconds)..."
curl -s -X POST "$BASE_URL/api/progress" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user_123",
    "chapter_id": 1,
    "content_type": "video",
    "video_timestamp": 33,
    "is_completed": false
  }' | jq .
echo -e "\n"

# Test 7: Save Quiz Progress
echo "7. Saving Quiz Progress (Chapter 1, Question 2)..."
curl -s -X POST "$BASE_URL/api/progress" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user_123",
    "chapter_id": 1,
    "content_type": "quiz",
    "quiz_question_index": 2,
    "is_completed": false
  }' | jq .
echo -e "\n"

# Test 8: Get User Progress
echo "8. Getting User Progress..."
curl -s "$BASE_URL/api/progress/user/test_user_123" | jq .
echo -e "\n"

# Test 9: Get All User Progress
echo "9. Getting All User Progress..."
curl -s "$BASE_URL/api/progress/user/test_user_123/all" | jq .
echo -e "\n"

# Test 10: Get Chapter-Specific Progress
echo "10. Getting Chapter 1 Progress..."
curl -s "$BASE_URL/api/progress/user/test_user_123/chapter/1" | jq .
echo -e "\n"

echo "All tests completed!"
