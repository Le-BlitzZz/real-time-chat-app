
### Real-Time-Chat-Application - README

---

#### **Overview**

This is a real-time chat application built to support user registration, login, friend management, and private or global chat functionalities. It is designed to handle WebSocket communication for seamless real-time messaging.

---

#### **Tech Stack**

- **Backend:** Golang
  - Frameworks/Libraries: `Gin`, `Gorilla/WebSocket`, `Gorm`
- **Frontend:** HTML, JavaScript, CSS
- **Database:**
  - **SQL:** MySQL (used with Gorm for relational data like user accounts and friend management).
  - **NoSQL:** Redis (used for real-time chat data storage, including messages and chat metadata).
- **Containerization:** Docker

---

#### **Workflow**

1. **User Registration & Login:**
   - Users can register and log in with secure credentials.

2. **Friend Management:**
   - Users can send, accept, or reject friend requests.
   - Friend lists are dynamically updated.

3. **Global and Private Chats:**
   - **Global Chat:** Open to all users.
   - **Private Chat:** Created between two users upon request. Both users must join the chat for messages to be visible.

4. **Real-Time Messaging:**
   - Uses WebSocket connections to deliver instant messages.
   - Chats and their messages are stored in Redis for quick retrieval.

---

#### **Running the Application**

**To start the application using Docker, run:**

```bash
docker-compose -f docker-compose.docker.yml --project-name chatapp-docker up
```

**Application Workflow:**

- The `app` service runs on **port 8080** (mapped to localhost).
- MySQL database runs on **port 3306**.
- Redis runs on **port 6379**.

Access the application in your browser via: [http://localhost:8080](http://localhost:8080).

---

#### **Ports Used**

| Service  | Port    |
|----------|---------|
| App      | 8080    |
| MySQL    | 3306    |
| Redis    | 6379    |

---

#### **Requirements**

- Docker (with Docker Compose)

---

#### **Notes**

- The `config.docker.yml` file contains database connection details and default user credentials.
- To ensure the application works as expected, ensure that the required ports (8080, 3306, 6379) are available on your system.
